package modelutils

import (
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

// This is used to propagate the window size from a parent to a child and sub childs
type ResizeMsg struct {
	Width  int
	Height int
}

type ErrorMsg struct {
	Error error
}
