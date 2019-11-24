# Leader Elector

A simple leader election sidecar container for Kubernetes based on the recent [client-go](https://github.com/kubernetes/client-go) library. It aims to provide and alternative to the widely used but outdated `election` component of [kubernetes/contrib](https://github.com/kubernetes-retired/contrib).

## Work in Progress
* Configuration via flags
* Configurable lock type (LeaseLock, ConfigMap, Endpoint)

## Usage
...
