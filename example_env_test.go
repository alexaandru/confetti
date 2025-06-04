package confetti_test

import (
	"fmt"
	"os"
	"time"

	"github.com/alexaandru/confetti"
)

type ExampleConfig struct {
	Host   string
	Port   int
	Debug  bool
	Nested struct {
		Value string
		Deep  struct {
			Foo string
		}
	}
	Strs []string
	Ints []int
}

type ComplexConfig struct {
	Str    string
	Int    int
	Uint   uint
	Bool   bool
	Float  float64
	Strs   []string
	Ints   []int
	Uints  []uint
	Bools  []bool
	Floats []float64
	Dur    time.Duration
	Nested struct {
		Strs []string
		Deep struct {
			Int int
		}
	}
}

func ExampleLoad_env() {
	os.Setenv("MYAPP1_HOST", "127.0.0.1")
	os.Setenv("MYAPP1_PORT", "1234")
	os.Setenv("MYAPP1_DEBUG", "true")
	os.Setenv("MYAPP1_NESTED_VALUE", "bar")
	os.Setenv("MYAPP1_NESTED_DEEP_FOO", "baz")
	os.Setenv("MYAPP1_STRS", "a,b,c")
	os.Setenv("MYAPP1_INTS", "1,2,3")

	cfg := &ExampleConfig{}
	if err := confetti.Load(cfg, confetti.WithEnv("MYAPP1")); err != nil {
		panic(err)
	}

	fmt.Printf("Host=%s\n", cfg.Host)
	fmt.Printf("Port=%d\n", cfg.Port)
	fmt.Printf("Debug=%v\n", cfg.Debug)
	fmt.Printf("Nested.Value=%s\n", cfg.Nested.Value)
	fmt.Printf("Nested.Deep.Foo=%s\n", cfg.Nested.Deep.Foo)
	fmt.Printf("Strs=%#v\n", cfg.Strs)
	fmt.Printf("Ints=%#v\n", cfg.Ints)
	// Output:
	// Host=127.0.0.1
	// Port=1234
	// Debug=true
	// Nested.Value=bar
	// Nested.Deep.Foo=baz
	// Strs=[]string{"a", "b", "c"}
	// Ints=[]int{1, 2, 3}
}

func ExampleLoad_env_complex() {
	os.Setenv("CPLX_STR", "foo")
	os.Setenv("CPLX_INT", "42")
	os.Setenv("CPLX_UINT", "7")
	os.Setenv("CPLX_BOOL", "true")
	os.Setenv("CPLX_FLOAT", "3.14")
	os.Setenv("CPLX_STRS", "a,b,c")
	os.Setenv("CPLX_INTS", "1,2,3")
	os.Setenv("CPLX_UINTS", "4,5,6")
	os.Setenv("CPLX_FLOATS", "1.1,2.2,3.3")
	os.Setenv("CPLX_DUR", "1h30m")
	os.Setenv("CPLX_NESTED_STRS", "x,y")
	os.Setenv("CPLX_NESTED_DEEP_INT", "99")

	cfg := &ComplexConfig{}
	if err := confetti.Load(cfg, confetti.WithEnv("CPLX")); err != nil {
		panic(err)
	}

	fmt.Printf("Str=%s\n", cfg.Str)
	fmt.Printf("Int=%d\n", cfg.Int)
	fmt.Printf("Uint=%d\n", cfg.Uint)
	fmt.Printf("Bool=%v\n", cfg.Bool)
	fmt.Printf("Float=%.2f\n", cfg.Float)
	fmt.Printf("Strs=%#v\n", cfg.Strs)
	fmt.Printf("Ints=%#v\n", cfg.Ints)
	fmt.Printf("Uints=%#v\n", cfg.Uints)
	fmt.Printf("Bools=%#v\n", cfg.Bools)
	fmt.Printf("Floats=%#v\n", cfg.Floats)
	fmt.Printf("Dur=%s\n", cfg.Dur)
	fmt.Printf("Nested.Strs=%#v\n", cfg.Nested.Strs)
	fmt.Printf("Nested.Deep.Int=%d\n", cfg.Nested.Deep.Int)
	// Output:
	// Str=foo
	// Int=42
	// Uint=7
	// Bool=true
	// Float=3.14
	// Strs=[]string{"a", "b", "c"}
	// Ints=[]int{1, 2, 3}
	// Uints=[]uint{0x4, 0x5, 0x6}
	// Bools=[]bool(nil)
	// Floats=[]float64{1.1, 2.2, 3.3}
	// Dur=1h30m0s
	// Nested.Strs=[]string{"x", "y"}
	// Nested.Deep.Int=99
}

