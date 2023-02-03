package activities

import (
	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Activities struct {
	client *sdk.APIClient
}

func New(client *sdk.APIClient) Activities {
	return Activities{
		client: client,
	}
}

func executeActivity(ctx workflow.Context, activity any, ret any, request any) error {
	if err := workflow.ExecuteActivity(ctx, activity, request).Get(ctx, ret); err != nil {
		var timeoutError *temporal.TimeoutError
		if errors.As(err, &timeoutError) {
			return errors.New(timeoutError.Message())
		}
		var applicationError *temporal.ApplicationError
		if errors.As(err, &applicationError) {
			return errors.New(applicationError.Message())
		}
		return err
	}
	return nil
}