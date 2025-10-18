package listener

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/processor"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/proto/validatingpolicy/v1alpha1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type policyListener struct {
	controlPlaneAddr            string
	clientAddr                  string
	client                      protov1alpha1.ValidatingPolicyServiceClient
	conn                        *grpc.ClientConn
	processors                  []processor.Processor
	connEstablished             bool
	controlPlaneReconnectWait   time.Duration
	controlPlaneMaxDialInterval time.Duration
	healthCheckInterval         time.Duration
	logger                      *logrus.Logger
}

// can two storage entities share the underlying connection of the policy listener?
// store the normal policies and have it message
func NewPolicyListener(
	controlPlaneAddr string,
	clientAddr string,
	processors []processor.Processor,
	logger *logrus.Logger,
	controlPlaneReconnectWait,
	controlPlaneMaxDialInterval,
	healthCheckInterval time.Duration) *policyListener {
	return &policyListener{
		controlPlaneAddr:            controlPlaneAddr,
		logger:                      logger,
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
	l.logger.Infof("Connecting to control plane at %s", l.controlPlaneAddr)
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
	l.logger.Info("Establishing validation channel...")

	// Establish the stream
	stream, err := l.client.ValidatingPoliciesStream(ctx)
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
				l.logger.Info("Stopping policy listener due to context cancellation")
				if err := stream.CloseSend(); err != nil {
					l.logger.Errorf("Error closing stream: %v", err)
				}

				if l.conn != nil {
					if err := l.conn.Close(); err != nil {
						l.logger.Errorf("Error closing connection: %v", err)
					}
				}
				return
			default:
				if !l.connEstablished {
					if err := stream.Send(&protov1alpha1.ValidatingPolicyStreamRequest{ClientAddress: l.clientAddr}); err != nil {
						l.logger.Error("Error sending to stream")
						return
					}
					l.connEstablished = true
				}
				req, err := stream.Recv()
				if err == io.EOF {
					l.logger.Errorf("Policy sender closed the stream")
					return
				}
				if err != nil {
					l.logger.Errorf("Error receiving policy request: %v", err)
					return
				}

				l.logger.Infof("Received validating policy request: %s, Delete: %t", req.Name, req.Delete)
				go func() {
					for _, p := range l.processors {
						p.Process(req)
					}
				}()
			}
		}
	}()

	l.logger.Info("Policy listener running...")
	wg.Wait()
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
				l.logger.Debugf("Health check failed: %v", err)
			}
			continue
		}
	}
}
