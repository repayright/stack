package modelutils

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/config"
)

type OpenPromptMsg struct{}
type ClosePromptMsg struct{}

type BlurMsg struct{}

type ChangeContextMsg struct{}

type ChangeViewMsg struct {
	Controller config.Controller
}

type ConfirmAskMsg struct {
	Question string
}

type ConfirmMsg struct {
	Confirm bool
}

type RenderMsg struct{}

type WindowSizeMsg struct {
	Width  int
	Height int
}

type ResizeMsg struct {
	Width  int
	Height int
}

type UpdateRendererMsg struct {
	Msg tea.Msg
}
