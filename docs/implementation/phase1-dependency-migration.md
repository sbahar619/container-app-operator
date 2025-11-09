# Phase 1 Implementation: Dependency Migration

## Overview

This document provides detailed implementation steps for Phase 1 of the CNAMERecord migration from provider-dns v0.1.3 (cluster-scoped) to provider-dns-v2 v1.0.1 (namespaced).

**Objective:** Update all Go module dependencies and import paths to use provider-dns-v2, including Crossplane Runtime v2 API changes.

**Success Criteria:** 
- `make build` completes successfully
- All import statements reference provider-dns-v2
- Crossplane Runtime v2 APIs properly integrated
- No compilation errors related to DNS types or Crossplane APIs

## API Compatibility Verification

**Status:** ✅ Verified compatible

The provider-dns-v2 namespaced API has been verified to contain all required types:

| Type | Status | Notes |
|------|--------|-------|
| CNAMERecord | ✅ Exists | Namespaced resource (`scope=Namespaced`) |
| CNAMERecordSpec | ✅ Compatible | Standard Crossplane resource spec |
| CNAMERecordStatus | ✅ Compatible | Contains `ResourceStatus` + `AtProvider` fields |
| CNAMERecordList | ✅ Exists | Standard Kubernetes list type |

**Import Path:** `github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1`

**API Group:** `record.dns-v2.m.crossplane.io/v1alpha1`

## Prerequisites

- Understanding of Go module management
- Access to provider-dns-v2 v1.0.1 repository
- Familiarity with operator codebase structure
- Review [Coding Standards](../project/CODING_STANDARDS.md) before implementation

## Implementation Steps

### Step 1: Update go.mod Dependency

**File:** `go.mod`

**Current State (line 11):**
```go
github.com/dana-team/provider-dns v0.1.3
```

**Target State:**
```go
github.com/dana-team/provider-dns-v2 v1.0.1
```

**Actions:**
1. Open `go.mod`
2. Locate the `require` block (starting around line 5)
3. Find the line containing `github.com/dana-team/provider-dns v0.1.3`
4. Replace with `github.com/dana-team/provider-dns-v2 v1.0.1`
5. Save the file

**Note:** Do NOT run `go mod tidy` yet - wait until all imports are updated.

---

### Step 2: Update Import Paths in Source Files

The following 10 files require import path updates. Each must be updated from:
- **FROM:** `github.com/dana-team/provider-dns/apis/record/v1alpha1`
- **TO:** `github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1`

#### 2.1 Core Operator Entry Point

**File:** `cmd/main.go`

**Location:** Line 29

**Current:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
```

**Target:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
```

**Context:** This import is used for scheme registration (line 71) to ensure the Kubernetes client can recognize CNAMERecord types.

---

#### 2.2 API Types Definition

**File:** `api/v1alpha1/capp_types.go`

**Location:** Import section (find exact line with grep)

**Current:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
```

**Target:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
```

**Context:** Used in `DNSRecordObjectStatus` struct (around line 196) to embed CNAMERecordStatus.

**Verification Completed:** ✅ CNAMERecordStatus type exists in v2 with compatible structure:
```go
type CNAMERecordStatus struct {
	v1.ResourceStatus `json:",inline"`
	AtProvider        CNAMERecordObservation `json:"atProvider,omitempty"`
}
```

---

#### 2.3 DNS Record Manager

**File:** `internal/kinds/capp/resourcemanagers/dnsrecord.go`

**Location:** Line 8

**Current:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
```

**Target:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
```

**Context:** This is the core reconciliation logic that creates, updates, and deletes CNAMERecord resources.

**Impact:** Used throughout the file for CNAMERecord, CNAMERecordSpec, CNAMERecordList types.

---

#### 2.4 Resource Preparers (Helper Functions)

**File:** `internal/kinds/capp/resourceclient/resourcepreparers.go`

**Location:** Line 6

**Current:**
```go
dnsvrecord1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
```

**Target:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
```

**Special Note:** 
- Current import alias has a typo: `dnsvrecord1alpha1` (extra 'v')
- Change to: `dnsrecordv1alpha1` (consistent with other files)
- This aligns with coding standards: **clear variable names** and **no code duplication**
- Consistent naming improves code readability and maintainability

**Impact:** Affects `GetBareDNSRecord` helper function (line 74).

---

#### 2.5 Controller Setup

**File:** `internal/kinds/capp/controllers/controller.go`

**Location:** Import section (find with grep)

**Current:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
```

