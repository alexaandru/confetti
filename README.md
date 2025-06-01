# Confetti 🎊

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Build and Test](https://github.com/alexaandru/confetti/actions/workflows/ci.yml/badge.svg)](https://github.com/alexaandru/confetti/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/alexaandru/confetti/badge.svg)](https://coveralls.io/github/alexaandru/confetti)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexaandru/confetti)](https://goreportcard.com/report/github.com/alexaandru/confetti)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexaandru/confetti.svg)](https://pkg.go.dev/github.com/alexaandru/confetti)
[![Socket.dev](https://socket.dev/api/badge/go/package/github.com/alexaandru/confetti)](https://socket.dev/go/package/github.com/alexaandru/confetti)

An opinionated take on Go configuration: put your secrets in SSM, put the rest in ENV vars
and this package will load them all. And you can also load from JSON files because... why not?

## Why Confetti?

- **Minimal API, maximal power:** One function (`Load`) and three types of loaders;
- **Bring Your Own Loader:** If builtin loaders don't fit your needs, you can easily implement
  your own loader by implementing the `Loader` interface;
- **Composability:** Layer environment variables, SSM, and JSON loaders in any order—later
  loaders override earlier ones; The env variable names are inferred from the struct field name
  and the passed prefix (if non empty) and can also be overriden on a per-field basis using the
  struct tag `env` (e.g. `env:"MYAPP_FOO"`);
- **Robust type support:** Handles primitives, slices, nested structs, booleans (with many/common
  string forms such as t/f, yes/no, etc.) and time durations out of the box;
- **Testable by example:** Code coverage is achieved with concise, real-world examples that
  double as documentation;
- **Minimal dependencies:** Only SSM loader pulls in AWS SDK v2, stdlib for everything else.
- **Unknown field/var detection:** Optionally error if unknown fields/vars are present in the data
  but not in the target config.

```go
cfg := MyConfig{}
err := confetti.Load(&cfg, confetti.WithErrOnUnknown(), confetti.WithJSON("./config.json"))
if errors.Is(err, confetti.ErrUnknownFields) {
    // handle unknown fields error
}
```

## Usage

```go
import "github.com/alexaandru/confetti"

var cfg MyConfig

func init() {
  if err := confetti.Load(&cfg, confetti.WithEnv("MYAPP")); err != nil {
    panic(err)
  }
}
```

For more examples see:

- [ENV Loader Example](example_env_test.go);
- [JSON Loader Example](example_json_test.go);
- [ENV+JSON Loader Example](example_both_test.go);
- [SSM Loader Example](example_ssm_test.go).

## License

MIT
