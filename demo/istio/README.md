# Istio Demo 

This Istio Demo is prototype of the kyverno envoy plugin .

## Overview 

The goal of the demo to show user how kyverno-envoy-plugin will work with istio and how it can be used to enforce policies to the traffic between services. The Kyverno-envoy-plugin allows configuring these Envoy proxies to query Kyverno-json for policy decisions on incoming requests.

## Contains
 
- A manifests folder with everything we need to run the demo . 
- bootstrap.sh creates the cluster and installs istio . 

## Architecture
The below architecture illustrates a scenario where no service mesh or Envoy-like components have been pre-installed or already installed.

![Architecture](architecture1.png)


The below architecture illustrates a scenario where a service mesh or Envoy-like components have been pre-installed or already installed.
![Architecture](architecture2.png)

## Requirements

- Istio Authorizationpolicy manifest  to add "extension provider " concept in MeshConfig to specify Where/how to talk to envoy ext-authz service 
-
-