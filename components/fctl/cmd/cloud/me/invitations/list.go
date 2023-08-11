package invitations

import (
	"flag"
	"os"
	"time"

	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	statusFlag       = "status"
	organizationFlag = "organization"
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
func NewListConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useList, flag.ExitOnError)
	flags.String(statusFlag, "", "Filter invitations by status")
	flags.String(organizationFlag, "", "Filter invitations by organization")
	return fctl.NewControllerConfig(
		useList,
		shortList,
		shortList,
		[]string{
			"ls", "l",
		},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*ListStore] = (*ListController)(nil)

type ListController struct {
	store  *ListStore
	config fctl.ControllerConfig
}

func NewListController(config fctl.ControllerConfig) *ListController {
	return &ListController{
		store:  NewListStore(),
		config: config,
	}
}

func (c *ListController) GetStore() *ListStore {
	return c.store
}

func (c *ListController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *ListController) Run() (fctl.Renderable, error) {

	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}
	client, err := fctl.NewMembershipClient(flags, ctx, cfg)
	if err != nil {
		return nil, err
	}

	listInvitationsResponse, _, err := client.DefaultApi.
		ListInvitations(ctx).
		Status(fctl.GetString(flags, statusFlag)).
		Organization(fctl.GetString(flags, organizationFlag)).
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
	tableData = fctl.Prepend(tableData, []string{"ID", "Email", "Status", "CreationDate"})
	return pterm.DefaultTable.
		WithHasHeader().
		WithWriter(c.config.GetOut()).
		WithData(tableData).
		Render()

}

func NewListCommand() *cobra.Command {
	config := NewListConfig()
	return fctl.NewCommand("list",
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController[*ListStore](NewListController(*config)),
	)
}
