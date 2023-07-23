package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/ui/modelutils"
)

type PlainOutput struct {
	content string
}

func NewPlainOutput(content string) *PlainOutput {
	return &PlainOutput{
		content: content,
	}
}

func (po *PlainOutput) GetListKeyMapHandler() *modelutils.KeyMapHandler {
	return nil
}

func (po *PlainOutput) Init() tea.Cmd {
	return nil
}

func (po *PlainOutput) Update(msg tea.Msg) (Model, tea.Cmd) {
	return nil, nil
}

func (po *PlainOutput) View() string {
	return po.content
}
