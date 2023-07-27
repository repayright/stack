package stack

import (
	"flag"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/formancehq/fctl/pkg/config"

	"github.com/formancehq/fctl/cmd/stack/internal"
	"github.com/formancehq/fctl/membershipclient"
	fctl "github.com/formancehq/fctl/pkg"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	unprotectFlag = "unprotect"
	regionFlag    = "region"
	nowaitFlag    = "no-wait"
	useCreate     = "create [name]"
	shortCreate   = "Create a new stack"
)

type CreateStore struct {
	Stack    *membershipclient.Stack     `json:"stack"`
	Versions *shared.GetVersionsResponse `json:"versions"`
}

func NewCreateStore() *CreateStore {
	return &CreateStore{
		Stack:    &membershipclient.Stack{},
		Versions: &shared.GetVersionsResponse{},
	}
}

func NewStackConfig() *config.ControllerConfig {
	flags := flag.NewFlagSet(useCreate, flag.ExitOnError)
	flags.Bool(unprotectFlag, false, "Unprotect stacks (no confirmation on write commands)")
	flags.String(regionFlag, "", "Region on which deploy the stack")
	flags.Bool(nowaitFlag, false, "Not wait stack availability")

	return config.NewControllerConfig(
		useCreate,
		shortCreate,
		shortCreate,
		[]string{
			"cr",
			"c",
		},
		flags,
		config.Organization,
	)
}

var _ config.Controller = (*CreateController)(nil)

type CreateController struct {
	store   *CreateStore
	profile *fctl.Profile
	config  *config.ControllerConfig
}

func (c *CreateController) GetKeyMapAction() *config.KeyMapHandler {
	return nil
}

func NewStackController(config *config.ControllerConfig) *CreateController {
	return &CreateController{
		store:  NewCreateStore(),
		config: config,
	}
}

func (c *CreateController) GetStore() any {
	return c.store
}

func (c *CreateController) GetConfig() *config.ControllerConfig {
	return c.config
}

func (c *CreateController) Run() (config.Renderer, error) {
	flags := c.config.GetAllFLags()
	ctx := c.config.GetContext()

	cfg, err := fctl.GetConfig(flags)
	if err != nil {
		return nil, err
	}

	organization, err := fctl.ResolveOrganizationID(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	apiClient, err := fctl.NewMembershipClient(flags, ctx, cfg, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	protected := !config.GetBool(flags, unprotectFlag)
	metadata := map[string]string{
		fctl.ProtectedStackMetadata: fctl.BoolPointerToString(&protected),
	}

	name := ""
	if len(c.config.GetArgs()) > 0 {
		name = c.config.GetArgs()[0]
	} else {
		name, err = pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show("Enter a name")
		if err != nil {
			return nil, err
		}
	}

	region := config.GetString(flags, regionFlag)
	if region == "" {
		regions, _, err := apiClient.DefaultApi.ListRegions(ctx, organization).Execute()
		if err != nil {
			return nil, errors.Wrap(err, "listing regions")
		}

		var options []string
		for _, region := range regions.Data {
			privacy := "Private"
			if region.Public {
				privacy = "Public "
			}
			name := "<noname>"
			if region.Name != "" {
				name = region.Name
			}
			options = append(options, fmt.Sprintf("%s | %s | %s", region.Id, privacy, name))
		}

		printer := pterm.DefaultInteractiveSelect.WithOptions(options)
		selectedOption, err := printer.Show("Please select a region")
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(options); i++ {
			if selectedOption == options[i] {
				region = regions.Data[i].Id
				break
			}
		}
	}

	stackResponse, _, err := apiClient.DefaultApi.CreateStack(ctx, organization).CreateStackRequest(membershipclient.CreateStackRequest{
		Name:     name,
		Metadata: metadata,
		RegionID: region,
	}).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "creating stack")
	}

	profile := fctl.GetCurrentProfile(flags, cfg)

	if !config.GetBool(flags, nowaitFlag) {
		var spinner *pterm.SpinnerPrinter
		if config.GetString(flags, config.OutputFlag) == "plain" {
			spinner, err = pterm.DefaultSpinner.Start("Waiting services availability")
			if err != nil {
				return nil, err
			}
		}

		if err := waitStackReady(ctx, c.config.GetOut(), flags, profile, stackResponse.Data); err != nil {
			return nil, err
		}

		if spinner != nil {
			if err := spinner.Stop(); err != nil {
				return nil, err
			}
		}
	}

	stackClient, err := fctl.NewStackClient(flags, ctx, cfg, stackResponse.Data, c.config.GetOut())
	if err != nil {
		return nil, err
	}

	versions, err := stackClient.GetVersions(ctx)
	if err != nil {
		return nil, err
	}
	if versions.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when reading versions", versions.StatusCode)
	}

	c.store.Stack = stackResponse.Data
	c.store.Versions = versions.GetVersionsResponse
	c.profile = profile

	return c, nil
}

func (c *CreateController) Render() (tea.Model, error) {

	if config.GetString(c.config.GetAllFLags(), config.OutputFlag) == "plain" {
		fctl.BasicTextCyan.WithWriter(c.config.GetOut()).Printfln("Your dashboard will be reachable on: %s",
			c.profile.ServicesBaseUrl(c.store.Stack).String())
	}

	model, err := internal.PrintStackInformation(c.config.GetOut(), c.config.GetAllFLags(), c.profile, c.store.Stack, c.store.Versions)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func NewCreateCommand() *cobra.Command {
	config := NewStackConfig()
	return fctl.NewCommand(config.GetUse(),
		fctl.WithArgs(cobra.RangeArgs(0, 1)),
		fctl.WithController(NewStackController(config)),
	)
}
