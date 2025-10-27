package discovery

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"sort"
	"sync"
	"time"

	protov1alpha1 "github.com/kyverno/kyverno-envoy-plugin/pkg/control-plane/proto/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type PolicyDiscoveryService struct {
	protov1alpha1.UnimplementedValidatingPolicyServiceServer

	// Policy storage
	polMu    *sync.RWMutex
	policies map[string]*protov1alpha1.ValidatingPolicy

	// Connection management
	cxnMu       *sync.RWMutex
	connections map[string]*clientConnection

	// Version tracking
	versionMu      *sync.RWMutex
	currentVersion string // Hash of all stored policies
	currentNonce   string // Current nonce for this version
	pendingVersion *pendingVersionState

	// Health checks
	healthMu        *sync.RWMutex
	lastHealthCheck map[string]time.Time

	// Configuration
	ctx                 context.Context
	healthCheckTimeout  time.Duration
	clientPruneInterval time.Duration
	ackTimeout          time.Duration
}

type clientConnection struct {
	stream         grpc.BidiStreamingServer[protov1alpha1.PolicyDiscoveryRequest, protov1alpha1.PolicyDiscoveryResponse]
	clientAddr     string
	currentVersion string
	mu             *sync.Mutex
}

type pendingVersionState struct {
	version      string
	nonce        string
	ackedClients map[string]bool
	mu           *sync.Mutex
	allAcked     chan struct{}
}

// NewPolicyDiscoveryService creates a new PolicyDiscoveryService instance
func NewPolicyDiscoveryService(
	ctx context.Context,
	healthCheckTimeout time.Duration,
	clientPruneInterval time.Duration,
	ackTimeout time.Duration,
) *PolicyDiscoveryService {
	return &PolicyDiscoveryService{
		UnimplementedValidatingPolicyServiceServer: protov1alpha1.UnimplementedValidatingPolicyServiceServer{},
		polMu:               &sync.RWMutex{},
		policies:            make(map[string]*protov1alpha1.ValidatingPolicy),
		cxnMu:               &sync.RWMutex{},
		connections:         make(map[string]*clientConnection),
		versionMu:           &sync.RWMutex{},
		currentVersion:      "",
		currentNonce:        "",
		pendingVersion:      nil,
		healthMu:            &sync.RWMutex{},
		lastHealthCheck:     make(map[string]time.Time),
		ctx:                 ctx,
		healthCheckTimeout:  healthCheckTimeout,
		clientPruneInterval: clientPruneInterval,
		ackTimeout:          ackTimeout,
	}
}

func (s *PolicyDiscoveryService) PolicyDiscoveryStream(stream grpc.BidiStreamingServer[protov1alpha1.PolicyDiscoveryRequest, protov1alpha1.PolicyDiscoveryResponse]) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	clientAddr := req.GetClientAddress()
	if clientAddr == "" {
		return status.Errorf(codes.InvalidArgument, "client address required")
	}

	cxn := &clientConnection{
		stream:         stream,
		clientAddr:     clientAddr,
		currentVersion: "",
		mu:             &sync.Mutex{},
	}

	// Register client
	s.cxnMu.Lock()
	s.connections[clientAddr] = cxn
	s.cxnMu.Unlock()

	// Send initial policy snapshot
	s.polMu.RLock()
	policies := make([]*protov1alpha1.ValidatingPolicy, 0, len(s.policies))
	for _, pol := range s.policies {
		policies = append(policies, pol)
	}
	s.polMu.RUnlock()

	s.versionMu.RLock()
	initialVersion := s.currentVersion
	initialNonce := s.currentNonce
	s.versionMu.RUnlock()

	resp := &protov1alpha1.PolicyDiscoveryResponse{
		VersionInfo: initialVersion,
		Policies:    policies,
		Nonce:       initialNonce,
	}
	if sendErr := stream.Send(resp); sendErr != nil {
		// Clean up registration
		s.cxnMu.Lock()
		delete(s.connections, clientAddr)
		s.cxnMu.Unlock()
		return sendErr
	}
	cxn.currentVersion = initialVersion

	// Main loop
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			// On disconnect with error
			break
		}

		// Process the discovery request and get response if needed
		resp, err := s.processDiscoveryRequest(clientAddr, req)
		if err != nil {
			log.Printf("Error processing discovery request from %s: %v", clientAddr, err)
			break
		}

		// If a response is needed, send it
		if resp != nil {
			// Update tracked version
			cxn.currentVersion = resp.VersionInfo

			if err := stream.Send(resp); err != nil {
				log.Printf("Error sending response to %s: %v", clientAddr, err)
				break
			}
		}
	}

	// On disconnect
	s.cxnMu.Lock()
	delete(s.connections, clientAddr)
	s.cxnMu.Unlock()

	s.versionMu.Lock()
	if s.pendingVersion != nil {
		delete(s.pendingVersion.ackedClients, clientAddr)
	}
	s.versionMu.Unlock()

	return nil
}

