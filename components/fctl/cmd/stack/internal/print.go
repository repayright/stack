package internal

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/iancoleman/strcase"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func getContent(out io.Writer, stack *membershipclient.Stack, versions *shared.GetVersionsResponse) (string, error) {
	uiInfo, err := printInformation(out, stack)
	if err != nil {
		return "", err
	}

	uiVersion, err := printVersion(out, versions, stack)
	if err != nil {
		return "", err
	}

	uiMetadata, err := printMetadata(out, stack)
	if err != nil {
		return "", err
	}

	return uiInfo.View() + uiVersion.View() + uiMetadata.View(), nil
}

func printInformation(out io.Writer, stack *membershipclient.Stack) (*ui.ListModel, error) {

	items := []list.Item{}

	items = append(items, ui.NewItem(pterm.LightCyan("ID"), stack.Id))
	items = append(items, ui.NewItem(pterm.LightCyan("Name"), stack.Name))
	items = append(items, ui.NewItem(pterm.LightCyan("Region"), stack.RegionID))

	if ui, err := ui.NewDefaultListModel(items, false); err != nil {
		return nil, err
	} else {
		return ui.WithTitle("Information"), nil
	}
}

func printVersion(out io.Writer, versions *shared.GetVersionsResponse, stack *membershipclient.Stack) (*ui.ListModel, error) {
	items := []list.Item{}

	for _, service := range versions.Versions {
		items = append(items, ui.NewItem(
			pterm.LightCyan(fmt.Sprintf("%s (%s)", strcase.ToCamel(service.Name), pterm.Yellow(service.Version))),
			"",
		))
	}

	if ui, err := ui.NewDefaultListModel(items, false); err != nil {
		return nil, nil
	} else {
		return ui.WithTitle("Version"), nil
	}
}

func printMetadata(out io.Writer, stack *membershipclient.Stack) (*ui.ListModel, error) {
	items := []list.Item{}

	for k, v := range stack.Metadata {
		items = append(items, ui.NewItem(
			pterm.LightCyan(k),
			v,
		))

	}
	if ui, err := ui.NewDefaultListModel(items, false); err != nil {
		return nil, nil
	} else {
		return ui.WithTitle("Metadata"), nil
	}
}

func PrintStackInformation(cmd *cobra.Command, profile *fctl.Profile, stack *membershipclient.Stack, versions *shared.GetVersionsResponse) error {
	out := cmd.OutOrStdout()
	content, err := getContent(out, stack, versions)
	if err != nil {
		return err
	}

	// Static
	if flag := fctl.GetString(cmd, fctl.OutputFlag); flag == "static" {
		fctl.Println(content)
		return nil
	}

	// Dynamic
	model, err := ui.NewViewPortManager(content, out, profile, stack, versions)
	if err != nil {
		return err
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}
