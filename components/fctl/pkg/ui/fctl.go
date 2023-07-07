package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
	"github.com/formancehq/fctl/pkg/ui/theme"
	"github.com/pterm/pterm"
)

type FctlModel struct {
	content string
}

// Generated with: https://patorjk.com/software/taag/#p=display&f=Doom&t=FCTL
func GetDefaultFCTLASCII(version string) string {
	t := "______ _____ _____ _  \n|  ___/  __ \\_   _| |\n| |_  | /  \\/ | | | |\n|  _| | |     | | | | \n| |   | \\__/\\ | | | |____\n\\_|    \\____/ \\_/ \\_____/\n"
	t = ApplyStyleToString(t)

	return t + "\n" + "Version : " + pterm.Red(version)
	// Get the version from the config file then compare it to the latest version, if semver is a major change display red, if minor display orange, if tiny display yellow
}

// Export Style as a function argument, then export the function in utils
func ApplyStyleToString(header string) string {
	lines := strings.Split(header, "\n")

	for i := 0; i < len(lines); i++ {
		lines[i] = theme.HeaderStyle.Sprint(lines[i])
	}

	return strings.Join(lines, "\n")
}

func NewFctlModel() *FctlModel {
	return &FctlModel{
		content: GetDefaultFCTLASCII("0.0.1"),
	}
}

func GetListKeyMapHandler() *modelutils.KeyMapHandler {
	return nil
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
