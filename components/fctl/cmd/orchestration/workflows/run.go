package workflows

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"
	"strings"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	"github.com/formancehq/fctl/cmd/orchestration/internal"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	variableFlag = "variable"
	waitFlag     = "wait"
	useRun       = "run <id>"
	shortRun     = "Run a workflow"
)

type RunStore struct {
	WorkflowInstance shared.WorkflowInstance `json:"workflowInstance"`
}

func NewRunStore() *RunStore {
	return &RunStore{}
}

func NewRunConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useRun, flag.ExitOnError)
	flags.Bool(waitFlag, false, "Wait end of the run")
	flags.String(variableFlag, "", "Variable to pass to the workflow")

	c := config.NewControllerConfig(
		useRun,
		shortRun,
		shortRun,
		[]string{
			"r",
		},
		flags,
		config.Organization, config.Stack,
	)

	return c
}

type RunController struct {
	store  *RunStore
	client *formance.Formance
	config *config.ControllerConfig
}

var _ config.Controller[*RunStore] = (*RunController)(nil)

func NewRunController(config *config.ControllerConfig) *RunController {
	return &RunController{
		store:  NewRunStore(),
		config: config,
	}
}

func (c *RunController) GetStore() *RunStore {
	return c.store
}

func (c *RunController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *RunController) Run() (modelutils.Renderable, error) {

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

	variables := make(map[string]string)
	for _, variable := range fctl.GetStringSlice(flags, variableFlag) {
		parts := strings.SplitN(variable, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("malformed flag: " + variable)
		}
		variables[parts[0]] = parts[1]
	}
	wait := fctl.GetBool(flags, waitFlag)
	response, err := client.Orchestration.
		RunWorkflow(ctx, operations.RunWorkflowRequest{
			RequestBody: variables,
			Wait:        &wait,
			WorkflowID:  args[0],
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

	c.store.WorkflowInstance = response.RunWorkflowResponse.Data
	c.client = client
	return c, nil
}

func (c *RunController) Render() error {
	flags := c.config.GetFlags()
	out := c.config.GetOut()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()
	wait := fctl.GetBool(flags, waitFlag)

	pterm.Success.WithWriter(out).Printfln("Workflow instance created with ID: %s", c.store.WorkflowInstance.ID)
	if wait {
		w, err := c.client.Orchestration.GetWorkflow(ctx, operations.GetWorkflowRequest{
			FlowID: args[0],
		})
		if err != nil {
			panic(err)
		}

		return internal.PrintWorkflowInstance(out, w.GetWorkflowResponse.Data, c.store.WorkflowInstance)
	}
	return nil
}

func NewRunCommand() *cobra.Command {
	c := NewRunConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*RunStore](NewRunController(c)),
	)
}
