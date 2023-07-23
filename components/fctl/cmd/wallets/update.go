package wallets

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
	useUpdate   = "update <wallet-id>"
	shortUpdate = "Update a wallets"
)

type UpdateStore struct {
	Success bool `json:"success"`
}

func NewDefaultUpdateStore() *UpdateStore {
	return &UpdateStore{
		Success: false,
	}
}
func NewUpdateConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useUpdate, flag.ExitOnError)
	flags.String("metadata", "", "Metadata to use")
	fctl.WithConfirmFlag(flags)

	c := fctl.NewControllerConfig(
		useUpdate,
		shortUpdate,
		[]string{
			"up",
		},
		os.Stdout,
		flags,
	)

	c.SetShortDescription(shortUpdate)
	return c
}

var _ fctl.Controller[*UpdateStore] = (*UpdateController)(nil)

type UpdateController struct {
	store        *UpdateStore
	metadataFlag string
	config       fctl.ControllerConfig
}

func NewUpdateController(config fctl.ControllerConfig) *UpdateController {
	return &UpdateController{
		store:        NewDefaultUpdateStore(),
		metadataFlag: "metadata",
		config:       config,
	}
}

func (c *UpdateController) GetStore() *UpdateStore {
	return c.store
}

func (c *UpdateController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *UpdateController) Run() (fctl.Renderable, error) {
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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to update a wallets") {
		return nil, fctl.ErrMissingApproval
	}

	client, err := fctl.NewStackClient(flags, ctx, cfg, stack)
	if err != nil {
		return nil, errors.Wrap(err, "creating stack client")
	}

	metadata, err := fctl.ParseMetadata(fctl.GetStringSlice(flags, c.metadataFlag))
	if err != nil {
		return nil, err
	}

	if len(c.config.GetArgs()) == 0 {
		return nil, fmt.Errorf("wallet id is required")
	}

	response, err := client.Wallets.UpdateWallet(ctx, operations.UpdateWalletRequest{
		RequestBody: &operations.UpdateWalletRequestBody{
			Metadata: metadata,
		},
		ID: c.config.GetArgs()[0],
	})
	if err != nil {
		return nil, errors.Wrap(err, "updating wallet")
	}

	if response.WalletsErrorResponse != nil {
		return nil, fmt.Errorf("%s: %s", response.WalletsErrorResponse.ErrorCode, response.WalletsErrorResponse.ErrorMessage)
	}

	if response.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	c.store.Success = true
	return c, nil
}

func (c *UpdateController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Wallet updated successfully!")
	return nil
}
func NewUpdateCommand() *cobra.Command {
	c := NewUpdateConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithShortDescription(c.GetDescription()),
		fctl.WithAliases(c.GetAliases()...),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*UpdateStore](NewUpdateController(*c)),
	)
}
