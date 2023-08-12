package fctl

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/formancehq/fctl/membershipclient"
	"github.com/formancehq/fctl/pkg/config"
	"github.com/formancehq/formance-sdk-go"
)

func NewMembershipClient(flags *flag.FlagSet, ctx context.Context, cfg *config.Config, out io.Writer) (*membershipclient.APIClient, error) {
	profile := config.GetCurrentProfile(flags, cfg)

	httpClient := GetHttpClient(flags, map[string][]string{}, out)

	configuration := membershipclient.NewConfiguration()

	token, err := profile.GetToken(ctx, httpClient)
	if err != nil {
		return nil, err
	}

	configuration.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	configuration.HTTPClient = httpClient
	configuration.UserAgent = "fctl/" + config.Version
	configuration.Servers[0].URL = profile.GetMembershipURI()

	return membershipclient.NewAPIClient(configuration), nil
}

func NewStackClient(flags *flag.FlagSet, ctx context.Context, cfg *config.Config, stack *membershipclient.Stack, out io.Writer) (*formance.Formance, error) {
	profile := config.GetCurrentProfile(flags, cfg)
	httpClient := GetHttpClient(flags, map[string][]string{}, out)

	token, err := profile.GetStackToken(ctx, httpClient, stack)
	if err != nil {
		return nil, err
	}

	return formance.New(
		formance.WithServerURL(stack.Uri),
		formance.WithClient(
			GetHttpClient(flags, map[string][]string{
				"Authorization": {fmt.Sprintf("Bearer %s", token)},
				"User-Agent":    {"fctl/" + config.Version},
			},
				out,
			),
		),
	), nil
}
