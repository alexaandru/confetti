//nolint:testpackage // ok
package confetti

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type stubSSMClient struct{}

func (stubSSMClient) GetParameter(_ context.Context, _ *ssm.GetParameterInput, _ ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	return nil, nil //nolint:nilnil // stub, unused in this test
}

func TestWithSSMClient(t *testing.T) {
	t.Parallel()

	client := &stubSSMClient{}
	ld := WithSSMClient(client)

	c := &confetti{}
	if err := ld.Load(nil, c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.ssmClient != client {
		t.Fatalf("expected ownConfig.ssmClient to be set to the given client, got %v", c.ssmClient)
	}
}
