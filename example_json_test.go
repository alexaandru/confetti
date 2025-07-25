package confetti_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/alexaandru/confetti"
)

var jsonData = `{"Host":"localhost","Port":8080,"Debug":true,"Nested":{"Value":"foo", "Deep":{"Foo":"baz", "Unused":"unknown"}}}`

func ExampleLoad_json_file() {
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

func ExampleLoad_json_bytes() {
	cfg := &ExampleConfig{}
	if err := confetti.Load(cfg, confetti.WithJSON([]byte(jsonData))); err != nil {
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

func ExampleLoad_json_reader() {
	cfg := &ExampleConfig{}

	r := bytes.NewBufferString(jsonData)
	if err := confetti.Load(cfg, confetti.WithJSON(r)); err != nil {
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

func ExampleLoad_json_readseeker() {
	cfg := &ExampleConfig{}

	r := bytes.NewReader([]byte(jsonData))
	if err := confetti.Load(cfg, confetti.WithJSON(r)); err != nil {
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

func ExampleLoad_json_unsupported_type() {
	cfg := &ExampleConfig{}
	err := confetti.Load(cfg, confetti.WithJSON(123))
	fmt.Printf("Error: %v\n", err)
	// Output:
	// Error: unsupported type for WithJSON: int
}

func ExampleLoad_json_file_not_found() {
	cfg := &ExampleConfig{}
	err := confetti.Load(cfg, confetti.WithJSON("no_such_file.json"))
	fmt.Printf("Error: %v\n", err)
	// Output:
	// Error: open no_such_file.json: no such file or directory
}
