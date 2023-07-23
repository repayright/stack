package login

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/config"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useLogin         = "login"
	descriptionLogin = "Login to the service"
)

type Store struct {
	DeviceCode string `json:"deviceCode"`
	LoginURI   string `json:"loginUri"`
	BrowserURL string `json:"browserUrl"`
	Success    bool   `json:"success"`
}

func NewStore() *Store {
	return &Store{
		DeviceCode: "",
		LoginURI:   "",
		BrowserURL: "",
		Success:    false,
	}
}

func NewConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useLogin, flag.ExitOnError)
	flags.String(config.MembershipURIFlag, "", "service url")

	return config.NewControllerConfig(
		useLogin,
		descriptionLogin,
		descriptionLogin,
		[]string{
			"log",
		},
		flags,
	)
}

var _ config.Controller = (*LoginController)(nil)

type LoginController struct {
	store  *Store
	config *config.ControllerConfig
}

func (c *LoginController) GetKeyMapAction() *config.KeyMapHandler {
	return nil
}

func NewController(config *config.ControllerConfig) *LoginController {
	return &LoginController{
		store:  NewStore(),
		config: config,
	}
}

func (c *LoginController) GetStore() any {
	return c.store
}

func (c *LoginController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *LoginController) Run() (config.Renderer, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)
	membershipUri := config.GetString(flags, config.MembershipURIFlag)
	if membershipUri == "" {
		membershipUri = profile.GetMembershipURI()
	}

	relyingParty, err := fctl.GetAuthRelyingParty(fctl.GetHttpClient(flags, map[string][]string{}, c.config.GetOut()), membershipUri)
	if err != nil {
		return nil, err
	}

	ret, err := LogIn(ctx, DialogFn(func(uri, code string) {
		c.store.DeviceCode = code
		c.store.LoginURI = uri
	}), relyingParty)

	// Other relying error not related to browser
	if err != nil && err.Error() != "error_opening_browser" {
		return nil, err
	}

	// Browser not found
	if err == nil {
		c.store.Success = true
	}

	profile.SetMembershipURI(membershipUri)
	profile.UpdateToken(ret)

	currentProfileName := fctl.GetCurrentProfileName(flags, cfg)

	cfg.SetCurrentProfile(currentProfileName, profile)

	return c, cfg.Persist()
}

func (c *LoginController) Render() (tea.Model, error) {
	out := c.config.GetOut()
	fmt.Fprintln(out, "Please enter the following code on your browser:", c.store.DeviceCode)
	fmt.Fprintln(out, "Link:", c.store.LoginURI)

	if !c.store.Success && c.store.BrowserURL != "" {
		fmt.Fprintf(out, "Unable to find a browser, please open the following link: %s", c.store.BrowserURL)
		return nil, nil
	}

	if c.store.Success {
		pterm.Success.WithWriter(c.config.GetOut()).Printfln("Logged!")
	}

	return nil, nil

}

func NewCommand() *cobra.Command {
	c := NewConfig()
	return fctl.NewCommand(c.GetUse(),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithController(NewController(c)),
	)
}
