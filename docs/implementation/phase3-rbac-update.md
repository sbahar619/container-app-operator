# Phase 3 Implementation: RBAC Update

## Overview

This phase updates RBAC to align with the namespaced `CNAMERecord` CR from `provider-dns-v2`. It changes the API group from `record.dns.crossplane.io` to `record.dns-v2.crossplane.io` everywhere RBAC is declared or generated.

**Objective:** Ensure the operator has correct permissions to manage `cnamerecords` in the new API group.

**Success Criteria:**
- Generated RBAC manifests reference `record.dns-v2.crossplane.io`
- Helm chart RBAC template updated to `record.dns-v2.crossplane.io`
- `make manifests` regenerates RBAC successfully (or acceptable alternative on Go 1.24.x)


## Current State After Phase 2

**Completed:**
- Dependency migration (Phase 1)
- Resource scoping (Phase 2) for namespaced `CNAMERecord` operations

**Remaining for Phase 3:**
- Update `+kubebuilder:rbac` markers
- Regenerate RBAC manifests
- Update Helm chart RBAC template


## Implementation Steps

### Step 1: Update kubebuilder RBAC markers

- File: `internal/kinds/capp/controllers/controller.go`

Change the `+kubebuilder:rbac` marker that references the DNS provider group:

```go
// OLD
// +kubebuilder:rbac:groups="record.dns.crossplane.io",resources=cnamerecords,verbs=get;list;watch;update;create;delete

// NEW
// +kubebuilder:rbac:groups="record.dns-v2.crossplane.io",resources=cnamerecords,verbs=get;list;watch;update;create;delete
```

Notes:
- Keep verbs unchanged unless your deployment model requires fewer permissions.
- We remain on ClusterRole (per HLD) to support multi-namespace management.


### Step 2: Regenerate manifests

Run:

```bash
make manifests
```

If your environment uses Go 1.24.x and generation fails due to swiss map errors, use one of:

```bash
# Option A: Disable swiss maps just for generation
GOEXPERIMENT=noswissmap make manifests

# Option B: Run controller-gen via Docker with Go 1.23.x (example)
docker run --rm -v "$PWD":/workspace -w /workspace golang:1.23 \
  bash -lc 'go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.16.4 && \
  $(go env GOPATH)/bin/controller-gen rbac:roleName=manager-role crd:allowDangerousTypes=true webhook paths="./..." output:crd:artifacts:config=config/crd/bases'
```


### Step 3: Verify generated RBAC

- File: `config/rbac/role.yaml`

Confirm the DNS API group is updated:

```yaml
rules:
  - apiGroups:
      - record.dns-v2.crossplane.io  # UPDATED
    resources:
      - cnamerecords
    verbs: [ create, delete, get, list, update, watch ]
```


### Step 4: Update Helm chart RBAC template

- File: `charts/container-app-operator/templates/manager-rbac.yaml`

Change the DNS API group:

```yaml
# OLD
- apiGroups:
  - record.dns.crossplane.io
  resources:
  - cnamerecords

# NEW
- apiGroups:
  - record.dns-v2.crossplane.io
  resources:
  - cnamerecords
```


### Step 5: Sanity grep for any remaining references

Run:

```bash
rg --hidden --line-number "record\\.dns\\.crossplane\\.io"
```

All matches should be removed or updated to `record.dns-v2.crossplane.io`.


## Build

```bash
go fmt ./...
go vet ./...
go build -o bin/manager cmd/main.go
```


## Validation

1. Check generated and chart RBAC files contain the new API group:
   - `config/rbac/role.yaml`
   - `charts/container-app-operator/templates/manager-rbac.yaml`
2. Validate permissions locally (replace NAMESPACE/SA name accordingly):

```bash
kubectl auth can-i --as=system:serviceaccount:<NAMESPACE>:<RELEASE>-controller-manager \
  --verb=get --resource=cnamerecords --api-group=record.dns-v2.crossplane.io
kubectl auth can-i --as=system:serviceaccount:<NAMESPACE>:<RELEASE>-controller-manager \
  --verb=create --resource=cnamerecords --api-group=record.dns-v2.crossplane.io
```

3. Deploy to a test cluster (Phase 5 will formalize) and create a `Capp`; verify:
   - `CNAMERecord` reconcile succeeds without RBAC denials
   - Events/Logs show successful CRUD on `cnamerecords` using `record.dns-v2.crossplane.io`


## Scope

**Includes:**
- Update kubebuilder RBAC markers
- Regenerate manifests
- Update Helm RBAC template

**Excludes (deferred to later phases):**
- Test updates (Phase 4)
- Deployment verification (Phase 5)


