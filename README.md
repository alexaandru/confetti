# Confetti 🎊

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Test](https://github.com/alexaandru/confetti/actions/workflows/ci.yml/badge.svg)](https://github.com/alexaandru/confetti/actions/workflows/ci.yml)
![Coverage](coverage-badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexaandru/confetti)](https://goreportcard.com/report/github.com/alexaandru/confetti)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexaandru/confetti.svg)](https://pkg.go.dev/github.com/alexaandru/confetti)
[![Socket.dev](https://socket.dev/api/badge/go/package/github.com/alexaandru/confetti)](https://socket.dev/go/package/github.com/alexaandru/confetti)

An opinionated take on handling Go configs: put your secrets in SSM, put the rest in ENV vars
and this package will load them all. And you can also load from JSON because... why not?

## Why Confetti?

- **Minimal API, maximal power:** One function (`Load()`), 3 `Loader` and 6 ways to load the data
  (ENV vars, SSM var holding a JSON, a local JSON file, a `[]byte` slice, an `io.Reader` or
  (preferred over `io.Reader`) an `io.ReadSeeker`);
- **Minimal dependencies:** Only SSM loader pulls in AWS SDK v2, stdlib for everything else;
- **Composability:** Layer environment variables, SSM, and JSON loaders in any order—later
  loaders override earlier ones;
- **Robust ENV var support:** The env variable names are inferred from the struct field name
  and the passed prefix (if non empty) and can also be overriden on a per-field basis using the
  struct tag `env` (e.g. `env:"MYAPP_FOO"`).
  Field names in CamelCase are converted to UPPER_SNAKE_CASE for environment variable lookup. Acronyms are handled so that `AWSRegion` becomes `AWS_REGION`, and `MyID` becomes `MY_ID`.
- **Robust type support:** When loading from env it handles primitives, slices, nested structs,
  booleans (with many/common string forms such as t/f, yes/no, etc.) and time durations out of the box;
- **Testable by example:** Code coverage is achieved with concise, real-world examples that
  double as documentation;
- **Bring Your Own Loader:** If builtin loaders don't fit your needs, you ~~can easily implement
  your own loader by implementing the `Loader` interface~~ (not at the moment, TBD);
- **Unknown field/var detection:** Optionally error if unknown fields/vars are present in the data
  but not in the target config, but ONLY AFTER the data has been loaded, so you can
  still use the config and just warn about the unknown fields.

## Available Loaders

| Loader           | Source Type         | Example Usage                                      |
| ---------------- | ------------------- | -------------------------------------------------- |
| WithErrOnUnknown | N/A                 | This sets the option to err on unknown fields/vars |
| WithEnv          | ENV prefix (string) | `WithEnv("MYAPP")`                                 |
| WithSSM          | SSM key (string)    | `WithSSM("/my/key", "us-east-1")`                  |
| WithJSON         | file path (string)  | `WithJSON("config.json")`                          |
| WithJSON         | []byte              | `WithJSON([]byte(jsonData))`                       |
| WithJSON         | io.ReadSeeker       | `WithJSON(bytes.NewReader(data))`                  |
| WithJSON         | io.Reader           | `WithJSON(os.Stdin)`                               |

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

### Strict Mode

```go
cfg := MyConfig{}
err := confetti.Load(&cfg, confetti.WithErrOnUnknown(), confetti.WithJSON("./config.json"))
if errors.Is(err, confetti.ErrUnknownFields) {
    // handle unknown fields error
}
```

### Default Values

No direct support for default values however, you can provide the `cfg` pre-populated
and you can also [Make the zero value useful](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=385s)
which together should cover most use cases.

Another option would be to use go:embed to embed a JSON file with the defaults,
while using other loaders to override it, i.e.

```go
//go:embed defaults.json
var defaultConfig []byte

func init() {
  if err := confetti.Load(&cfg,
    confetti.WithJSON(defaultConfig),
    confetti.WithJSON(".env.production.json")); err != nil {
    panic(err)
  }
}
```

For more examples see: [ENV](example_env_test.go), [JSON](example_json_test.go),
[ENV+JSON](example_both_test.go) and [SSM](example_ssm_test.go) Loader examples.

## License

[MIT](LICENSE)
