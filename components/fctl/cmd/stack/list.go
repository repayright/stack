package stack

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	deletedFlag = "deleted"
	outputFlag  = "output"
	orgKey      = "organization"
	staticTable = true
)

const (
	maxLengthOrganizationId = 20
	maxLengthStackId        = 20
	maxLengthStackName      = 20
	maxLengthStackStatus    = 20
	maxLengthStackCreatedAt = 20
	maxLengthStackDeletedAt = 20
)

func NewListCommand() *cobra.Command {
	return fctl.NewMembershipCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List stacks"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithBoolFlag(deletedFlag, false, "Display deleted stacks"),
		fctl.WithRunE(listCommand),
		fctl.WrapOutputPostRunE(viewStackTable),
	)
}

func listCommand(cmd *cobra.Command, args []string) error {
	cfg, err := fctl.GetConfig(cmd)
	if err != nil {
		return err
	}

	profile := fctl.GetCurrentProfile(cmd, cfg)

	organization, err := fctl.ResolveOrganizationID(cmd, cfg)
	if err != nil {
		return errors.Wrap(err, "searching default organization")
	}

	apiClient, err := fctl.NewMembershipClient(cmd, cfg)
	if err != nil {
		return err
	}

	rsp, _, err := apiClient.DefaultApi.ListStacks(cmd.Context(), organization).
		Deleted(fctl.GetBool(cmd, deletedFlag)).
		Execute()
	if err != nil {
		return errors.Wrap(err, "listing stacks")
	}

	if len(rsp.Data) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No stacks found.")
		return nil
	}

	//Create map for addtionnal data containing the organization
	additionalData := map[string]interface{}{
		orgKey: organization,
	}

	fctl.SetSharedData(rsp.Data, profile, cfg, additionalData)

	return nil
}

func viewStackTable(cmd *cobra.Command, args []string) error {
	data := fctl.GetSharedData().([]membershipclient.Stack)
	organization := fctl.GetSharedAdditionnalData(orgKey).(string)

	// Default Columns
	columns := ui.NewArrayColumn(
		ui.NewColumn("Organization ID", maxLengthOrganizationId),
		ui.NewColumn("Stack ID", maxLengthStackId),
		ui.NewColumn("Name", maxLengthStackName),
		ui.NewColumn("API URL", maxLengthStackStatus),
		ui.NewColumn("Region", maxLengthStackCreatedAt),
		ui.NewColumn("Created At", maxLengthStackCreatedAt),
	)

	// Add Deleted At column if --deleted flag is set
	if fctl.GetBool(cmd, deletedFlag) {
		columns = columns.AddColumn("Deleted At", maxLengthStackDeletedAt)
	}

	// Create table data
	tableData := fctl.Map(data, func(stack membershipclient.Stack) table.Row {
		data := []string{
			organization,
			stack.Id,
			stack.Name,
			fctl.GetSharedProfile().ServicesBaseUrl(&stack).String(),
			stack.RegionID,
			stack.CreatedAt.Format(time.RFC3339),
		}

		if fctl.GetBool(cmd, deletedFlag) {
			if stack.DeletedAt != nil {
				data = append(data, stack.DeletedAt.Format(time.RFC3339))
			} else {
				data = append(data, "")
			}
		}

		return data
	})

	opts := []table.Option{
		table.WithColumns(columns),
		table.WithRows(tableData),
		table.WithFocused(true),
	}

	// Add static table option if --static flag is set
	opt := ui.WithStaticTable(len(tableData), staticTable)
	if opt != nil {
		opts = append(opts, opt)
	}

	t := ui.NewTableModel(opt != nil, opts...)

	fmt.Println(t.View())

	return nil
}
