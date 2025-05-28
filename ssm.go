package confetti

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// ssmLoader loads config from an AWS SSM parameter containing a JSON string.
type ssmLoader struct {
	key       string
	awsRegion string
	profile   string
}

const DefaultAWSRegion = "us-east-1"

func (s ssmLoader) Load(config any) (err error) {
	cfgOpts := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(cmp.Or(s.awsRegion, DefaultAWSRegion))}

	if s.profile != "" {
		cfgOpts = append(cfgOpts, awsconfig.WithSharedConfigProfile(s.profile))
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(), cfgOpts...)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	svc := ssm.NewFromConfig(cfg)
	decrypted := true

	resp, err := svc.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: &s.key, WithDecryption: &decrypted,
	})
	if err != nil {
		return fmt.Errorf("failed to get SSM parameter %s: %w", s.key, err)
	}

	if resp.Parameter == nil || resp.Parameter.Value == nil {
		return fmt.Errorf("parameter %s not found or has no value", s.key)
	}

	return json.Unmarshal([]byte(*resp.Parameter.Value), config)
}
