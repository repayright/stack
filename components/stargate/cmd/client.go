package cmd

import (
	"github.com/formancehq/stack/components/stargate/internal/client"
	"github.com/formancehq/stack/components/stargate/internal/client/interceptors"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	app "github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/grpc/keepalive"
)

const (
	organizationIDFlag = "organization-id"
	stackIDFlag        = "stack-id"

	stargateServerURLFlag = "stargate-server-url"
	gatewayURLFlag        = "gateway-url"

	workerPoolMaxWorkersFlag = "worker-pool-max-worker"
	workerPoolMaxTasksFlag   = "worker-pool-max-tasks"

	ClientChanSizeFlag                = "client-chan-size"
	HTTPClientTimeoutFlag             = "http-client-timeout"
	HTTPClientMaxIdleConnsFlag        = "http-client-max-idle-conns"
	HTTPClientMaxIdleConnsPerHostFlag = "http-client-max-idle-conns-per-host"

	KeepAliveClientParamTimeFlag                = "keepalive-client-param-time"
	KeepAliveClientParamTimeoutFlag             = "keepalive-client-param-timeout"
	KeepAliveClientParamPermitWithoutStreamFlag = "keepalive-client-param-permit-without-stream"

	AuthRefreshTokenDurationBeforeExpireTimeFlag = "auth-refresh-token-duration-before-expire-time"
	AuthClientIDFlag                             = "auth-client-id"
	AuthClientSecretFlag                         = "auth-client-secret"
	AuthEndpointFlag                             = "auth-endpoint"
)

func newClient() *cobra.Command {
	return &cobra.Command{
		Use:          "client",
		Short:        "Launch client",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.New(cmd.OutOrStdout(), resolveClientOptions(
				viper.GetViper(),
			)...).Run(cmd.Context())
		},
	}
}

func resolveClientOptions(v *viper.Viper) []fx.Option {
	options := make([]fx.Option, 0)
	options = append(options, fx.NopLogger)

	options = append(options,
		otlptraces.CLITracesModule(viper.GetViper()),
		otlpmetrics.CLIMetricsModule(viper.GetViper()),

		fx.Provide(func() client.WorkerPoolConfig {
			return client.NewWorkerPoolConfig(
				viper.GetInt(workerPoolMaxWorkersFlag),
				viper.GetInt(workerPoolMaxTasksFlag),
			)
		}),
		fx.Provide(func() client.Config {
			return client.NewClientConfig(
				viper.GetString(organizationIDFlag),
				viper.GetString(stackIDFlag),
				viper.GetInt(ClientChanSizeFlag),
				viper.GetString(gatewayURLFlag),
				viper.GetDuration(HTTPClientTimeoutFlag),
				viper.GetInt(HTTPClientMaxIdleConnsFlag),
				viper.GetInt(HTTPClientMaxIdleConnsPerHostFlag),
			)
		}),

		fx.Provide(func() keepalive.ClientParameters {
			return client.NewKeepAliveClientParams(
				viper.GetDuration(KeepAliveClientParamTimeFlag),
				viper.GetDuration(KeepAliveClientParamTimeoutFlag),
				viper.GetBool(KeepAliveClientParamPermitWithoutStreamFlag),
			)
		}),
		fx.Provide(func() interceptors.Config {
			return interceptors.NewConfig(
				viper.GetDuration(AuthRefreshTokenDurationBeforeExpireTimeFlag),
				viper.GetString(AuthClientIDFlag),
				viper.GetString(AuthClientSecretFlag),
				viper.GetString(AuthEndpointFlag),
			)
		}),
		client.Module(viper.GetString(stargateServerURLFlag)),
	)

	return options
}