**Target:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
```

**Context:** Used in controller's `SetupWithManager` to watch CNAMERecord resources.

---

#### 2.6 Status Sync Logic

**File:** `internal/kinds/capp/status/route.go`

**Location:** Import section (find with grep)

**Current:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
```

**Target:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
```

**Context:** Used to sync CNAMERecord status into Capp status.

---

#### 2.7 Route Utilities

**File:** `internal/kinds/capp/utils/route.go`

**Location:** Import section (find with grep)

**Current:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
```

**Target:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
```

**Context:** Utility functions for route/DNS operations.

---

#### 2.8 E2E Test Helper

**File:** `test/e2e_tests/helper.go`

**Location:** Import section (find with grep)

**Current:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
```

**Target:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
```

**Context:** Test setup and helper functions for E2E tests.

---

#### 2.9 E2E Test Mocks

**File:** `test/e2e_tests/mocks/route.go`

**Location:** Import section (find with grep)

**Current:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns/apis/record/v1alpha1"
```

**Target:**
```go
dnsrecordv1alpha1 "github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1"
```

**Context:** Mock objects for route-related tests.

**Note:** File `test/e2e_tests/utils/route_adapter.go` also uses this import and should be updated similarly.

---

### Step 2.5: Update Crossplane Runtime API References

**CRITICAL:** Provider-dns-v2 uses **Crossplane Runtime v2** instead of v1. This requires additional code changes beyond import updates.

#### Breaking Change: API Version Upgrade

Provider-dns-v2 `CNAMERecordSpec` embeds `v2.ManagedResourceSpec` (not `v1.ResourceSpec`), which means:
- Different field access patterns
- Different condition checking APIs
- Need to import both v1 and v2 APIs from crossplane-runtime

#### 2.5.1 Update DNS Record Manager

**File:** `internal/kinds/capp/resourcemanagers/dnsrecord.go`

**Current (lines ~10, 79-83):**
```go
import (
    xpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"  // OLD - will be removed
    // ... other imports ...
)

// In prepareResource function:
Spec: dnsrecordv1alpha1.CNAMERecordSpec{
    ForProvider: dnsrecordv1alpha1.CNAMERecordParameters{
        Name:  &recordName,
        Zone:  &zone,
        Cname: &cname,
    },
    ResourceSpec: xpcommonv1.ResourceSpec{
        ProviderConfigReference: &xpcommonv1.Reference{
            Name: xpProvider,
        },
    },
},
```

**Target:**
```go
import (
    // REMOVE: xpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
    xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"  // ADD this
    // ... other imports remain unchanged ...
)

// In prepareResource function (around line 63-87):
dnsRecord := dnsrecordv1alpha1.CNAMERecord{
    TypeMeta: metav1.TypeMeta{},
    ObjectMeta: metav1.ObjectMeta{
        Name: resourceName,
        // Namespace will be set in Step 2.6.1
        Labels: map[string]string{
            utils.CappResourceKey:   capp.Name,
            utils.CappNamespaceKey:  capp.Namespace,
            utils.ManagedByLabelKey: utils.CappKey,
        },
    },
    Spec: dnsrecordv1alpha1.CNAMERecordSpec{
        ForProvider: dnsrecordv1alpha1.CNAMERecordParameters{
            Name:  &recordName,
            Zone:  &zone,
            Cname: &cname,
        },
    },
}
// Set ProviderConfigReference on the embedded ManagedResourceSpec
dnsRecord.Spec.ProviderConfigReference = &xpv1.ProviderConfigReference{Name: xpProvider}

return dnsRecord, nil
```

**Changes:**
1. **REMOVE** old import: `xpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"`
2. **ADD** new import: `xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"` (note: v1 API from v2 module)
3. Remove the inline `ResourceSpec` struct literal from CNAMERecordSpec
4. Set `ProviderConfigReference` on the embedded field after constructing the main struct
5. Use `xpv1.ProviderConfigReference` instead of `xpcommonv1.Reference`

**Note:** Namespace support (setting `Namespace` in ObjectMeta) will be addressed in **Phase 2** according to the HLD.

---

### Step 2.6: Update Crossplane Runtime API References (Continued)

