package cloud

import (
	"flag"
	"fmt"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/pkg/ui/modelutils"

	fctl "github.com/formancehq/fctl/pkg"
	"github.com/spf13/cobra"
)

const (
	useGPT         = "generate-personal-token"
	shortGPT       = "Generate a personal bearer token"
	descriptionGPT = "Generate a personal bearer token"
)

type Store struct {
	Token string `json:"token"`
}

func NewStore() *Store {
	return &Store{}
}
func NewConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useGPT, flag.ExitOnError)
	return config.NewControllerConfig(
		useGPT,
		descriptionGPT,
		shortGPT,
		[]string{
			"gpt",
		},
		flags,
		config.Stack,
		config.Organization,
	)
}

var _ config.Controller[*Store] = (*Controller)(nil)

type Controller struct {
	store  *Store
	config *config.ControllerConfig
}

func NewController(config *config.ControllerConfig) *Controller {
	return &Controller{
		store:  NewStore(),
		config: config,
	}
}

func (c *Controller) GetStore() *Store {
	return c.store
}

func (c *Controller) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *Controller) Run() (modelutils.Renderable, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}
	profile := fctl.GetCurrentProfile(flags, cfg)

	organizationID, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	stack, err := fctl.ResolveStack(flags, ctx, cfg, organizationID, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	token, err := profile.GetStackToken(ctx, fctl.GetHttpClient(flags, map[string][]string{}, c.config.GetOut()), stack)
	if err != nil {
		return nil, err
	}

	c.store.Token = token

	return c, nil
}

func (c *Controller) Render() error {

	fmt.Fprintln(c.config.GetOut(), c.store.Token)
	return nil
}

func NewGeneratePersonalTokenCommand() *cobra.Command {
	config := NewConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithController[*Store](NewController(config)),
	)
}