func (s *PolicyDiscoveryService) StorePolicy(pol *protov1alpha1.ValidatingPolicy) error {
	s.polMu.Lock()
	defer s.polMu.Unlock()
	s.policies[pol.Name] = pol

	// Compute new hash of all policies and update version/nonce.
	s.versionMu.Lock()
	defer s.versionMu.Unlock()

	// Collect all policies in a deterministic order
	policyNames := make([]string, 0, len(s.policies))
	for name := range s.policies {
		policyNames = append(policyNames, name)
	}
	sort.Strings(policyNames)

	h := sha256.New()
	for _, name := range policyNames {
		pol := s.policies[name]
		// Assume protov1alpha1.ValidatingPolicy implements proto.Message
		b, err := proto.Marshal(pol)
		if err != nil {
			continue // skip problematic policies (should not happen)
		}
		h.Write([]byte(name))
		h.Write(b)
	}
	newVersion := fmt.Sprintf("%x", h.Sum(nil))

	// If version changed, update version and generate new nonce
	if s.currentVersion != newVersion {
		s.currentVersion = newVersion
		s.currentNonce = fmt.Sprintf("%x", time.Now().UnixNano())
		// reset pendingVersion state
		s.pendingVersion = &pendingVersionState{
			version:      newVersion,
			nonce:        s.currentNonce,
			ackedClients: make(map[string]bool),
			mu:           &sync.Mutex{},
			allAcked:     make(chan struct{}),
		}
	}

	// Send update to all connected clients
	s.cxnMu.RLock()
	defer s.cxnMu.RUnlock()
	for _, cxn := range s.connections {
		go func(c *clientConnection) {
			c.mu.Lock()
			defer c.mu.Unlock()
			// The client will receive the update in the main loop of PolicyDiscoveryStream
		}(cxn)
	}
	return nil
}

func (s *PolicyDiscoveryService) DeletePolicy(polName string) error {
	s.polMu.Lock()
	defer s.polMu.Unlock()
	delete(s.policies, polName)

	// Recompute version hash after deletion
	s.versionMu.Lock()
	defer s.versionMu.Unlock()

	// Collect all policies in a deterministic order
	policyNames := make([]string, 0, len(s.policies))
	for name := range s.policies {
		policyNames = append(policyNames, name)
	}
	sort.Strings(policyNames)

	h := sha256.New()
	for _, name := range policyNames {
		pol := s.policies[name]
		b, err := proto.Marshal(pol)
		if err != nil {
			continue
		}
		h.Write([]byte(name))
		h.Write(b)
	}
	newVersion := fmt.Sprintf("%x", h.Sum(nil))

	// If version changed, update version and generate new nonce
	if s.currentVersion != newVersion {
		s.currentVersion = newVersion
		s.currentNonce = fmt.Sprintf("%x", time.Now().UnixNano())
		// reset pendingVersion state
		s.pendingVersion = &pendingVersionState{
			version:      newVersion,
			nonce:        s.currentNonce,
			ackedClients: make(map[string]bool),
			mu:           &sync.Mutex{},
			allAcked:     make(chan struct{}),
		}
	}

	return nil
}

func (s *PolicyDiscoveryService) GetPolicies() ([]*protov1alpha1.ValidatingPolicy, error) {
	s.polMu.RLock()
	defer s.polMu.RUnlock()
	policies := make([]*protov1alpha1.ValidatingPolicy, 0, len(s.policies))
	for _, pol := range s.policies {
		policies = append(policies, pol)
	}
	return policies, nil
}

// LoadInitialPolicies loads a list of policies into the discovery service.
// This should be called at startup before clients connect to ensure they receive
// the complete policy set. If policies already exist, they will be replaced with
// the new set and version tracking will be updated.
func (s *PolicyDiscoveryService) LoadInitialPolicies(policies []*protov1alpha1.ValidatingPolicy) error {
	s.polMu.Lock()
	defer s.polMu.Unlock()

	// Replace all existing policies with the new set
	s.policies = make(map[string]*protov1alpha1.ValidatingPolicy)
	for _, pol := range policies {
		s.policies[pol.Name] = pol
	}

	// Update version tracking
	s.versionMu.Lock()
	defer s.versionMu.Unlock()

	// Collect all policies in a deterministic order
	policyNames := make([]string, 0, len(s.policies))
	for name := range s.policies {
		policyNames = append(policyNames, name)
	}
	sort.Strings(policyNames)

	// Compute hash of all policies
	h := sha256.New()
	for _, name := range policyNames {
		pol := s.policies[name]
		b, err := proto.Marshal(pol)
		if err != nil {
			continue // skip problematic policies (should not happen)
		}
		h.Write([]byte(name))
		h.Write(b)
	}
	newVersion := fmt.Sprintf("%x", h.Sum(nil))

	// Update version and generate new nonce
	s.currentVersion = newVersion
	s.currentNonce = fmt.Sprintf("%x", time.Now().UnixNano())

	// Reset pending version state
	s.pendingVersion = &pendingVersionState{
		version:      newVersion,
		nonce:        s.currentNonce,
		ackedClients: make(map[string]bool),
		mu:           &sync.Mutex{},
		allAcked:     make(chan struct{}),
	}

	return nil
}

