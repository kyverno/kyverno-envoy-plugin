package sidecarserver

import (
	"context"
	"fmt"
	"os"
	"time"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/authz"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/http"
	vpolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/compiler"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/processor"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/httpauth"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/probes"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/stream/listener"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/util/wait"
)

type startFunc = func(context.Context) error

func Command() *cobra.Command {
	var probesAddress string
	var httpAuthAddress string
	var controlPlaneAddr string
	var controlPlaneReconnectWait time.Duration
	var controlPlaneMaxDialInterval time.Duration
	var healthCheckInterval time.Duration
	var nestedRequest bool
	var grpcAddress string
	var grpcNetwork string
	command := &cobra.Command{
		Use:   "sidecar-authz-server",
		Short: "Start the Kyverno Authz Server as a sidecar",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// track errors
				var probesErr, connErr, grpcAuthErr, mgrErr, httpAuthErr error
				err := func(ctx context.Context) error {
					logger := logrus.New()
					// create a cancellable context
					ctx, cancel := context.WithCancel(ctx)
					// cancel context at the end
					defer cancel()
					// create a wait group
					var group wait.Group
					// wait all tasks in the group are over
					defer group.Wait()

					clientAddr := os.Getenv("POD_IP")
					if clientAddr == "" {
						return fmt.Errorf("can't start auth server, no POD_IP has been passed")
					}

					cfg, err := rest.InClusterConfig()
					if err != nil {
						return err
					}

					// initialize kubernetes client
					dyn, err := dynamic.NewForConfig(cfg)
					if err != nil {
						return err
					}

					envoyCompiler := vpolcompiler.NewCompiler[dynamic.Interface, *authv3.CheckRequest, *authv3.CheckResponse]()
					httpCompiler := vpolcompiler.NewCompiler[dynamic.Interface, *http.Request, *http.Response]()
					httpAccessor := processor.NewPolicyAccessor(httpCompiler, logger)
					envoyAccessor := processor.NewPolicyAccessor(envoyCompiler, logger)
					processors := []processor.Processor{httpAccessor, envoyAccessor}

					provider := listener.NewPolicyListener(controlPlaneAddr,
						clientAddr, processors,
						logger, controlPlaneReconnectWait,
						controlPlaneMaxDialInterval,
						healthCheckInterval)

					// create http and grpc server
					probesServer := probes.NewServer(probesAddress)
					authorizer := httpauth.NewAuthorizer(dyn, httpAccessor, nestedRequest, logger)
					httpAuthServer := httpauth.NewServer(httpAuthAddress, authorizer)
					grpc := authz.NewServer(grpcNetwork, grpcAddress, envoyAccessor, dyn, nil)

					// run servers
					group.StartWithContext(ctx, func(ctx context.Context) {
						// probes
						defer cancel()
						probesErr = probesServer.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						// auth server
						defer cancel()
						httpAuthErr = httpAuthServer.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						// grpc auth server
						defer cancel()
						grpcAuthErr = grpc.Run(ctx)
					})
					group.StartWithContext(ctx, func(ctx context.Context) {
						// control plane connection
						for {
							select {
							case <-ctx.Done():
								return
							default:
								if connErr = provider.Start(ctx); connErr != nil {
									logger.Error("error connecting to the control plane, sleeping 10 seconds then retrying")
									time.Sleep(time.Second * 10)
								}
								continue
							}
						}
					})

					return nil
				}(ctx)
				return multierr.Combine(err, probesErr, connErr, grpcAuthErr, mgrErr, httpAuthErr)
			})
		},
	}
	command.Flags().StringVar(&grpcAddress, "grpc-address", ":9081", "Address to listen on")
	command.Flags().StringVar(&grpcNetwork, "grpc-network", "tcp", "Network to listen on")
	command.Flags().DurationVar(&controlPlaneReconnectWait, "control-plane-reconnect-wait", 3*time.Second, "Duration to wait before retrying connecting to the control plane")
	command.Flags().DurationVar(&controlPlaneMaxDialInterval, "control-plane-max-dial-interval", 8*time.Second, "Duration to wait before stopping attempts of sending a policy to a client")
	command.Flags().DurationVar(&healthCheckInterval, "health-check-interval", 30*time.Second, "Interval for sending health checks")
	command.Flags().StringVar(&probesAddress, "probes-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&httpAuthAddress, "http-auth-server-address", ":9083", "Address to serve the http authorization server on")
	command.Flags().StringVar(&controlPlaneAddr, "control-plane-address", "", "Control plane address")
	command.Flags().BoolVar(&nestedRequest, "nested-request", false, "Expect the requests to validate to be in the body of the original request")

	_ = command.MarkFlagRequired("control-plane-address")
	return command
}
