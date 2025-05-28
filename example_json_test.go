package confetti_test

import (
	"fmt"
	"os"

	"github.com/alexaandru/confetti"
)

func ExampleLoad_from_json_file() {
	jsonData := `{"Host":"localhost","Port":8080,"Debug":true,"Nested":{"Value":"foo", "Deep":{"Foo":"baz"}}}`

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
