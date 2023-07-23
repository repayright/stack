package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/theme"
)

type Logo struct {
	content string
}

// Generated with: https://patorjk.com/software/taag/#p=display&f=Doom&t=FCTL
func GetDefaultFCTLASCII() string {
	t := "______ _____ _____ _  \n|  ___/  __ \\_   _| |\n| |_  | /  \\/ | | | |\n|  _| | |     | | | | \n| |   | \\__/\\ | | | |____\n\\_|    \\____/ \\_/ \\_____/"
	return t
}

func NewLogo() *Logo {
	return &Logo{
		content: GetDefaultFCTLASCII(),
	}
}

func (f Logo) GetMaxPossibleHeight() int {
	return len(strings.Split(f.content, "\n"))
}
func (f Logo) GetMaxPossibleWidth() int {
	tab := strings.Split(f.content, "\n")
	max := 0
	for _, line := range tab {
		if len(line) >= max {
			max = len(line)
		}
	}
	return max
}

func (f Logo) Init() tea.Cmd {
	return nil
}
func (f Logo) Update(msg tea.Msg) (Logo, tea.Cmd) {
	return f, nil
}

// Export Style as a function argument, then export the function in utils
func Style(header string) string {
	lines := strings.Split(header, "\n")
	style := lipgloss.NewStyle().Foreground(theme.LogoColor)

	for i := 0; i < len(lines); i++ {
		lines[i] = style.Render(lines[i])
	}

	//padding := lipgloss.NewStyle().PaddingTop(1)

	return strings.Join(lines, "\n")
}

func (f Logo) View() string {
	return lipgloss.Place(f.GetMaxPossibleWidth(), f.GetMaxPossibleHeight(), lipgloss.Left, lipgloss.Top, Style(f.content))
}