This continues the Crossplane Runtime v2 API updates from Step 2.5.

---

#### 2.6.1 Update Route Utilities

**File:** `internal/kinds/capp/utils/route.go`

**Current (lines ~12, 31-34):**
```go
import (
    xpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"  // OLD - will be removed
    // ... other imports ...
)

// In IsDNSRecordAvailable function:
if dnsRecord.Status.Conditions != nil {
    readyCondition := dnsRecord.Status.GetCondition(xpcommonv1.TypeReady)
    available = readyCondition.Equal(xpcommonv1.Available())
}
```

**Target:**
```go
import (
    // REMOVE: xpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
    xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"  // ADD this
    corev1 "k8s.io/api/core/v1"  // ADD this - for ConditionStatus constants
    // ... other imports remain unchanged ...
)

// In IsDNSRecordAvailable function (around line 31-34):
if dnsRecord.Status.Conditions != nil {
    readyCondition := dnsRecord.Status.GetCondition(xpv1.TypeReady)
    available = readyCondition.Status == corev1.ConditionTrue
}
```

**Changes:**
1. **REMOVE** old import: `xpcommonv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"`
2. **ADD** new import: `xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"` (note: v1 API but from v2 module path)
3. **ADD** new import: `corev1 "k8s.io/api/core/v1"` for Kubernetes standard condition constants
4. Replace all uses of `xpcommonv1` with `xpv1`
5. Update condition checking: use `.Status == corev1.ConditionTrue` instead of `.Equal()` method
6. **IMPORTANT**: Crossplane Runtime v2 uses standard Kubernetes `corev1.ConditionStatus` type instead of crossplane-specific condition constants

---

#### 2.6.2 Update go.mod (if needed)

Check if crossplane-runtime v2 is already in your dependencies:

```bash
grep "crossplane-runtime" go.mod
```

If you see only v1 (e.g., `github.com/crossplane/crossplane-runtime v1.x.x`), you need to add v2:

```bash
go get github.com/crossplane/crossplane-runtime/v2@latest
```

**Note:** You may have both v1 and v2 dependencies - this is normal as they're separate module paths.

---

### Step 3: Resolve Dependencies

After all import paths are updated:

```bash
# Clean module cache to avoid conflicts
go clean -modcache

# Download new dependencies
go mod download

# Tidy up go.mod and go.sum
go mod tidy
```

**Expected Outcome:**
- `go.mod` shows `github.com/dana-team/provider-dns-v2 v1.0.1` in require block
- `go.sum` contains checksums for provider-dns-v2
- No references to old provider-dns v0.1.3 remain (except possibly in indirect dependencies)

---

### Step 4: Verify API Compatibility

Before building, verify that v2 API types are compatible:

```bash
# Check what types are available in v2
go doc github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1

# Specifically verify these types exist:
go doc github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1.CNAMERecord
go doc github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1.CNAMERecordSpec
go doc github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1.CNAMERecordStatus
go doc github.com/dana-team/provider-dns-v2/apis/namespaced/record/v1alpha1.CNAMERecordList
```

**If types are missing or incompatible:**
- Document the differences
- May need struct field mapping
- May require additional code changes beyond import updates

---

### Step 5: Build Verification

**Note:** Phase 1 focuses on dependency migration and compilation. If you encounter issues with `make generate` (e.g., Go toolchain bugs), you can skip it if no API types were modified.

#### Option A: Full Build with Makefile (Recommended)

```bash
# Full build including code generation
make build
```

This runs `make generate` → `make manifests` → build automatically.

#### Option B: Direct Build (If make generate fails)

If `make generate` fails due to toolchain issues and you **only changed import paths** (no API type changes):

```bash
# Format and vet
go fmt ./...
go vet ./...

# Build directly
go build -o bin/manager cmd/main.go
```

**Expected Output:**
- No compilation errors
- Binary created successfully at `bin/manager`

**Common Issues:**

1. **`make generate` fails with "mapiterinit redeclared" errors**
   - Cause: Go 1.24.x experimental versions have known bugs with swiss maps
   - Fix: Use Option B (direct build) or downgrade to stable Go 1.23.x
   - This is a Go toolchain issue, not a code issue

2. **"undefined: dnsrecordv1alpha1.CNAMERecordStatus"**
   - Cause: v2 API renamed or removed this type
   - Fix: Check v2 API structure and update code accordingly

