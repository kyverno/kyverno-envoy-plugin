# Quick start

The Kyverno Envoy Plugin is a powerful tool that integrates with the Envoy proxy.

It allows you to enforce Kyverno policies on incoming and outgoing traffic in a service mesh environment, providing an additional layer of security and control over your applications.

## Overview 

[Envoy](https://www.envoyproxy.io/docs/envoy/latest/intro/what_is_envoy) is a Layer 7 proxy and communication bus tailored for large-scale, modern service-oriented architectures. Starting from version 1.7.0, Envoy includes an [External Authorization filter](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html) that interfaces with an authorization service to determine the legitimacy of incoming requests.

This functionality allows authorization decisions to be offloaded to an external service, which can access the request context. The request context includes details such as the origin and destination of the network activity, as well as specifics of the network request (e.g., HTTP request). This information enables the external service to make a well-informed decision regarding the authorization of the incoming request processed by Envoy.

## What is the Kyverno Envoy Plugin?

The [Kyverno Envoy Plugin](https://github.com/kyverno/kyverno-envoy-plugin) is gRPC server that implements [Envoy External Authorization API](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html).

This allows you to enforce Kyverno policies on incoming and outgoing traffic in a service mesh environment, providing an additional layer of security and control over your applications. You can use this version of Kyverno to enforce fine-grained, context-aware access control policies with Envoy without modifying your microservice.

## How does this work?

In addition to the Envoy sidecar, your application pods will include a Kyverno Authz Server component, either as a sidecar or as a separate pod. When Envoy receives an API request intended for your microservice, it consults the Kyverno Authz Server to determine whether the request should be permitted or not.

Performing policy evaluations locally with Envoy is advantageous, as it eliminates the need for an additional network hop for authorization checks, thus enhancing both performance and availability.

!!! info 

    The Kyverno Envoy Plugin is frequently deployed in Kubernetes environments as a sidecar container or as a separate pod. Additionally, it can be used in other environments as a standalone process running alongside Envoy.

## Additional Resources 

See the following pages on [envoyproxy.io](https://www.envoyproxy.io/) for more information on external authorization:

- [External Authorization](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html) to learn about the External Authorization filter.
- [Network](https://www.envoyproxy.io/docs/envoy/latest/configuration/listeners/network_filters/ext_authz_filter#config-network-filters-ext-authz) and [HTTP](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter#config-http-filters-ext-authz) for details on configuring the External Authorization filter.

