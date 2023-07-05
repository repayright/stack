package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

var (
	// Document <=> Wrap Content
	DocStyle = lipgloss.NewStyle().Padding(1, 2, 0, 2)

	// Tabs
	BaseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	TabBorderColor    = lipgloss.Color("240")
	InactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = tabBorderWithBottom("┘", " ", "└")
	InactiveTabStyle  = lipgloss.NewStyle().Border(InactiveTabBorder, true).BorderForeground(HighlightColor).Padding(0, 1)
	ActiveTabStyle    = InactiveTabStyle.Copy().Border(ActiveTabBorder, true)

	SelectedColorForeground           = lipgloss.Color("229")
	SelectedColorForegroundBackground = lipgloss.Color("57")

	// Viewport
	HighlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	ImportantColor = pterm.Yellow
	WindowStyle    = lipgloss.NewStyle().BorderForeground(lipgloss.NewStyle().GetBackground()).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240"))

	// List
	TitleStyle        = lipgloss.NewStyle().PaddingLeft(2)
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	QuitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	PaginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
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
