package confetti

import (
	"cmp"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMAPI is the minimal interface for SSM GetParameter used by ssmLoader.
type SSMAPI interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

// ssmLoader loads config from an AWS SSM parameter containing a JSON string.
type ssmLoader struct {
	key       string
	awsRegion string
	profile   string
}

const DefaultAWSRegion = "us-east-1"

func (s ssmLoader) Load(config any, ownConfig *confetti) (err error) {
	cfgOpts := []func(*awsconfig.LoadOptions) error{awsconfig.WithRegion(cmp.Or(s.awsRegion, DefaultAWSRegion))}

	if s.profile != "" {
		cfgOpts = append(cfgOpts, awsconfig.WithSharedConfigProfile(s.profile))
	}

	var svc SSMAPI

	if ownConfig != nil && ownConfig.mockedSSM != nil {
		svc = ownConfig.mockedSSM
	} else {
		var cfg aws.Config

		cfg, err = awsconfig.LoadDefaultConfig(context.Background(), cfgOpts...)
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}

		svc = ssm.NewFromConfig(cfg)
	}

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

	errOnUnknown := false

	if ownConfig != nil {
		errOnUnknown = ownConfig.errOnUnknown
	}

	return loadJSON(strings.NewReader(*resp.Parameter.Value), config, errOnUnknown)
}
