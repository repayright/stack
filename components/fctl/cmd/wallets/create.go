package wallets

import (
	"flag"
	"fmt"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useCreateWallet         = "create <name>"
	descriptionCreateWallet = "Create a new wallet"
)

type CreateStore struct {
	WalletID string `json:"walletId"`
}
type CreateController struct {
	store  *CreateStore
	config fctl.ControllerConfig
}

var _ fctl.Controller[*CreateStore] = (*CreateController)(nil)

func NewDefaultCreateStore() *CreateStore {
	return &CreateStore{}
}
func NewCreateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useCreateWallet, flag.ExitOnError)
	fctl.WithMetadataFlag(flags)
	fctl.WithConfirmFlag(flags)
	return fctl.NewControllerConfig(
		useCreateWallet,
		descriptionCreateWallet,
		[]string{
			"cr",
		},
		os.Stdout,
		flags,
	)
}
func NewCreateController(config fctl.ControllerConfig) *CreateController {
	return &CreateController{
		store:  NewDefaultCreateStore(),
		config: config,
	}
}

func (c *CreateController) GetStore() *CreateStore {
	return c.store
}

func (c *CreateController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *CreateController) Run() (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to create a wallet") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(flags, metadataFlag))
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
		fctl.WithController[*CreateStore](NewCreateController(*c)),
	)
}
