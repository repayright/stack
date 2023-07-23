package wallets

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"

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

func NewUpdateStore() *UpdateStore {
	return &UpdateStore{
		Success: false,
	}
}
func NewUpdateConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useUpdate, flag.ExitOnError)
	config.WithMetadataFlag(flags)
	config.WithConfirmFlag(flags)

	c := config.NewControllerConfig(
		useUpdate,
		shortUpdate,
		shortUpdate,
		[]string{
			"up",
		},
		flags,
		config.Organization, config.Stack,
	)

	return c
}

var _ config.Controller[*UpdateStore] = (*UpdateController)(nil)

type UpdateController struct {
	store  *UpdateStore
	config *config.ControllerConfig
}

func NewUpdateController(config *config.ControllerConfig) *UpdateController {
	return &UpdateController{
		store:  NewUpdateStore(),
		config: config,
	}
}

func (c *UpdateController) GetStore() *UpdateStore {
	return c.store
}

func (c *UpdateController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *UpdateController) Run() (config.Renderer, error) {
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

	if !fctl.CheckStackApprobation(flags, stack, "You are about to update a wallets") {
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
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController(NewUpdateController(c)),
	)
}
