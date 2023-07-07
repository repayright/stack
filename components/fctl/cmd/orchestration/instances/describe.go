package instances

import (
	"fmt"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type InstancesDescribeStore struct {
	WorkflowInstancesHistory []shared.WorkflowInstanceHistory `json:"workflow_instance_history"`
}
type InstancesDescribeController struct {
	store  *InstancesDescribeStore
	client *formance.Formance
}

var _ fctl.Controller[*InstancesDescribeStore] = (*InstancesDescribeController)(nil)

func NewDefaultInstancesDescribeStore() *InstancesDescribeStore {
	return &InstancesDescribeStore{}
}

func NewInstancesDescribeController() *InstancesDescribeController {
	return &InstancesDescribeController{
		store: NewDefaultInstancesDescribeStore(),
	}
}

func NewDescribeCommand() *cobra.Command {
	c := NewInstancesDescribeController()
	return fctl.NewCommand("describe <instance-id>",
		fctl.WithShortDescription("Describe a specific workflow instance"),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*InstancesDescribeStore](c),
	)
}

func (c *InstancesDescribeController) GetStore() *InstancesDescribeStore {
	return c.store
}

func (c *InstancesDescribeController) Run(cmd *cobra.Command, args []string) (fctl.Renderable, error) {
	soc, err := fctl.GetStackOrganizationConfig(cmd)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewStackClient(cmd, soc.Config, soc.Stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	response, err := client.Orchestration.GetInstanceHistory(cmd.Context(), operations.GetInstanceHistoryRequest{
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

func (c *InstancesDescribeController) Render(cmd *cobra.Command, args []string) (ui.Model, error) {

	for i, history := range c.store.WorkflowInstancesHistory {
		if err := printStage(cmd, i, c.client, args[0], history); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
