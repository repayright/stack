package workflows

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	useCreate   = "create <file>|-"
	shortCreate = "Create a workflow"
)

type CreateStore struct {
	WorkflowId string `json:"workflowId"`
}

func NewCreateStore() *CreateStore {
	return &CreateStore{}
}
func NewCreateConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useCreate, flag.ExitOnError)

	c := config.NewControllerConfig(
		useCreate,
		shortCreate,
		shortCreate,
		[]string{
			"cr", "c",
		},
		flags,
		config.Organization, config.Stack,
	)

	return c
}

type CreateController struct {
	store  *CreateStore
	config *config.ControllerConfig
}

var _ config.Controller[*CreateStore] = (*CreateController)(nil)

func NewCreateController(config *config.ControllerConfig) *CreateController {
	return &CreateController{
		store:  NewCreateStore(),
		config: config,
	}
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) GetConfig() *config.ControllerConfig {
	return c.config
}
func (c *CreateController) Run() (modelutils.Renderable, error) {
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

	script, err := fctl.ReadFile(flags, soc.Stack, args[0])
	if err != nil {
		return nil, err
	}

	config := shared.WorkflowConfig{}
	if err := yaml.Unmarshal([]byte(script), &config); err != nil {
		return nil, err
	}

	//nolint:gosimple
	response, err := client.Orchestration.
		CreateWorkflow(ctx, shared.CreateWorkflowRequest{
			Name:   config.Name,
			Stages: config.Stages,
		})
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%s: %s", response.Error.ErrorCode, response.Error.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.WorkflowId = response.CreateWorkflowResponse.Data.ID

	return c, nil
}

func (c *CreateController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Workflow created with ID: %s", c.store.WorkflowId)

	return nil
}

func NewCreateCommand() *cobra.Command {
	config := NewCreateConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*CreateStore](NewCreateController(config)),
	)
}
