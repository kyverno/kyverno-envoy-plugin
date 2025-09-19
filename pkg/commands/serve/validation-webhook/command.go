package validationwebhook

import (
	"context"
	"fmt"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	apolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/apol/compiler"
	vpolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/vpol/compiler"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/probes"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/signals"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/webhook/validation"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

func Command() *cobra.Command {
	var probesAddress string
	var metricsAddress string
	var kubeConfigOverrides clientcmd.ConfigOverrides
	command := &cobra.Command{
		Use:   "validation-webhook",
		Short: "Start the validation webhook",
		RunE: func(cmd *cobra.Command, args []string) error {
			// setup signals aware context
			return signals.Do(context.Background(), func(ctx context.Context) error {
				// track errors
				var httpErr, grpcErr, mgrErr error
				err := func(ctx context.Context) error {
					// create a rest config
					kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
						clientcmd.NewDefaultClientConfigLoadingRules(),
						&kubeConfigOverrides,
					)
					config, err := kubeConfig.ClientConfig()
					if err != nil {
						return err
					}
					// create a wait group
					var group wait.Group
					// wait all tasks in the group are over
					defer group.Wait()
					// create a controller manager
					scheme := runtime.NewScheme()
					if err := v1alpha1.Install(scheme); err != nil {
						return err
					}
					if err := vpol.Install(scheme); err != nil {
						return err
					}
					mgr, err := ctrl.NewManager(config, ctrl.Options{
						Scheme: scheme,
						Metrics: metricsserver.Options{
							BindAddress: metricsAddress,
						},
						LeaderElection: false,
					})
					if err != nil {
						return fmt.Errorf("failed to construct manager: %w", err)
					}
					// create apol compiler
					apolCompiler := apolcompiler.NewCompiler()
					apolCompileFunc := func(policy *v1alpha1.AuthorizationPolicy) field.ErrorList {
						_, err := apolCompiler.Compile(policy)
						ctrl.LoggerFrom(ctx).Error(err.ToAggregate(), "Authorization policy compilation error")
						return err
					}

					k8sClient, err := dynamic.NewForConfig(config)
					if err != nil {
						return err
					}

					// create vpol compiler
					vpolCompiler := vpolcompiler.NewCompiler(k8sClient)
					vpolCompileFunc := func(policy *vpol.ValidatingPolicy) field.ErrorList {
						_, err := vpolCompiler.Compile(policy)
						ctrl.LoggerFrom(ctx).Error(err.ToAggregate(), "Validating policy compilation error")
						return err
					}
					v := validation.NewValidator(apolCompileFunc, vpolCompileFunc)
					if err := ctrl.NewWebhookManagedBy(mgr).For(&v1alpha1.AuthorizationPolicy{}).WithValidator(v).Complete(); err != nil {
						return fmt.Errorf("failed to create webhook: %w", err)
					}
					if err := ctrl.NewWebhookManagedBy(mgr).For(&vpol.ValidatingPolicy{}).WithValidator(v).Complete(); err != nil {
						return fmt.Errorf("failed to create webhook: %w", err)
					}
					// create a cancellable context
					ctx, cancel := context.WithCancel(ctx)
					// start manager
					group.StartWithContext(ctx, func(ctx context.Context) {
						// cancel context at the end
						defer cancel()
						mgrErr = mgr.Start(ctx)
					})
					// create http and grpc servers
					http := probes.NewServer(probesAddress)
					// run servers
					group.StartWithContext(ctx, func(ctx context.Context) {
						// cancel context at the end
						defer cancel()
						httpErr = http.Run(ctx)
					})
					return nil
				}(ctx)
				return multierr.Combine(err, httpErr, grpcErr, mgrErr)
			})
		},
	}
	command.Flags().StringVar(&probesAddress, "probes-address", ":9080", "Address to listen on for health checks")
	command.Flags().StringVar(&metricsAddress, "metrics-address", ":9082", "Address to listen on for metrics")
	clientcmd.BindOverrideFlags(&kubeConfigOverrides, command.Flags(), clientcmd.RecommendedConfigOverrideFlags("kube-"))
	return command
}
