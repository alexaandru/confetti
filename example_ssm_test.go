//go:build testaws

package confetti_test

import (
	"context"
	"fmt"

	"github.com/alexaandru/confetti"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// ExampleConfig is assumed to be defined in another test file and available here.

func ExampleLoad_from_ssm() {
	jsonValue := `{"Host":"ssmhost","Port":9000,"Debug":true,"Nested":{"Value":"ssmval","Deep":{"Foo":"ssmdeep", "Unknown":"unknown"}}}`
	ssmName := "CONFETTI_TEST"
	region := "us-east-1"

	awsCfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		panic("failed to load AWS config: " + err.Error())
	}

	ssmClient := ssm.NewFromConfig(awsCfg)
	overwrite := true

	_, err = ssmClient.PutParameter(context.Background(), &ssm.PutParameterInput{
		Name:      &ssmName,
		Type:      "String",
		Value:     &jsonValue,
		Overwrite: &overwrite,
	})
	if err != nil {
		panic("failed to put SSM parameter: " + err.Error())
	}

	cfg := &ExampleConfig{}
	err = confetti.Load(cfg, confetti.WithErrOnUnknown(), confetti.WithSSM(ssmName, region))

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
