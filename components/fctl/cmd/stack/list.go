package stack

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	deletedFlag = "deleted"
	outputFlag  = "output"
)

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func NewListCommand() *cobra.Command {
	return fctl.NewMembershipCommand("list",
		fctl.WithAliases("ls", "l"),
		fctl.WithShortDescription("List stacks"),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithBoolFlag(deletedFlag, false, "Display deleted stacks"),
		fctl.WithStringFlag(outputFlag, "", "Output to the resquested format (json)"),
		fctl.WithRunE(listCommand),
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

	//Output to requested format
	if flags := fctl.GetString(cmd, outputFlag); flags != "" {
		switch flags {
		case "json":
			// Marshal to JSON then print to stdout
			m, err := rsp.MarshalJSON()
			if (err) != nil {
				return errors.Wrap(err, "marshalling to json")
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), string(m))

			return err
		}
	}

	// Default Columns
	columns := []table.Column{
		{Title: "Organization ID", Width: 15},
		{Title: "Stack ID", Width: 10},
		{Title: "Name", Width: 10},
		{Title: "API URL", Width: 50},
		{Title: "Region", Width: 40},
		{Title: "Created At", Width: 25},
	}

	tableData := fctl.Map(rsp.Data, func(stack membershipclient.Stack) table.Row {
		data := []string{
			organization,
			stack.Id,
			stack.Name,
			profile.ServicesBaseUrl(&stack).String(),
			stack.RegionID,
			stack.CreatedAt.Format(time.RFC3339),
		}

		if fctl.GetBool(cmd, deletedFlag) {
			columns = append(columns, table.Column{Title: "Deleted At", Width: 25})
			if stack.DeletedAt != nil {
				data = append(data, stack.DeletedAt.Format(time.RFC3339))
			} else {
				data = append(data, "")
			}
		}

		return data
	})

	//Table max height of 10
	maxHeight := If(len(tableData)+1 > 10, 10, len(tableData)+1)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(tableData),
		table.WithFocused(true),
		table.WithHeight(maxHeight),
	)

	//Default styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	// Print model
	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
	}

	return err
}