3. **"cannot find package"**
   - Cause: Import path incorrect or module not downloaded
   - Fix: Verify import path matches v2 repo structure

4. **"unknown field ResourceSpec in CNAMERecordSpec"**
   - Cause: v2 embeds `v2.ManagedResourceSpec` differently than v1
   - Fix: Set `ProviderConfigReference` after struct construction (see Step 2.5.1)

5. **"undefined: xpcommonv1.ResourceSpec"**
   - Cause: Crossplane Runtime v2 API changes
   - Fix: Use `xpv1.ProviderConfigReference` from crossplane-runtime/v2 (see Step 2.5)

6. **".Equal undefined (type xpv1.Condition has no field or method Equal)"**
   - Cause: Crossplane Runtime v2 uses standard Kubernetes condition types
   - Fix: Use `.Status == corev1.ConditionTrue` with `corev1 "k8s.io/api/core/v1"` import (see Step 2.6.1)

---

## Validation Checklist

After completing all steps:

- [ ] `go.mod` contains `provider-dns-v2 v1.0.1`
- [ ] `go.mod` does NOT contain `provider-dns v0.1.3` in direct dependencies
- [ ] All 10 source files updated with new import path
- [ ] Import alias corrected from `dnsvrecord1alpha1` to `dnsrecordv1alpha1`
- [ ] Crossplane Runtime v2 imports updated in `dnsrecord.go` (`xpv1` from v2 module)
- [ ] Crossplane Runtime v2 imports updated in `route.go` (`xpv1` from v2 module)
- [ ] `corev1` import added to `route.go` for condition status constants
- [ ] `ProviderConfigReference` set correctly using `xpv1.ProviderConfigReference`
- [ ] Condition checking updated to use `corev1.ConditionTrue`
- [ ] `go mod tidy` completes without errors
- [ ] `make build` OR `go build -o bin/manager cmd/main.go` completes successfully
- [ ] No compilation errors related to DNS types or Crossplane APIs
- [ ] Generated binary exists in expected location

---

## Rollback Procedure

If issues arise and rollback is needed:

```bash
# Revert all changes
git checkout go.mod
git checkout cmd/main.go
git checkout api/v1alpha1/capp_types.go
git checkout internal/kinds/capp/resourcemanagers/dnsrecord.go
git checkout internal/kinds/capp/resourceclient/resourcepreparers.go
git checkout internal/kinds/capp/controllers/controller.go
git checkout internal/kinds/capp/status/route.go
git checkout internal/kinds/capp/utils/route.go
git checkout test/e2e_tests/helper.go
git checkout test/e2e_tests/mocks/route.go
git checkout test/e2e_tests/utils/route_adapter.go  # If updated

# Clean and restore dependencies
go clean -modcache
go mod download
go mod tidy
```

---

## Next Steps

After Phase 1 completes successfully:
1. Commit changes with descriptive message
2. Create PR for review
3. Proceed to Phase 2: Resource Scoping Implementation
4. Document any API compatibility issues discovered

---

## Coding Standards Compliance

This implementation follows project [Coding Standards](../project/CODING_STANDARDS.md):

**Minimal Changes:**
- Phase 1 only updates imports and dependency declarations
- No logic changes, maintaining existing behavior
- Smallest necessary modifications to achieve build compatibility

**Clear Variable Names:**
- Fixes import alias typo: `dnsvrecord1alpha1` → `dnsrecordv1alpha1`
- Maintains consistent naming convention across all files

**Single Responsibility:**
- Phase focuses solely on dependency migration
- Functional changes deferred to subsequent phases

**No Code Duplication:**
- Standardizes import alias naming across all files
- Removes inconsistency in naming conventions

## Notes for Phase 2

Phase 1 focuses on API compatibility - it updates:
- Import paths to provider-dns-v2
- Crossplane Runtime v1 → v2 API usage
- Type compatibility for CNAMERecord structs

Phase 1 does NOT change functional behavior:
- CNAMERecord namespace is not yet added to ObjectMeta (Phase 2)
- Get/List/Delete operations remain cluster-scoped (Phase 2)  
- Helper function signatures not changed for namespace support (Phase 2)
- RBAC configurations remain unchanged (Phase 3)

These changes are intentionally deferred to maintain clear separation of concerns and enable incremental validation.

