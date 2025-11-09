# Phase 1 Implementation: Dependency Migration

## Overview

This document provides detailed implementation steps for Phase 1 of the CNAMERecord migration from provider-dns v0.1.3 (cluster-scoped) to provider-dns-v2 v1.0.1 (namespaced).

**Objective:** Update all Go module dependencies and import paths to use provider-dns-v2 without changing functionality.

**Success Criteria:** 
- `make build` completes successfully
- All import statements reference provider-dns-v2
- No compilation errors related to DNS types

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

```bash
# Regenerate code if needed
make generate

# Build the operator
make build
```

**Expected Output:**
- No compilation errors
- Binary created successfully

**Common Issues:**

1. **"undefined: dnsrecordv1alpha1.CNAMERecordStatus"**
   - Cause: v2 API renamed or removed this type
   - Fix: Check v2 API structure and update code accordingly

2. **"cannot find package"**
   - Cause: Import path incorrect or module not downloaded
   - Fix: Verify import path matches v2 repo structure

3. **Field compatibility errors**
   - Cause: v2 changed struct fields
   - Fix: Update field references to match v2 API

---

## Validation Checklist

After completing all steps:

- [ ] `go.mod` contains `provider-dns-v2 v1.0.1`
- [ ] `go.mod` does NOT contain `provider-dns v0.1.3` in direct dependencies
- [ ] All 10 source files updated with new import path
- [ ] Import alias corrected from `dnsvrecord1alpha1` to `dnsrecordv1alpha1`
- [ ] `go mod tidy` completes without errors
- [ ] `make build` completes successfully
- [ ] No compilation errors related to DNS types
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

Phase 1 only updates imports - it does NOT:
- Add namespace to CNAMERecord ObjectMeta
- Change Get/List/Delete operation scopes
- Modify helper function signatures
- Update RBAC configurations

These changes are intentionally deferred to Phase 2 to maintain clear separation of concerns and enable incremental validation.

