package login

import (
	"flag"
	"fmt"
	"os"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	useLogin         = "login"
	descriptionLogin = "Login to the service"
)

type Dialog interface {
	DisplayURIAndCode(uri, code string)
}
type DialogFn func(uri, code string)

func (fn DialogFn) DisplayURIAndCode(uri, code string) {
	fn(uri, code)
}

type LoginStore struct {
	profile    *fctl.Profile `json:"-"`
	DeviceCode string        `json:"deviceCode"`
	LoginURI   string        `json:"loginUri"`
	BrowserURL string        `json:"browserUrl"`
	Success    bool          `json:"success"`
}

func NewDefaultLoginStore() *LoginStore {
	return &LoginStore{
		profile:    nil,
		DeviceCode: "",
		LoginURI:   "",
		BrowserURL: "",
		Success:    false,
	}
}

func NewLoginControllerConfig() *fctl.ControllerConfig {
	flags := flag.NewFlagSet(useLogin, flag.ExitOnError)
	flags.String(fctl.MembershipURIFlag, "", "service url")

	return fctl.NewControllerConfig(
		useLogin,
		descriptionLogin,
		[]string{},
		os.Stdout,
		flags,
	)
}

var _ fctl.Controller[*LoginStore] = (*LoginController)(nil)

type LoginController struct {
	store  *LoginStore
	config fctl.ControllerConfig
}

func NewLoginController(config fctl.ControllerConfig) *LoginController {
	return &LoginController{
		store:  NewDefaultLoginStore(),
		config: config,
	}
}

func (c *LoginController) GetStore() *LoginStore {
	return c.store
}

func (c *LoginController) GetConfig() fctl.ControllerConfig {
	return c.config
}

func (c *LoginController) Run() (fctl.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	profile := fctl.GetCurrentProfile(flags, cfg)
	membershipUri := fctl.GetString(flags, fctl.MembershipURIFlag)
	if membershipUri == "" {
		membershipUri = profile.GetMembershipURI()
	}

	relyingParty, err := fctl.GetAuthRelyingParty(fctl.GetHttpClient(flags, map[string][]string{}), membershipUri)
	if err != nil {
		return nil, err
	}

	c.store.profile = profile

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

func (c *LoginController) Render() error {

	fmt.Println("Please enter the following code on your browser:", c.store.DeviceCode)
	fmt.Println("Link:", c.store.LoginURI)

	if !c.store.Success && c.store.BrowserURL != "" {
		fmt.Printf("Unable to find a browser, please open the following link: %s", c.store.BrowserURL)
		return nil
	}

	if c.store.Success {
		pterm.Success.WithWriter(c.config.GetOut()).Printfln("Logged!")
	}

	return nil

}

func NewCommand() *cobra.Command {
	config := NewLoginControllerConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithShortDescription(config.GetDescription()),
		fctl.WithArgs(cobra.ExactArgs(0)),
		fctl.WithGoFlagSet(config.GetFlags()),
		fctl.WithPersistentGoFlagSet(config.GetPFlags()),
		fctl.WithController[*LoginStore](NewLoginController(*config)),
	)
}
