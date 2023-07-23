package invitations

import (
	"flag"
	"github.com/formancehq/fctl/pkg/config"
	"time"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	statusFlag = "status"
)
const (
	useList   = "list"
	shortList = "List all invitations"
)

type Invitations struct {
	Id           string    `json:"id"`
	UserEmail    string    `json:"userEmail"`
	Status       string    `json:"status"`
	CreationDate time.Time `json:"creationDate"`
}

type ListStore struct {
	Invitations []Invitations `json:"invitations"`
}

func NewListStore() *ListStore {
	return &ListStore{}
}
func NewListConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	flags.String(statusFlag, "", "Filter invitations by status")

	return config.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"ls", "l",
		},
		flags,
		config.Organization, config.Stack,
	)
}

var _ config.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config *config.ControllerConfig
}

func NewListController(config *config.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (modelutils.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	listInvitationsResponse, _, err := apiClient.DefaultApi.
		ListOrganizationInvitations(ctx, organizationID).
		Status(fctl.GetString(flags, statusFlag)).
		Execute()
	if err != nil {
		return nil, err
	}

	c.store.Invitations = fctl.Map(listInvitationsResponse.Data, func(i membershipclient.Invitation) Invitations {
		return Invitations{
			Id:           i.Id,
			UserEmail:    i.UserEmail,
			Status:       i.Status,
			CreationDate: i.CreationDate,
		}
	})

	return c, nil
}

func (c *ListController) Render() error {
	tableData := fctl.Map(c.store.Invitations, func(i Invitations) []string {
		return []string{
			i.Id,
			i.UserEmail,
			i.Status,
			i.CreationDate.Format(time.RFC3339),
		}
	})

	tableData = fctl.Prepend(tableData, []string{"ID", "Email", "Status", "Creation date"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewListCommand() *cobra.Command {
	config := NewListConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(config)),
	)
}
