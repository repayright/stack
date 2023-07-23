package instances

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	"github.com/formancehq/fctl/cmd/orchestration/instances/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	useDescribe         = "describe <instance-id>"
	descriptionDescribe = "Describe a specific workflow instance"
)

type DescribeStore struct {
	WorkflowInstancesHistory []shared.WorkflowInstanceHistory `json:"workflowInstanceHistory"`
}

func NewDescribeStore() *DescribeStore {
	return &DescribeStore{}
}
func NewDescribeConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useDescribe, flag.ExitOnError)

	c := config.NewControllerConfig(
		useDescribe,
		descriptionDescribe,
		descriptionDescribe,
		[]string{
			"des",
		},
		flags,
		config.Organization, config.Stack,
	)

	return c
}

type DescribeController struct {
	store  *DescribeStore
	client *formance.Formance
	config *config.ControllerConfig
}

var _ config.Controller[*DescribeStore] = (*DescribeController)(nil)

func NewDescribeController(config *config.ControllerConfig) *DescribeController {
	return &DescribeController{
		store:  NewDescribeStore(),
		config: config,
	}
}

func (c *DescribeController) GetStore() *DescribeStore {
	return c.store
}

func (c *DescribeController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *DescribeController) Run() (modelutils.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	soc, err := fctl.GetStackOrganizationConfig(flags, ctx, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(flags, ctx, soc.Config, soc.Stack, c.config.GetOut())
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.GetInstanceHistory(ctx, operations.GetInstanceHistoryRequest{
		InstanceID: args[0],
	})
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.client = client
	c.store.WorkflowInstancesHistory = response.GetWorkflowInstanceHistoryResponse.Data

	return c, nil
}

func (c *DescribeController) Render() error {

	for i, history := range c.store.WorkflowInstancesHistory {
		if err := internal.PrintStage(c.config.GetOut(), c.config.GetContext(), i, c.client, c.config.GetArgs()[0], history); err != nil {
			return err
		}
	}

	return nil
}

func NewDescribeCommand() *cobra.Command {
	config := NewDescribeConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DescribeStore](NewDescribeController(config)),
	)
}
