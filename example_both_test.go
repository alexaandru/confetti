package confetti_test

import (
	"fmt"
	"os"

	"github.com/alexaandru/confetti"
)

// ExampleConfig is assumed to be defined in another test file and available here.

func ExampleLoad_json_and_env() {
	// JSON provides Host and Nested.Value, ENV provides Port, Debug, and Nested.Deep.Foo
	jsonData := `{"Host":"localhost","Nested":{"Value":"foo"}}`

	file := "test_config.json"
	if err := os.WriteFile(file, []byte(jsonData), 0o644); err != nil {
		panic("failed to write test file: " + err.Error())
	}
	defer os.Remove(file)

	os.Setenv("MYAPP_PORT", "8080")
	os.Setenv("MYAPP_DEBUG", "true")
	os.Setenv("MYAPP_NESTED_DEEP_FOO", "baz")

	cfg := &ExampleConfig{}
	if err := confetti.Load(cfg,
		confetti.WithJSON(file),
		confetti.WithEnv("MYAPP"),
	); err != nil {
		panic(err)
	}

	fmt.Printf("Host=%s\n", cfg.Host)
	fmt.Printf("Port=%d\n", cfg.Port)
	fmt.Printf("Debug=%v\n", cfg.Debug)
	fmt.Printf("Nested.Value=%s\n", cfg.Nested.Value)
	fmt.Printf("Nested.Deep.Foo=%s\n", cfg.Nested.Deep.Foo)
	// Output:
	// Host=localhost
	// Port=8080
	// Debug=true
	// Nested.Value=foo
	// Nested.Deep.Foo=baz
}
