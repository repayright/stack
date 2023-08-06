package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/ui/theme"
)

type Style struct {
	Header         lipgloss.Style
	HeaderSelected lipgloss.Style
	Row            lipgloss.Style
	RowSelected    lipgloss.Style

	Wrapper lipgloss.Style
	Body    lipgloss.Style
}

func NewStyle() *Style {
	return &Style{
		Header: lipgloss.NewStyle().
			BorderForeground(theme.TabBorderColor).
			Bold(false).PaddingTop(0).MarginTop(0).PaddingLeft(1).MarginRight(1),
		HeaderSelected: lipgloss.NewStyle().
			BorderForeground(theme.TabBorderColor).
			Bold(false).PaddingTop(0).MarginTop(0).PaddingLeft(1).MarginRight(1).Foreground(theme.SelectedColorForeground).
			Background(theme.SelectedColorForegroundBackground),
		Row: lipgloss.NewStyle().Foreground(theme.TabBorderColor).PaddingLeft(1),
		RowSelected: lipgloss.NewStyle().PaddingLeft(1).
			Foreground(theme.SelectedColorForeground).
			Background(theme.SelectedColorForegroundBackground).
			Bold(false).PaddingTop(0).MarginTop(0),
		Wrapper: lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(theme.TabBorderColor).Bold(false).PaddingTop(0).MarginTop(0),
		Body:    lipgloss.NewStyle().BorderForeground(theme.TabBorderColor).Bold(false).PaddingTop(0).MarginTop(0),
	}
}
