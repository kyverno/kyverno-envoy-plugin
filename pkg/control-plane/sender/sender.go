package sender

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	"go.uber.org/multierr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	ctrl "sigs.k8s.io/controller-runtime"
)

type PolicySender struct {
	protov1alpha1.UnimplementedValidatingPolicyServiceServer
	polMu                     *sync.Mutex
	policies                  map[string]*protov1alpha1.ValidatingPolicy
	healthCheckMap            map[string]time.Time
	cxnMu                     *sync.Mutex
	cxnsMap                   map[string]grpc.BidiStreamingServer[protov1alpha1.ValidatingPolicyStreamRequest, protov1alpha1.ValidatingPolicy]
	ctx                       context.Context
	initialSendPolicyWait     time.Duration // how long to wait before the second attempt of a failed policy send
	maxSendPolicyInterval     time.Duration // the maximum duration to wait before stopping attempts of a policy send
	clientFlushInterval       time.Duration // how often we remove unhealthy clients from the map
	maxClientInactiveDuration time.Duration // how long should we wait before deciding this client is unhealthy
}

func NewPolicySender(
	ctx context.Context,
	initialSendPolicyWait time.Duration,
	maxSendPolicyInterval time.Duration,
	clientFlushInterval time.Duration,
	maxClientInactiveDuration time.Duration,
) *PolicySender {
	return &PolicySender{
		polMu:                     &sync.Mutex{},
		cxnMu:                     &sync.Mutex{},
		ctx:                       ctx,
		policies:                  make(map[string]*protov1alpha1.ValidatingPolicy),
		healthCheckMap:            make(map[string]time.Time),
		cxnsMap:                   make(map[string]grpc.BidiStreamingServer[protov1alpha1.ValidatingPolicyStreamRequest, protov1alpha1.ValidatingPolicy]),
		initialSendPolicyWait:     initialSendPolicyWait,
		maxSendPolicyInterval:     maxSendPolicyInterval,
		clientFlushInterval:       clientFlushInterval,
		maxClientInactiveDuration: maxClientInactiveDuration,
	}
}

func (s *PolicySender) StartHealthCheckMonitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(s.clientFlushInterval):
			s.deleteInactive()
		}
	}
}

func (s *PolicySender) SendPolicy(pol *protov1alpha1.ValidatingPolicy) {
	errCh := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(s.cxnsMap))
	// send to clients, but don't wait on any of them
	for _, stream := range s.cxnsMap {
		go func() {
			defer wg.Done()
			errCh <- s.sendWithBackoff(stream, pol)
		}()
	}

	wg.Wait()
	close(errCh)

	errs := make([]error, len(errCh))
	for e := range errCh {
		errs = append(errs, e)
	}
	if len(errs) > 0 {
		ctrl.LoggerFrom(nil).Error(multierr.Combine(errs...), "failed to send policy")
	}
}

func (s *PolicySender) StorePolicy(pol *protov1alpha1.ValidatingPolicy) {
	s.polMu.Lock()
	defer s.polMu.Unlock()
	s.policies[pol.Name] = pol
}

func (s *PolicySender) DeletePolicy(polName string) {
	s.polMu.Lock()
	defer s.polMu.Unlock()
	delete(s.policies, polName)
}

func (s *PolicySender) HealthCheck(ctx context.Context, r *protov1alpha1.HealthCheckRequest) (*protov1alpha1.HealthCheckResponse, error) {
	if r.ClientAddress == "" || r.Time == nil {
		return nil, nil // invalid request, do nothing
	}
	// s.logger.Debugf("got health check message from %s, time: %s", r.ClientAddress, r.Time.AsTime().Format(time.RFC3339))
	t, ok := s.healthCheckMap[r.ClientAddress]
	if !ok || r.Time.AsTime().After(t) {
		s.healthCheckMap[r.ClientAddress] = r.Time.AsTime()
	}
	return &protov1alpha1.HealthCheckResponse{}, nil
}

func (s *PolicySender) ValidatingPoliciesStream(stream grpc.BidiStreamingServer[protov1alpha1.ValidatingPolicyStreamRequest, protov1alpha1.ValidatingPolicy]) error {
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			req, err := stream.Recv()
			if err == io.EOF {
				if p, ok := peer.FromContext(stream.Context()); ok {
					ctrl.LoggerFrom(nil).Info(fmt.Sprintf("Receiver at %s closed the stream", p.Addr))
				} else {
					ctrl.LoggerFrom(nil).Info("Receiver closed the stream")
				}
				return nil
			}
			if err != nil {
				if p, ok := peer.FromContext(stream.Context()); ok {
					ctrl.LoggerFrom(nil).Info(fmt.Sprintf("Receiver at %s errored", p.Addr))
				} else {
					ctrl.LoggerFrom(nil).Info("Receiver errored")
				}
				return err
			}
			ctrl.LoggerFrom(nil).Info(fmt.Sprintf("Received validating policy stream request from: %s", req.ClientAddress))
			for _, pol := range s.policies {
				// send each policy in a goroutine to avoid blocking the receive loop
				go func(p *protov1alpha1.ValidatingPolicy) {
					if err := s.sendWithBackoff(stream, p); err != nil {
						ctrl.LoggerFrom(nil).Error(err, "Error sending policy with backoff")
					}
				}(pol)
			}
			s.cxnMu.Lock()
			s.cxnsMap[req.ClientAddress] = stream
			s.cxnMu.Unlock()
		}
	}
}

func (s *PolicySender) sendWithBackoff(stream grpc.BidiStreamingServer[protov1alpha1.ValidatingPolicyStreamRequest, protov1alpha1.ValidatingPolicy], pol *protov1alpha1.ValidatingPolicy) error {
	operation := func() error {
		if err := stream.Send(pol); err != nil {
			return err
		}
		return nil
	}
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = s.initialSendPolicyWait
	b.MaxInterval = s.maxSendPolicyInterval
	return backoff.Retry(operation, b)
}

func (s *PolicySender) deleteInactive() {
	defer s.cxnMu.Unlock()
	inactiveClients := s.getInactiveClients()
	s.cxnMu.Lock()
	for _, c := range inactiveClients {
		ctrl.LoggerFrom(nil).Info(fmt.Sprintf("deleting %s from active clients", c))
		delete(s.cxnsMap, c)
		delete(s.healthCheckMap, c)
	}
}

func (s *PolicySender) getInactiveClients() []string {
	clientsToDelete := []string{}
	for c, t := range s.healthCheckMap {
		if elapsed := time.Since(t); elapsed > s.maxClientInactiveDuration {
			clientsToDelete = append(clientsToDelete, c)
		}
	}
	return clientsToDelete
}
