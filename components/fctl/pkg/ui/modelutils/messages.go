package modelutils

import "github.com/formancehq/fctl/pkg/config"

type BlurMsg struct{}

type ChangeViewMsg struct {
	controller config.Controller
}

type ConfirmAskMsg struct {
	question string
}

type ConfirmMsg struct {
	confirm bool
}