func ExampleLoad_env_unsupported_bools() {
	os.Setenv("CPLX_BOOLS", "true,false,yes,0,n")

	cfg := &ComplexConfig{}
	err := confetti.Load(cfg, confetti.WithEnv("CPLX"))
	fmt.Printf("%#v\n", cfg.Bools)
	fmt.Println(err)
	// Output:
	// []bool{true, false, true, false, false}
	// <nil>
}

func ExampleLoad_env_error_not_pointer() {
	cfg := ComplexConfig{} // not a pointer
	err := confetti.Load(cfg, confetti.WithEnv("MYAPP2"))

	fmt.Println(err)
	// Output:
	// config must be a pointer to a struct (got confetti_test.ComplexConfig)
}

func ExampleLoad_env_error_not_struct() {
	var x int

	err := confetti.Load(&x, confetti.WithEnv("MYAPP3"))

	fmt.Println(err)
	// Output:
	// config must be a pointer to a struct (got *int)
}

func ExampleLoad_env_error_parse_int() {
	os.Setenv("CPLX_INT", "notanint")

	cfg := &ComplexConfig{}
	err := confetti.Load(cfg, confetti.WithEnv("CPLX"))

	fmt.Println(err)
	// Output:
	// env CPLX_INT: strconv.ParseInt: parsing "notanint": invalid syntax
}

func ExampleLoad_env_struct_tag_override() {
	os.Setenv("CUSTOM_PORT", "9999")
	os.Setenv("MYAPP4_DEBUG", "true")
	os.Setenv("MYAPP4_UNUSED", "unknown") // should be ignored
	os.Setenv("CUSTOM_NESTED_VALUE", "tagged")

	type TaggedConfig struct {
		Port   int `env:"CUSTOM_PORT"`
		Debug  bool
		Nested struct {
			Value string `env:"CUSTOM_NESTED_VALUE"`
		}
	}

	cfg := &TaggedConfig{}
	err := confetti.Load(cfg, confetti.WithEnv("MYAPP4"))
	fmt.Printf("Port=%d Debug=%v Nested.Value=%s\n", cfg.Port, cfg.Debug, cfg.Nested.Value)
	fmt.Println(err)
	// Output:
	// Port=9999 Debug=true Nested.Value=tagged
	// <nil>
}

func ExampleLoad_env_error_on_unknown() {
	os.Setenv("MYAPP4_DEBUG", "true")
	os.Setenv("MYAPP4_UNUSED", "unknown") // should trigger an error
	os.Setenv("CUSTOM_PORT", "9999")
	os.Setenv("CUSTOM_NESTED_VALUE", "tagged")

	type TaggedConfig struct {
		Port   int `env:"CUSTOM_PORT"`
		Debug  bool
		Nested struct {
			Value string `env:"CUSTOM_NESTED_VALUE"`
		}
	}

	cfg := &TaggedConfig{}
	err := confetti.Load(cfg, confetti.WithErrOnUnknown(), confetti.WithEnv("MYAPP4"))
	fmt.Printf("Port=%d Debug=%v Nested.Value=%s\n", cfg.Port, cfg.Debug, cfg.Nested.Value)
	fmt.Println(err)
	// Output:
	// Port=9999 Debug=true Nested.Value=tagged
	// unknown environment variables: [MYAPP4_UNUSED]
}
