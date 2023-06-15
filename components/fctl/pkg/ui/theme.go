package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

var (
	BaseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	TabBorderColor                    = lipgloss.Color("240")
	SelectedColorForeground           = lipgloss.Color("229")
	SelectedColorForegroundBackground = lipgloss.Color("57")
	InactiveTabBorder                 = tabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder                   = tabBorderWithBottom("┘", " ", "└")
	DocStyle                          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	HighlightColor                    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	ImportantColor                    = pterm.Yellow
	InactiveTabStyle                  = lipgloss.NewStyle().Border(InactiveTabBorder, true).BorderForeground(HighlightColor).Padding(0, 1)
	ActiveTabStyle                    = InactiveTabStyle.Copy().Border(ActiveTabBorder, true)
	WindowStyle                       = lipgloss.NewStyle().BorderForeground(HighlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()

	TitleStyle        = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	QuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	//List Styles
	PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

// / COLOR PALETTE IMPORTANCY ///
var (
	InformationTitle = lipgloss.Color("#6c757d")
	InformationDesc  = pterm.NewStyle(pterm.FgLightCyan)
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
