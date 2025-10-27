package listener

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/processor"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	ctrl "sigs.k8s.io/controller-runtime"
)

type policyListener struct {
	controlPlaneAddr            string
	clientAddr                  string
	client                      protov1alpha1.ValidatingPolicyServiceClient
	conn                        *grpc.ClientConn
	processors                  map[vpol.EvaluationMode]processor.Processor
	connEstablished             bool
	controlPlaneReconnectWait   time.Duration
	controlPlaneMaxDialInterval time.Duration
	healthCheckInterval         time.Duration

	// control-plane tracking
	currentVersion string
	currentNonce   string
	mu             sync.Mutex
}

// can two storage entities share the underlying connection of the policy listener?
// store the normal policies and have it message
func NewPolicyListener(
	controlPlaneAddr string,
	clientAddr string,
	processors map[vpol.EvaluationMode]processor.Processor,
	controlPlaneReconnectWait,
	controlPlaneMaxDialInterval,
	healthCheckInterval time.Duration) *policyListener {
	return &policyListener{
		controlPlaneAddr:            controlPlaneAddr,
		processors:                  processors,
		clientAddr:                  clientAddr,
		controlPlaneReconnectWait:   controlPlaneReconnectWait,
		controlPlaneMaxDialInterval: controlPlaneMaxDialInterval,
		healthCheckInterval:         healthCheckInterval,
	}
}

func (l *policyListener) Start(ctx context.Context) error {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = l.controlPlaneReconnectWait
	b.MaxInterval = l.controlPlaneReconnectWait
	if err := backoff.Retry(l.dial, b); err != nil {
		return err
	}
	go l.sendHealthChecks(ctx) // start sending health checks
	if err := l.listen(ctx); err != nil {
		return err
	}
	return nil
}

func (l *policyListener) dial() error {
	ctrl.LoggerFrom(nil).Info(fmt.Sprintf("Connecting to control plane at %s", l.controlPlaneAddr))
	l.connEstablished = false // set connection to false to mark a new connection
	conn, err := grpc.NewClient(l.controlPlaneAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	l.conn = conn
	l.client = protov1alpha1.NewValidatingPolicyServiceClient(conn)
	return nil
}

func (l *policyListener) listen(ctx context.Context) error {
	ctrl.LoggerFrom(nil).Info("Establishing validation channel...")

	// Establish the stream
	stream, err := l.client.PolicyDiscoveryStream(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				ctrl.LoggerFrom(nil).Info("Stopping policy listener due to context cancellation")
				if err := stream.CloseSend(); err != nil {
					ctrl.LoggerFrom(nil).Error(err, "Error closing stream")
				}

				if l.conn != nil {
					if err := l.conn.Close(); err != nil {
						ctrl.LoggerFrom(nil).Error(err, "Error closing connection")
					}
				}
				return
			default:
				// Send initial or ACK/NACK request
				if !l.connEstablished {
					// Initial request with empty version and nonce
					if err := stream.Send(&protov1alpha1.PolicyDiscoveryRequest{
						ClientAddress: l.clientAddr,
						VersionInfo:   "",
						ResponseNonce: "",
					}); err != nil {
						ctrl.LoggerFrom(nil).Error(err, "Error sending initial request")
						return
					}
					l.connEstablished = true
				}

				// Receive policy discovery response
				resp, err := stream.Recv()
				if err == io.EOF {
					ctrl.LoggerFrom(nil).Info("Policy discovery stream closed by server")
					return
				}
				if err != nil {
					ctrl.LoggerFrom(nil).Error(err, "Error receiving policy discovery response")
					return
				}

				ctrl.LoggerFrom(nil).Info(fmt.Sprintf("Received policy discovery response: version=%s, nonce=%s, policies=%d",
					resp.VersionInfo, resp.Nonce, len(resp.Policies)))

				// Process all policies from the response
				applyErr := l.applyPolicies(resp.Policies)

				// Update tracked version and nonce
				l.mu.Lock()
				l.currentVersion = resp.VersionInfo
				l.currentNonce = resp.Nonce
				l.mu.Unlock()

				// Send ACK or NACK
				ackReq := &protov1alpha1.PolicyDiscoveryRequest{
					ClientAddress: l.clientAddr,
					ResponseNonce: resp.Nonce,
				}

				if applyErr != nil {
					// NACK: Send error details
					ctrl.LoggerFrom(nil).Error(applyErr, "Failed to apply policies")
					ackReq.VersionInfo = l.currentVersion // Keep old version
					ackReq.ErrorDetail = &protov1alpha1.ErrorDetail{
						Message: applyErr.Error(),
					}
				} else {
					// ACK: Send new version
					ackReq.VersionInfo = resp.VersionInfo
				}

				if err := stream.Send(ackReq); err != nil {
					ctrl.LoggerFrom(nil).Error(err, "Error sending ACK/NACK")
					return
				}
			}
		}
	}()

	ctrl.LoggerFrom(nil).Info("Policy listener running...")
	wg.Wait()
	return nil
}

// applyPolicies processes a list of policies received from the discovery service
func (l *policyListener) applyPolicies(policies []*protov1alpha1.ValidatingPolicy) error {
	for _, pol := range policies {
		// Determine the evaluation mode
		if pol.Spec == nil {
			ctrl.LoggerFrom(nil).Info(fmt.Sprintf("Policy %s has no spec, skipping", pol.Name))
			continue
		}

		mode := vpol.EvaluationMode(pol.Spec.EvaluationMode)
		processor, ok := l.processors[mode]
		if !ok {
			ctrl.LoggerFrom(nil).Info(fmt.Sprintf("No processor for evaluation mode %s, skipping policy %s", mode, pol.Name))
			continue
		}

		// Process the policy (Process doesn't return an error)
		processor.Process(pol)
		ctrl.LoggerFrom(nil).Info(fmt.Sprintf("Successfully processed policy: %s", pol.Name))
	}
	return nil
}

func (l *policyListener) sendHealthChecks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(l.healthCheckInterval):
			if _, err := l.client.HealthCheck(ctx, &protov1alpha1.HealthCheckRequest{
				ClientAddress: l.clientAddr,
				Time:          timestamppb.Now()}); err != nil {
				ctrl.LoggerFrom(ctx).Error(err, "Health check failed")
			}
			continue
		}
	}
}
