package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

func NewRootCommand() *cobra.Command {
	viper.SetDefault("version", Version)

	root := &cobra.Command{
		Use:               "stargate",
		Short:             "stargate",
		DisableAutoGenTag: true,
	}

	version := newVersion()
	root.AddCommand(version)

	server := newServer()
	root.AddCommand(server)

	client := newClient()
	root.AddCommand(client)

	publish.InitCLIFlags(server)
	server.Flags().String(serviceHttpAddrFlag, "localhost:8080", "Listen address for http API")
	server.Flags().String(serviceGrpcAddrFlag, "localhost:3068", "Listen address for grpc API")
	server.Flags().String(jwksURLFlag, "", "JWKS URL")
	server.Flags().Duration(natsRequestTimeout, 10*time.Second, "NATS request timeout (in seconds)")
	server.Flags().Duration(KeepAlivePolicyMinTimeFlag, 5*time.Second, "Keepalive policy min time")
	server.Flags().Bool(KeepAlivePolicyPermitWithoutStreamFlag, true, "Keepalive policy permit without stream")
	server.Flags().Duration(KeepAliveServerParamMaxConnectionIdleFlag, 15*time.Second, "Keepalive policy permit without stream")
	server.Flags().Duration(KeepAliveServerParamMaxConnectionAgeFlag, 30*time.Second, "Keepalive policy permit without stream")
	server.Flags().Duration(KeepAliveServerParamMaxConnectionAgeGraceFlag, 5*time.Second, "Keepalive policy permit without stream")
	server.Flags().Duration(KeepAliveServerParamTimeFlag, 5*time.Second, "Keepalive policy permit without stream")
	server.Flags().Duration(KeepAliveServerParamTimeoutFlag, 1*time.Second, "Keepalive policy permit without stream")
	if err := viper.BindPFlags(server.Flags()); err != nil {
		panic(err)
	}

	client.Flags().String(organizationIDFlag, "", "Organization ID")
	client.Flags().String(stackIDFlag, "", "Stack ID")
	client.Flags().String(stargateServerURLFlag, "toto", "Stargate server URL")
	client.Flags().String(gatewayURLFlag, "", "Gateway URL")
	client.Flags().Int(workerPoolMaxWorkersFlag, 100, "Max worker pool size")
	client.Flags().Int(workerPoolMaxTasksFlag, 10000, "Max worker pool tasks")
	client.Flags().Int(ClientChanSizeFlag, 1024, "Client chan size")
	client.Flags().Duration(HTTPClientTimeoutFlag, 10*time.Second, "HTTP client timeout")
	client.Flags().Int(HTTPClientMaxIdleConnsFlag, 100, "HTTP client max idle conns")
	client.Flags().Int(HTTPClientMaxIdleConnsPerHostFlag, 2, "HTTP client max idle conns per host")
	client.Flags().Duration(KeepAliveClientParamTimeFlag, 10*time.Second, "Keepalive client param time")
	client.Flags().Duration(KeepAliveClientParamTimeoutFlag, time.Second, "Keepalive client param timeout")
	client.Flags().Bool(KeepAliveClientParamPermitWithoutStreamFlag, true, "Keepalive client param permit without stream")
	client.Flags().Duration(AuthRefreshTokenDurationBeforeExpireTimeFlag, 30*time.Second, "Auth refresh token duration")
	client.Flags().String(AuthClientIDFlag, "", "Auth client ID")
	client.Flags().String(AuthClientSecretFlag, "", "Auth client secret")
	client.Flags().String(AuthEndpointFlag, "", "Auth URL")
	if err := viper.BindPFlags(client.Flags()); err != nil {
		panic(err)
	}

	root.PersistentFlags().Bool(service.DebugFlag, false, "Debug mode")

	otlptraces.InitOTLPTracesFlags(root.Flags())
	otlpmetrics.InitOTLPMetricsFlags(root.Flags())

	if err := viper.BindPFlags(root.PersistentFlags()); err != nil {
		panic(err)
	}

	return root
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		if _, err = fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}

		os.Exit(1)
	}
}
