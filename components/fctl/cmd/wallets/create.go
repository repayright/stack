package wallets

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
)

const (
	useCreate   = "create <name>"
	shortCreate = "Create a new wallet"
)

type CreateStore struct {
	WalletID string `json:"walletId"`
}
type CreateController struct {
	store  *CreateStore
	config *config.ControllerConfig
}

var _ config.Controller[*CreateStore] = (*CreateController)(nil)

func NewCreateStore() *CreateStore {
	return &CreateStore{}
}
func NewCreateConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useCreate, flag.ExitOnError)
	config.WithMetadataFlag(flags)
	config.WithConfirmFlag(flags)
	return config.NewControllerConfig(
		useCreate,
		shortCreate,
		shortCreate,
		[]string{
			"cr",
		},
		flags,
		config.Organization, config.Stack,
	)
}
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
	out := c.config.GetOut()
	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config")
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, out)
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, out)
	if err != nil {
		return nil, err
	}

	if !fctl.CheckStackApprobation(flags, stack, "You are about to create a wallet") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack, out)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	metadata, err := fctl.ParseMetadata(config.GetStringSlice(flags, config.MetadataFlag))
	if err != nil {
		return nil, err
	}

	if len(c.config.GetArgs()) < 1 {
		return nil, fmt.Errorf("missing wallet name")
	}

	request := shared.CreateWalletRequest{
		Name:     c.config.GetArgs()[0],
		Metadata: metadata,
	}
	response, err := client.Wallets.CreateWallet(ctx, request)
	if err != nil {
		return nil, errors.Wrap(err, "creating wallet")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.WalletID = response.CreateWalletResponse.Data.ID

	return c, nil
}

func (c *CreateController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln(
		"Wallet created successfully with ID: %s", c.store.WalletID)
	return nil
}
func NewCreateCommand() *cobra.Command {
	c := NewCreateConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithShortDescription(c.GetDescription()),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*CreateStore](NewCreateController(c)),
	)
}
