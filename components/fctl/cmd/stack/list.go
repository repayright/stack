package stack

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
)

// Every column in the table is a column in the struct
// The column name is suffixed with the max length of the column
// This is used to align the columns in the table
// e.g. maxLengthOrganizationId = 20
//
//	maxLengthStackId = 4
//
// Where the max length is the number of characters in the column name and values
const (
	maxLengthOrganizationId = 15
	maxLengthStackId        = 8
	maxLengthStackName      = 6
	maxLengthApiUrl         = 48
	maxLengthStackRegion    = 21
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
		fctl.WrapOutputPostRunE(formatViewStackListAsTable),
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

func formatViewStackListAsTable(cmd *cobra.Command, args []string) error {
	data, ok := fctl.GetSharedData().([]membershipclient.Stack)
	if !ok {
		return errors.New("invalid shared data")
	}

	organization, ok := fctl.GetSharedAdditionnalData(orgKey).(string)
	if !ok {
		return errors.New("invalid shared additional data")
	}

	// Default Columns
	columns := ui.NewArrayColumn(
		ui.NewColumn("Organization Id", maxLengthOrganizationId),
		ui.NewColumn("Stack Id", maxLengthStackId),
		ui.NewColumn("Name", maxLengthStackName),
		ui.NewColumn("API URL", maxLengthApiUrl),
		ui.NewColumn("Region", maxLengthStackRegion),
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

	// Default table options
	opts := ui.NewDefaultOptions(columns, tableData)

	// Add plain table option if --plain flag is set

	isPlain := fctl.GetString(cmd, fctl.OutputFlag) == "plain"

	var opt table.Option
	if isPlain {
		opt = ui.WithHeight(len(tableData))
	} else {
		opt = ui.WithHeight(ui.MaxTableHeight)
	}

	opts = append(opts, opt)

	t := ui.NewTableModel(!isPlain, opts...)

	return displayStackList(cmd, t)
}

func displayStackList(cmd *cobra.Command, t *ui.TableModel) error {
	isPlain := fctl.GetString(cmd, fctl.OutputFlag)

	if isPlain == "plain" {
		fmt.Println(t.View())
		return nil
	}

	if _, err := tea.NewProgram(t, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}
