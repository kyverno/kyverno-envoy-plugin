# Hello World

The **Hello World** examples provide the simplest possible introduction to the **Kyverno Authorization Server**.
It’s designed to help you understand the basic workflow, configuration, and behavior of the server—without requiring Kubernetes or any external setup.

- [Envoy Authz Server](./envoy.md)
- [HTTP Authz Server](./http.md)

## Overview

Those example demonstrates how to:

- Start the Kyverno Authz Server locally.
- Load authorization policies from files.
- Send authorization requests.
- Observe how the server evaluates and responds.

!!! note
    This tutorial runs entirely on your local machine.  
    No Kubernetes cluster or other infrastructure is needed.

## Objectives

By the end of this example, you will be able to:

- Run the Authz Server using a local policy directory.
- Craft and send test authorization requests.
- Understand how the server returns `allowed` or `denied` results.
