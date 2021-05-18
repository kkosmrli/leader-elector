# Leader Elector

A simple leader election sidecar container for Kubernetes based on the recent [client-go](https://github.com/kubernetes/client-go) library. It aims to provide an alternative to the widely used but outdated `election` component of [kubernetes/contrib](https://github.com/kubernetes-retired/contrib).


## Configuration Flags

The following arguments can be passed on the command line, or by setting the environment variable named in brackets.

* `election` - name of the election and the corresponding lock resource (ELECTION_NAME).
* `namespace` - the Kubernetes namespace to run the election in (ELECTION_NAMESPACE).
* `locktype` - the resource type to use as the lock for this deployment (ELECTION_TYPE). Supported lock types are:
  * `configmaps` - default
  * `leases`
  * `endpoints`
* `port` - the port on which the election sidecar can be queried (ELECTION_PORT).

## Example
A working example deployment can be found under `example/deployment.yaml`. The necessary Role and RoleBindings must be applied prior via `example/rbac.yaml`.

### Sidecar Response
```
{
  "name": "election-example-7548f6f8f7-fml47"
}
```



