# Leader Elector

A simple leader election sidecar container for Kubernetes based on the recent [client-go](https://github.com/kubernetes/client-go) library. It aims to provide an alternative to the widely used but outdated `election` component of [kubernetes/contrib](https://github.com/kubernetes-retired/contrib).


## Configuration Flags

* `election` - name of the election and the corresponding lock resource.
* `namespace` - the Kubernetes namespace to run the election in.
* `locktype` - the leaselock resource type to use for this deployment. Supported lock types are:
  * `configmaps` - default
  * `leases`
  * `endpoints`
* `port` - the port on which the election sidecar can be queried.

## Example
A working example deployment can be found under `example/deployment.yaml`

## Work in Progress
* Configurable lock type (LeaseLock, ConfigMap, Endpoint)
* Make configuration fail safe

