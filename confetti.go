package confetti

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)

// Loader is the interface implemented by all config loaders (env, SSM, JSON).
// You can implement your own Loader to support custom sources.
type Loader interface {
	Load(targetConfig any, ownConfig *confetti) error
}

type confetti struct {
	mockedSSM    SSMAPI
	errOnUnknown bool
}

// Load applies one or more loader functions to populate the given config which MUST be
// a pointer to a struct.
//
// The first argument must be a pointer to a struct. Each loader (such as WithEnv, WithSSM, WithJSON)
// is applied in order, with later loaders overriding values from earlier ones.
//
// You can optionally pass options.
// Currently, the only supported option is WithErrOnUnknown, which controls whether to return an error
// when unknown fields are present in the source but not defined in the target config.
//
// Returns an error if the config pointer is nil, not a struct, or if any loader fails.
//
// Example usage:
//
//	cfg := MyConfig{}
//	err := confetti.Load(&cfg, confetti.WithJSON("./config.json"), confetti.WithEnv("MYAPP"))
//	if err != nil { panic(err) }
func Load(cfg any, ld Loader, opts ...Loader) (err error) {
	if cfg == nil {
		return errors.New("config pointer cannot be nil")
	}

	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config must be a pointer to a struct (got %T)", cfg)
	}

	c, optx, ldx := confetti{}, []Loader{}, []Loader{}

	for _, ld := range append([]Loader{ld}, opts...) {
		switch ld.(type) {
		case optsLoader:
			optx = append(optx, ld)
		default:
			ldx = append(ldx, ld)
		}
	}

	for _, ld := range append(optx, ldx...) {
		if err = ld.Load(cfg, &c); err != nil {
			return
		}
	}

	return
}

// WithErrOnUnknown sets whether to return an error if is present in the source but
// not defined in the config struct.
// NOTE: It currently only applies to json and ssm loaders.
func WithErrOnUnknown() optsLoader {
	return optsLoader{errOnUnknown: true}
}

// WithMockedSSM returns a loader that uses a mocked SSM client for testing.
func WithMockedSSM(client SSMAPI) optsMockedSSMLoader {
	return optsMockedSSMLoader{client: client}
}

// WithEnv returns a loader that populates struct fields from environment variables.
//
// The prefix is prepended to each field name (in UPPER_SNAKE_CASE) to form the env var name.
//
// The optional separator argument sets the delimiter for slice fields (default is ",").
// Supports primitive types and slices of primitives (string, int, uint, float, bool).
func WithEnv(prefix string, opts ...string) envLoader {
	separator := DefaultSeparator
	if len(opts) > 0 {
		separator = opts[0]
	}

	return envLoader{prefix: prefix, separator: separator}
}

// WithSSM returns a loader that loads the config struct from an AWS SSM parameter.
//
// The key is the SSM parameter name. The optional region and profile arguments override the default AWS region/profile.
// The SSM parameter value must be a JSON string matching the config struct.
//
// Usage:
//
//	confetti.WithSSM("/my/param", "us-west-2", "myprofile")
func WithSSM(key string, opts ...string) ssmLoader {
	awsRegion, profile := DefaultAWSRegion, ""
	if len(opts) > 0 {
		awsRegion = opts[0]
	}

	if len(opts) > 1 {
		profile = opts[1]
	}

	return ssmLoader{key: key, awsRegion: awsRegion, profile: profile}
}

// WithJSON returns a loader that loads the config struct from a JSON source,
// which can be: a file path (string), []byte, io.ReadSeeker or io.Reader.
func WithJSON(src any) jsonLoader {
	switch v := src.(type) {
	case string:
		f, err := os.Open(v) //nolint:gosec // this is the whole point of the library.
		if err != nil {
			return jsonLoader{err: err}
		}

		return jsonLoader{r: f, c: f}
	case []byte:
		return jsonLoader{r: bytes.NewReader(v)}
	case io.ReadSeeker:
		return jsonLoader{r: v}
	case io.Reader:
		b, err := io.ReadAll(v)
		if err != nil {
			return jsonLoader{err: err}
		}

		return jsonLoader{r: bytes.NewReader(b)}
	default:
		return jsonLoader{err: fmt.Errorf("unsupported type for WithJSON: %T", src)}
	}
}
