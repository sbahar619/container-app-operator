# Provider DNS v2: Enable Namespaced ProviderConfig for Namespaced CNAMERecord

## Goal

Make namespaced CNAMERecord (`record.dns-v2.m.crossplane.io/v1alpha1`) work with a namespaced `ProviderConfig` (`dns-v2.m.crossplane.io/v1beta1`) by ensuring the provider binary registers the namespaced config types in BOTH:
- The global/controller manager scheme
- The Terraform connector scheme used to resolve `ProviderConfig`

Symptoms when missing: CNAMERecord shows `Synced=False` with message “unknown GVK for ProviderConfig: no version dns-v2.m.crossplane.io/v1beta1 has been registered in scheme” during reconcile.


## What Must Be True

- CNAMERecord is namespaced and created in the Capp namespace.
- A namespaced `ProviderConfig` named `default` exists in the same namespace as the CNAMERecord.
- The provider binary registers:
  - `record.dns-v2.m.crossplane.io/v1alpha1` (CNAMERecord)
  - `dns-v2.m.crossplane.io/v1beta1` (ProviderConfig, ProviderConfigUsage)
- The Terraform connector resolves `ProviderConfig` using a scheme that includes `dns-v2.m.crossplane.io/v1beta1` (critical).


## Step-by-Step: Provider Code Changes

These edits are in the provider repo (`github.com/dana-team/provider-dns-v2`), not in this operator.

### 1) Register namespaced APIs in the global manager scheme

- File: `cmd/provider/main.go` (or wherever the root `scheme` is initialized)

Add imports:

```go
import (
	// ...
	configm "github.com/dana-team/provider-dns-v2/apis/namespaced/config/v1beta1"
	recordm "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)
```

Register into the same `scheme` used by the manager:

```go
func init() {
	// existing registrations ...
	utilruntime.Must(configm.AddToScheme(scheme)) // dns-v2.m.crossplane.io/v1beta1
	utilruntime.Must(recordm.AddToScheme(scheme)) // record.dns-v2.m.crossplane.io/v1alpha1
}
```


### 2) Ensure the ProviderConfig resolver uses a scheme that includes namespaced APIs

- File: `internal/clients/dns-v2.go`

In this provider, the resolver relies on the controller-runtime client's scheme (i.e., the manager's scheme), not a standalone scheme:

```go
// resolveModern(...)
pcRuntimeObj, err := crClient.Scheme().New(namespacedv1beta1.
	SchemeGroupVersion.WithKind(configRef.Kind))
```

This means step 1 is sufficient: registering `apis/namespaced` into the manager’s scheme via `apisNamespaced.AddToScheme(mgr.GetScheme())` makes the resolver aware of:
- `dns-v2.m.crossplane.io/v1beta1` (ProviderConfig, ProviderConfigUsage, ClusterProviderConfig)
- `record.dns-v2.m.crossplane.io/v1alpha1` (CNAMERecord)

Only if you introduce a custom, separate `runtime.Scheme` for decoding would you also need to explicitly call `AddToScheme` on that custom scheme for the namespaced packages.


## Build & Deploy the Provider

1) Build and push an image with the above changes (example):

```bash
IMAGE=ghcr.io/dana-team/provider-dns-v2:namespaced-fix
docker build -t $IMAGE .
docker push $IMAGE
```

2) Upgrade the provider in the cluster and restart:

```bash
kubectl -n crossplane-system patch providers.pkg.crossplane.io provider-dns \
  --type merge -p '{"spec":{"package":"'"$IMAGE"'"}}'
kubectl -n crossplane-system rollout restart deploy/provider-dns
```


## Runtime Validation

1) Provider logs should show registration and no unknown-GVK errors during reconcile:

```bash
kubectl -n crossplane-system logs deploy/provider-dns | grep -i 'dns-v2\.m\.crossplane\.io/v1beta1'
```

2) Ensure a namespaced ProviderConfig `default` exists alongside the CNAMERecord:

```bash
kubectl -n <capp-ns> get providerconfigs.dns-v2.m.crossplane.io default
```

3) A `ProviderConfigUsage` should be created in that namespace:

```bash
kubectl -n <capp-ns> get providerconfigusages.dns-v2.m.crossplane.io
```

4) Requeue the record and check conditions:

```bash
kubectl -n <capp-ns> annotate cnamerecords.record.dns-v2.m.crossplane.io <name> reconcilenow="$(date +%s)" --overwrite
kubectl -n <capp-ns> get cnamerecords.record.dns-v2.m.crossplane.io <name> -o jsonpath='{.status.conditions[*]}{"\n"}'
```

Expected: `Synced` progresses and the “unknown GVK … dns-v2.m.crossplane.io/v1beta1” error no longer appears.


## Operator Alignment (FYI)

In this operator:
- Keep `CNAMERecord` namespaced and set `ObjectMeta.Namespace` to the Capp namespace.
- Keep `spec.providerConfigRef.name: "default"` (no group/kind needed).
- Ensure RBAC uses the namespaced record API group (`record.dns-v2.m.crossplane.io`) and that the manager watches multiple namespaces if needed via ClusterRole.


## References

- Provider repo and examples: `dana-team/provider-dns-v2` [README](https://github.com/dana-team/provider-dns-v2)


