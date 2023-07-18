package holds

import (
	"flag"
	"fmt"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useVoid   = "void <hold-id>"
	shortVoid = "Void a hold"
)

type VoidStore struct {
	Success bool   `json:"success"`
	HoldId  string `json:"holdId"`
}
type VoidController struct {
	store  *VoidStore
	config fctl.ControllerConfig
}

var _ fctl.Controller[*VoidStore] = (*VoidController)(nil)

func NewDefaultVoidStore() *VoidStore {
	return &VoidStore{}
}

func NewVoidConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useVoid, flag.ExitOnError)

	return fctl.NewControllerConfig(
		useVoid,
		shortVoid,
		shortVoid,
		[]string{
			"deb",
		},
		os.Stdout,
		flags,
	)

}

func NewVoidController(config fctl.ControllerConfig) *VoidController {
	return &VoidController{
		store:  NewDefaultVoidStore(),
		config: config,
	}
}

func (c *VoidController) GetStore() *VoidStore {
	return c.store
}

func (c *VoidController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *VoidController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID)
	if err != nil {
		return nil, err
	}

	stackClient, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	if len(c.config.GetArgs()) != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", len(c.config.GetArgs()))
	}

	request := operations.VoidHoldRequest{
		HoldID: c.config.GetArgs()[0],
	}
	response, err := stackClient.Wallets.VoidHold(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "voiding hold")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true //Todo: check status code 204/200
	c.store.HoldId = c.config.GetArgs()[0]

	return c, nil
}

func (c *VoidController) Render() error {

	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Hold '%s' voided!", c.config.GetArgs()[0])

	return nil
}

func NewVoidCommand() *cobra.Command {

	c := NewVoidConfig()

	return fctl.NewCommand(c.GetUse(),
		fctl.WithShortDescription(c.GetDescription()),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*VoidStore](NewVoidController(*c)),
	)
}
