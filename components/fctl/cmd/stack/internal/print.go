package internal

import (
	"fmt"
	"io"

	blist "github.com/charmbracelet/bubbles/list"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui"
	"github.com/formancehq/fctl/pkg/ui/list"
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

func printInformation(out io.Writer, stack *membershipclient.Stack) (*list.ListModel, error) {

	items := []blist.Item{}

	items = append(items, list.NewItem(pterm.LightCyan("ID"), stack.Id))
	items = append(items, list.NewItem(pterm.LightCyan("Name"), stack.Name))
	items = append(items, list.NewItem(pterm.LightCyan("Region"), stack.RegionID))

	if ui, err := list.NewDefaultListModel(items, false); err != nil {
		return nil, err
	} else {
		return ui.WithTitle("Information"), nil
	}
}

func printVersion(out io.Writer, versions *shared.GetVersionsResponse, stack *membershipclient.Stack) (*list.ListModel, error) {
	items := []blist.Item{}

	for _, service := range versions.Versions {
		items = append(items, list.NewItem(
			pterm.LightCyan(fmt.Sprintf("%s (%s)", strcase.ToCamel(service.Name), pterm.Yellow(service.Version))),
			"",
		))
	}

	if ui, err := list.NewDefaultListModel(items, false); err != nil {
		return nil, nil
	} else {
		return ui.WithTitle("Version"), nil
	}
}

func printMetadata(out io.Writer, stack *membershipclient.Stack) (*list.ListModel, error) {
	items := []blist.Item{}

	for k, v := range stack.Metadata {
		items = append(items, list.NewItem(
			pterm.LightCyan(k),
			v,
		))

	}
	if ui, err := list.NewDefaultListModel(items, false); err != nil {
		return nil, nil
	} else {
		return ui.WithTitle("Metadata"), nil
	}
}

func PrintStackInformation(cmd *cobra.Command, profile *fctl.Profile, stack *membershipclient.Stack, versions *shared.GetVersionsResponse) (ui.Model, error) {
	out := cmd.OutOrStdout()
	content, err := getContent(out, stack, versions)
	if err != nil {
		return nil, err
	}

	// Plain
	if flag := fctl.GetString(cmd, fctl.OutputFlag); flag == "plain" {
		return ui.NewPlainOutput(content), nil
	}

	// Dynamic
	model, err := ui.NewViewPortManager(content, out)
	if err != nil {
		return nil, err
	}

	return model, nil
}
