package internal

import (
	"fmt"
	"io"
	"net/url"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/iancoleman/strcase"

	"github.com/pterm/pterm"
)

func getContent(out io.Writer, stack *membershipclient.Stack, url *url.URL, versions *shared.GetVersionsResponse) string {
	return printInformation(out, stack).View() + printVersion(out, url, versions, stack).View() + printMetadata(out, stack).View()
}

func PrintStackInformation(out io.Writer, profile *fctl.Profile, stack *membershipclient.Stack, versions *shared.GetVersionsResponse) error {
	baseUrlStr := profile.ServicesBaseUrl(stack)
	content := getContent(out, stack, baseUrlStr, versions)

	model, err := ui.NewModelManager(content, out, profile, stack, versions)
	if err != nil {
		return err
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}

func printInformation(out io.Writer, stack *membershipclient.Stack) *ui.ListModel {

	items := []list.Item{}

	items = append(items, ui.NewItem(pterm.LightCyan("ID"), stack.Id))
	items = append(items, ui.NewItem(pterm.LightCyan("Name"), stack.Name))
	items = append(items, ui.NewItem(pterm.LightCyan("Region"), stack.RegionID))

	return ui.NewListModel(items, ui.ItemDelegate{}, ui.ViewWidth, ui.ViewHeight)
}

func printVersion(out io.Writer, url *url.URL, versions *shared.GetVersionsResponse, stack *membershipclient.Stack) *ui.ListModel {
	items := []list.Item{}

	for _, service := range versions.Versions {

		items = append(items, ui.NewItem(
			pterm.LightCyan(fmt.Sprintf("%s (%s)", strcase.ToCamel(service.Name), pterm.Yellow(service.Version))),
			fmt.Sprintf("%s/api/%s", url.String(), service.Name),
		))

	}

	return ui.NewListModel(items, ui.ItemDelegate{}, ui.ViewWidth, ui.ViewHeight).WithMaxPossibleHeight().WithMaxPossibleWidth()
}

func printMetadata(out io.Writer, stack *membershipclient.Stack) *ui.ListModel {
	items := []list.Item{}

	for k, v := range stack.Metadata {

		items = append(items, ui.NewItem(
			pterm.LightCyan(k),
			v,
		))

	}
	return ui.NewListModel(items, ui.ItemDelegate{}, ui.ViewWidth, ui.ViewHeight).WithMaxPossibleHeight().WithMaxPossibleWidth()
}
