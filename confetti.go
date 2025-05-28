package confetti

import (
	"errors"
	"fmt"
	"reflect"
)

// Loader is the interface implemented by all config loaders (env, SSM, JSON).
// You can implement your own Loader to support custom sources.
type Loader interface {
	Load(config any) error
}

// Load applies one or more loader functions to populate the given config struct.
//
// The first argument must be a pointer to a struct. Each loader (such as WithEnv, WithSSM, WithJSON)
// is applied in order, with later loaders overriding values from earlier ones.
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

	for _, ld := range append([]Loader{ld}, opts...) {
		if err = ld.Load(cfg); err != nil {
			return
		}
	}

	return
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

// WithJSON returns a loader that loads the config struct from a JSON file at the given path.
func WithJSON(path string) jsonLoader {
	return jsonLoader(path)
}
