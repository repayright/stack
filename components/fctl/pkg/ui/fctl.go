package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
	"github.com/pterm/pterm"
)

type FctlModel struct {
	content string
}

func NewFctlModel() *FctlModel {
	return &FctlModel{
		content: GetDefaultFCTLASCII("0.0.1"),
	}
}

func GetListKeyMapHandler() *modelutils.KeyMapHandler {
	return nil
}

// Generated with: https://patorjk.com/software/taag/#p=display&f=Doom&t=FCTL
func GetDefaultFCTLASCII(version string) string {
	t := "______ _____ _____ _  \n|  ___/  __ \\_   _| |\n| |_  | /  \\/ | | | |\n|  _| | |     | | | | \n| |   | \\__/\\ | | | |____\n\\_|    \\____/ \\_/ \\_____/\n"
	t = ApplyStyleToHeader(t)

	return t + "\nVersion : " + pterm.Red(version)
	// Get the version from the config file then compare it to the latest version, if semver is a major change display red, if minor display orange, if tiny display yellow
}
func ApplyStyleToHeader(header string) string {

	lines := strings.Split(header, "\n")

	for i := 0; i < len(lines); i++ {
		lines[i] = fctl.HeaderStyle.Sprint(lines[i])
	}

	return strings.Join(lines, "\n")
}

func (f FctlModel) GetMaxPossibleHeight() int {
	return len(strings.Split(f.content, "\n"))
}

func (f FctlModel) Init() tea.Cmd {
	return nil
}
func (f FctlModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return f, nil
}
func (f FctlModel) View() string {
	return f.content
}
