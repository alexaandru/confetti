package confetti_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/alexaandru/confetti"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type mockSSM struct {
	value string
}

func ExampleLoad_ssm() {
	jsonValue := `{"Host":"ssmhost","Port":9000,"Debug":true,"Nested":{"Value":"ssmval","Deep":{"Foo":"ssmdeep", "Unknown":"unknown"}}}`
	ssmName := "CONFETTI_TEST"
	region := "us-east-1"

	cfg := &ExampleConfig{}
	err := confetti.Load(cfg,
		confetti.WithErrOnUnknown(),
		confetti.WithMockedSSM(&mockSSM{value: jsonValue}),
		confetti.WithSSM(ssmName, region),
	)

	fmt.Printf("Host=%s\n", cfg.Host)
	fmt.Printf("Port=%d\n", cfg.Port)
	fmt.Printf("Debug=%v\n", cfg.Debug)
	fmt.Printf("Nested.Value=%s\n", cfg.Nested.Value)
	fmt.Printf("Nested.Deep.Foo=%s\n", cfg.Nested.Deep.Foo)
	fmt.Printf("Error=%v\n", err)
	// Output:
	// Host=ssmhost
	// Port=9000
	// Debug=true
	// Nested.Value=ssmval
	// Nested.Deep.Foo=ssmdeep
	// Error=unknown fields in config: json: unknown field "Unknown"
}

func ExampleLoad_ssm_param_not_found() {
	cfg := &ExampleConfig{}
	err := confetti.Load(cfg,
		confetti.WithMockedSSM(&mockSSM{value: ""}),
		confetti.WithSSM("missing", "us-east-1"),
	)
	fmt.Printf("Error: %v\n", err)
	// Output:
	// Error: parameter missing not found or has no value
}

func ExampleLoad_ssm_error() {
	cfg := &ExampleConfig{}
	err := confetti.Load(cfg,
		confetti.WithMockedSSM(&mockSSM{value: "error: mock SSM error"}),
		confetti.WithSSM("fail", "us-east-1"),
	)
	fmt.Printf("Error: %v\n", err)
	// Output:
	// Error: failed to get SSM parameter fail: mock SSM error
}

func (m *mockSSM) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	if len(m.value) > 7 && m.value[:7] == "error: " {
		return nil, errors.New(m.value[7:])
	}

	if m.value == "" {
		return &ssm.GetParameterOutput{}, nil
	}

	return &ssm.GetParameterOutput{Parameter: &ssmtypes.Parameter{Value: &m.value}}, nil
}
