package invitations

import (
	"flag"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useDecline   = "decline <invitation-id>"
	shortDecline = "Decline an invitation"
)

type DeclineStore struct {
	Success      bool   `json:"success"`
	InvitationId string `json:"invitationId"`
}

func NewDeclineStore() *DeclineStore {
	return &DeclineStore{}
}

func NewDeclineConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useDecline, flag.ExitOnError)
	fctl.WithConfirmFlag(flags)

	return config.NewControllerConfig(
		useDecline,
		shortDecline,
		shortDecline,
		[]string{
			"dec", "d",
		},
		flags,
		config.Organization, config.Stack,
	)
}

type DeclineController struct {
	store  *DeclineStore
	config *config.ControllerConfig
}

var _ config.Controller[*DeclineStore] = (*DeclineController)(nil)

func NewDeclineController(config *config.ControllerConfig) *DeclineController {
	return &DeclineController{
		store:  NewDeclineStore(),
		config: config,
	}
}

func (c *DeclineController) GetStore() *DeclineStore {
	return c.store
}

func (c *DeclineController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *DeclineController) Run() (modelutils.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()
	args := c.config.GetArgs()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	client, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	if !fctl.CheckOrganizationApprobation(flags, "You are about to decline an invitation") {
		return nil, fctl.ErrMissingApproval
	}

	_, err = client.DefaultApi.DeclineInvitation(ctx, args[0]).Execute()
	if err != nil {
		return nil, err
	}

	c.store.InvitationId = args[0]
	c.store.Success = true

	return c, nil
}

func (c *DeclineController) Render() error {
	pterm.Success.WithWriter(c.config.GetOut()).Printfln("Invitation declined! %s", c.store.InvitationId)
	return nil
}
func NewDeclineCommand() *cobra.Command {
	config := NewDeclineConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(1)),
		fctl.WithController[*DeclineStore](NewDeclineController(config)),
	)
}
