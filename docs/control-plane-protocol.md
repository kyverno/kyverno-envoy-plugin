# Kyverno Envoy Plugin XDS-Like Protocol

## Overview

Bidirectional gRPC streaming protocol that distributes ValidatingPolicies from a control plane to multiple Envoy authorization plugin instances (sidecars). Similar to Envoy's xDS, but customized for ValidatingPolicy distribution.

## Core Concepts

### Versioning Model
- **Version**: SHA256 hash of all policies sorted by name (deterministic)
- **Nonce**: Timestamp-based unique identifier for each discovery response
- **PendingVersion**: Tracks ACK/NACK responses from all connected clients per version

### State Management

**Control Plane State:**
- `policies` - Map storing all ValidatingPolicy objects
- `currentVersion` - SHA256 hash of all policies
- `currentNonce` - Timestamp-based identifier
- `pendingVersion` - Tracks which clients have ACKed the current version
- `connections` - Map of active client streams

**Client State:**
- `currentVersion` - Last successfully applied version
- `currentNonce` - Last received nonce
- `connEstablished` - Connection state flag

## Protocol Messages

### PolicyDiscoveryRequest (Client → Control Plane)
- `client_address` - Unique client identifier
- `version_info` - Last successfully applied version
- `response_nonce` - Nonce from last received response
- `error_detail` - Error details (for NACK)

### PolicyDiscoveryResponse (Control Plane → Client)
- `version_info` - Current version hash
- `policies` - Full array of ValidatingPolicy objects
- `nonce` - Unique nonce for this response

## Request Processing Logic

### 1. ACK - Successful Application
**Condition:** `responseNonce == currentNonce && versionInfo == currentVersion`

**Behavior:**
- Record ACK in pendingVersion tracker
- Check if all connected clients have ACKed
- If all clients responded, close allAcked channel
- Return nil (no response sent)

**Edge Cases:**
- No pending version state → silently ignored
- Duplicate ACKs → tracker prevents double counting
- Disconnected clients still in tracker → counted until cleanup

### 2. NACK - Failed Application
**Condition:** `responseNonce == currentNonce && versionInfo != currentVersion && errorDetail exists`

**Behavior:**
- Log NACK with error details
- Mark client as "acked" in tracker (to unblock waiting)
- Close allAcked channel if all clients responded (including NACKs)
- Return nil (no retry sent)

**Edge Cases:**
- No automatic retry - control plane continues with current version
- Only affects that specific client
- Other clients unaffected

### 3. Stale Nonce - Missed Update
**Condition:** `responseNonce != "" && responseNonce != currentNonce`

**Behavior:**
- Detect client missed an update
- Send full current snapshot immediately

**Triggers:**
- Control plane restarted → new nonce generated
- Client missed update due to network issues
- Client restarted and reconnected

### 4. Initial Connection - Empty State
**Condition:** `responseNonce == "" && versionInfo == ""`

**Behavior:**
- Send full policy snapshot

**Triggers:**
- Client first connects
- Client pod restarts (state lost)
- Control plane restarted (EOF → client reconnects)

### 5. No Update Required
**Condition:** All other cases

**Behavior:**
- Return nil (no response)

**Examples:**
- Client already has current version
- Connection establishing
- Unmatched version/nonce combinations

## Update Flow

```
Kubernetes CRD Change
    ↓
Controller Reconcile
    ↓
StorePolicy() / DeletePolicy()
    ↓
Compute SHA256 Hash (sorted by policy name)
    ↓
Generate New Nonce
    ↓
Create PendingVersion State
    ↓
Wait for Client Requests (Lazy Notification)
    ↓
Send Full Policy Snapshot
    ↓
Client Applies via Processor
    ↓
Client Sends ACK or NACK
    ↓
Control Plane Tracks Responses
    ↓
All Clients Responded → Close allAcked Channel
```

## Policy Operations

### Adding/Updating a Policy
1. Controller detects ValidatingPolicy change
2. Converts to proto format
3. Calls StorePolicy()
4. Compute new SHA256 hash of all policies
5. Generate new nonce (timestamp)
6. Create new pendingVersion state
7. Clients receive update on next request
8. Clients apply and ACK/NACK

### Deleting a Policy
- Same flow as update
- Policy removed from map before hash computation
- Clients receive snapshot without deleted policy

### Concurrent Updates
- Latest version overwrites pendingVersion
- Clients may ACK intermediate versions
- ACKs for stale nonces are ignored (request-based matching)

## Connection Lifecycle

### Client Connection
1. Client establishes gRPC stream
2. Sends initial request with empty version/nonce
3. Control plane registers stream in connections map
4. Control plane sends initial snapshot
5. Client ACKs with version and nonce

### Client Disconnection
1. EOF detected in Recv loop
2. Log disconnect event
3. Remove from connections map
4. Remove from ackedClients tracker
5. Return from listen loop (triggers pod restart)

### Control Plane Restart
1. Clients detect EOF in Recv
2. Clients exit listen loop
3. Pod restart mechanism
4. New connection with empty state
5. Client requests full snapshot

## Health Checks

Separate HealthCheck RPC for liveness tracking:
- Periodic health check updates lastActive timestamp
- Control plane tracks lastHealthCheck per client
- Background task flushes inactive clients (maxClientInactiveDuration)
- Health check failures don't terminate stream

## Edge Cases & Scenarios

| Scenario | Behavior |
|----------|----------|
| Network partition | Client sends request with old nonce → control plane resends snapshot |
| Client restarts during update | New connection requests full snapshot, may ACK with old state |
| Control plane restarts during update | Clients detect EOF, reconnect, request snapshot, get latest version |
| Rapid sequential updates | Latest version replaces earlier pendingVersion, earlier ACKs ignored |
| Empty policy set | Control plane sends empty snapshot array |
| Stream errors | Logged, stream closed, client pod restarts to reconnect |
| Client NACKs update | Logged, marked as done, no retry, control plane continues |
| Early connection before policies loaded | Client may receive empty initial snapshot |
| Concurrent policy version changes | Latest version dominates in control plane |
| Pod restart after partial ACK | New connection triggers full resync |
| Leader election | Control plane may restart, clients reconnect |

## Implementation Details

### Version Computation Algorithm
```go
// Sort policy names deterministically
policyNames := sort(all policy names)

// Compute SHA256 hash
h := sha256.New()
for name in policyNames {
    pol := policies[name]
    marshal(pol)
    h.Write(name)
    h.Write(marshaled_policy)
}
version := h.Sum(nil)
```

### Lazy Notification Pattern
- Control plane waits for client requests
- Updates sent when client calls Recv()
- No active push - client-driven updates

### ACK Tracking
- pendingVersion.ackedClients tracks which clients have responded
- Includes both successful ACKs and NACKs
- allAcked channel closes when all connected clients responded

### NACK Treatment
- Treated as completion (not retry)
- Marks client in tracker as "done"
- Allows update to proceed if all clients responded

## Design Decisions

- **SHA256 for versioning** - Deterministic, consistent across restarts
- **Nonce-based verification** - Request/response correlation for reliability
- **Full snapshot updates** - Handles network partitions and lost messages
- **Lazy notification** - Client-driven prevents blocking
- **ACK tracking** - Can detect when updates propagate (future use: wait for all ACKs)
- **NACK equals completion** - Failed clients don't block others
