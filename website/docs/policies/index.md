# Policies

The Kyverno Authz Server uses `ValidatingPolicy` resources to define authorization rules. Policies can operate in two modes:

## Policy Guides

- **[Envoy Policy Breakdown](./envoy-policy-breakdown.md)** - Complete guide for writing policies that integrate with Envoy proxy
- **[HTTP Policy Breakdown](./http-policy-breakdown.md)** - Complete guide for writing policies for plain HTTP authorization

## Overview

A `ValidatingPolicy` is a Kubernetes custom resource that uses CEL (Common Expression Language) to evaluate authorization requests. The policy's evaluation mode determines whether it processes Envoy CheckRequests or plain HTTP requests.

### Key Concepts

- **Evaluation Mode**: Set to `Envoy` or `HTTP` to determine the request type
- **Failure Policy**: Controls behavior when policy evaluation fails (`Fail` or `Ignore`)
- **Match Conditions**: Optional CEL expressions for fine-grained request filtering
- **Variables**: Reusable named expressions available throughout the policy
- **Validation Rules**: CEL expressions that return authorization decisions
