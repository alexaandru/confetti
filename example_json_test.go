package confetti_test

import (
	"fmt"
	"os"

	"github.com/alexaandru/confetti"
)

var jsonData = `{"Host":"localhost","Port":8080,"Debug":true,"Nested":{"Value":"foo", "Deep":{"Foo":"baz", "Unused":"unknown"}}}`

func ExampleLoad_from_json_file() {
	file := "test_config.json"
	if err := os.WriteFile(file, []byte(jsonData), 0o644); err != nil {
		panic("failed to write test file: " + err.Error())
	}
	defer os.Remove(file)

	cfg := &ExampleConfig{}
	if err := confetti.Load(cfg, confetti.WithJSON(file)); err != nil {
		panic("Load failed: " + err.Error())
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

func ExampleLoad_json_unknown_fields() {
	file := "test_config.json"
	if err := os.WriteFile(file, []byte(jsonData), 0o644); err != nil {
		panic("failed to write test file: " + err.Error())
	}
	defer os.Remove(file)

	cfg := &ExampleConfig{}
	err := confetti.Load(cfg, confetti.WithErrOnUnknown(), confetti.WithJSON(file))

	fmt.Printf("Host=%s\n", cfg.Host)
	fmt.Printf("Port=%d\n", cfg.Port)
	fmt.Printf("Debug=%v\n", cfg.Debug)
	fmt.Printf("Nested.Value=%s\n", cfg.Nested.Value)
	fmt.Printf("Nested.Deep.Foo=%s\n", cfg.Nested.Deep.Foo)
	fmt.Printf("Error=%s\n", err)
	// Output:
	// Host=localhost
	// Port=8080
	// Debug=true
	// Nested.Value=foo
	// Nested.Deep.Foo=baz
	// Error=unknown fields in config: json: unknown field "Unused"
}