func (s *PolicyDiscoveryService) processDiscoveryRequest(clientAddr string, req *protov1alpha1.PolicyDiscoveryRequest) (*protov1alpha1.PolicyDiscoveryResponse, error) {
	s.versionMu.RLock()
	currentVersion := s.currentVersion
	currentNonce := s.currentNonce
	pv := s.pendingVersion
	s.versionMu.RUnlock()

	responseNonce := req.GetResponseNonce()
	versionInfo := req.GetVersionInfo()
	errorDetail := req.GetErrorDetail()

	// ACK: Client has applied current version
	if responseNonce == currentNonce && versionInfo == currentVersion {
		if pv != nil {
			pv.mu.Lock()
			defer pv.mu.Unlock()
			if !pv.ackedClients[clientAddr] {
				pv.ackedClients[clientAddr] = true
				// Check if all connected clients acked
				allAcked := true
				s.cxnMu.RLock()
				for addr := range s.connections {
					if !pv.ackedClients[addr] {
						allAcked = false
						break
					}
				}
				s.cxnMu.RUnlock()
				if allAcked {
					select {
					case <-pv.allAcked:
						// Already closed
					default:
						close(pv.allAcked)
					}
				}
			}
		}
		// No update to send; return nil means no response necessary
		return nil, nil
	}

	// NACK: Client failed with error on the last update
	if responseNonce == currentNonce && versionInfo != currentVersion && errorDetail != nil {
		log.Printf("[NACK] client %s: error: %s (wanted %s got %s nonce %s)", clientAddr, errorDetail.GetMessage(), currentVersion, versionInfo, currentNonce)
		// Treat as failed, do not wait for this client
		if pv != nil {
			pv.mu.Lock()
			defer pv.mu.Unlock()
			// Mark as acked to not block waiting for this client
			pv.ackedClients[clientAddr] = true
			// If all clients responded (including NACK), close allAcked
			allAcked := true
			s.cxnMu.RLock()
			for addr := range s.connections {
				if !pv.ackedClients[addr] {
					allAcked = false
					break
				}
			}
			s.cxnMu.RUnlock()
			if allAcked {
				select {
				case <-pv.allAcked:
					// Already closed
				default:
					close(pv.allAcked)
				}
			}
		}

		// On NACK, do nothing
		return nil, nil
	}

	// client missed an update: nonce doesn't match current nonce
	if responseNonce != "" && responseNonce != currentNonce {
		s.polMu.RLock()
		policies := make([]*protov1alpha1.ValidatingPolicy, 0, len(s.policies))
		for _, pol := range s.policies {
			policies = append(policies, pol)
		}
		s.polMu.RUnlock()

		resp := &protov1alpha1.PolicyDiscoveryResponse{
			VersionInfo: currentVersion,
			Policies:    policies,
			Nonce:       currentNonce,
		}
		return resp, nil
	}

	// Initial request (connected, no nonce or version)
	if responseNonce == "" && versionInfo == "" {
		s.polMu.RLock()
		policies := make([]*protov1alpha1.ValidatingPolicy, 0, len(s.policies))
		for _, pol := range s.policies {
			policies = append(policies, pol)
		}
		s.polMu.RUnlock()

		resp := &protov1alpha1.PolicyDiscoveryResponse{
			VersionInfo: currentVersion,
			Policies:    policies,
			Nonce:       currentNonce,
		}
		return resp, nil
	}

	// Default: no action necessary
	return nil, nil
}

func (s *PolicyDiscoveryService) HealthCheck(ctx context.Context, req *protov1alpha1.HealthCheckRequest) (*protov1alpha1.HealthCheckResponse, error) {
	s.healthMu.Lock()
	defer s.healthMu.Unlock()
	s.lastHealthCheck[req.GetClientAddress()] = time.Now()
	return &protov1alpha1.HealthCheckResponse{}, nil
}

// FlushInactive removes inactive clients from the health check map
func (s *PolicyDiscoveryService) FlushInactive() {
	s.healthMu.Lock()
	defer s.healthMu.Unlock()
	now := time.Now()
	for addr, lastActive := range s.lastHealthCheck {
		if now.Sub(lastActive) > s.healthCheckTimeout {
			delete(s.lastHealthCheck, addr)
		}
	}
}
