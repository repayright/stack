package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/formancehq/fctl/pkg/theme"
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
			Foreground(theme.YellowColor).
			BorderForeground(theme.TabBorderColor).
			Bold(true).
			PaddingTop(0).
			MarginTop(0).
			PaddingLeft(1).
			MarginRight(1),
		HeaderSelected: lipgloss.NewStyle().
			Foreground(theme.YellowColor).
			BorderForeground(theme.TabBorderColor).
			Bold(true).
			PaddingTop(0).
			MarginTop(0).
			MarginLeft(1).
			MarginRight(1).
			Background(theme.SelectedHeaderForegroundBackground),
		Row: lipgloss.NewStyle().
			Foreground(theme.GreyColor).
			PaddingLeft(1),
		RowSelected: lipgloss.NewStyle().
			MarginLeft(1).
			Foreground(theme.WhiteColor).
			Background(theme.SelectedColorForegroundBackground).
			Bold(false).
			PaddingTop(0).
			MarginTop(0),
		Wrapper: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(theme.TabBorderColor).
			Bold(false).
			PaddingTop(0).
			MarginTop(0),
		Body: lipgloss.NewStyle().
			BorderForeground(theme.TabBorderColor).
			Bold(false).
			PaddingTop(0).
			MarginTop(0),
	}
}
