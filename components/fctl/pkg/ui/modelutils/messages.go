package modelutils

import "github.com/formancehq/fctl/pkg/config"

type BlurMsg struct{}

type ChangeViewMsg struct {
	controller config.Controller
}
