# Phase 2 Implementation: Resource Scoping

## Overview

This document provides detailed implementation steps for Phase 2 of the CNAMERecord migration: adding namespace awareness to all CNAMERecord operations.

**Objective:** Ensure all CNAMERecord Get/List/Delete operations are namespace-scoped.

**Success Criteria:**
- All Delete operations target correct namespace
- List operations scoped to Capp namespace  
- No cluster-wide resource access
- Code compiles successfully

## Current State After Phase 1

**Completed:**
- Import paths migrated to provider-dns-v2 ✅
- Namespace set in CNAMERecord creation (line 67 in `dnsrecord.go`) ✅
- Get operations use `types.NamespacedName` with namespace ✅

**Remaining Work:**
- `GetBareDNSRecord()` helper function lacks namespace parameter
- List operations not explicitly namespace-scoped
- Delete operations using helper without namespace

## Implementation Steps

### Step 1: Update GetBareDNSRecord Helper Function

**File:** `internal/kinds/capp/resourceclient/resourcepreparers.go`

**Location:** Lines 73-80

**Current:**
```go
func GetBareDNSRecord(name string) dnsrecordv1alpha1.CNAMERecord {
    return dnsrecordv1alpha1.CNAMERecord{
        ObjectMeta: metav1.ObjectMeta{
            Name: name,
        },
    }
}
```

**Target:**
```go
func GetBareDNSRecord(name, namespace string) dnsrecordv1alpha1.CNAMERecord {
    return dnsrecordv1alpha1.CNAMERecord{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: namespace,
        },
    }
}
```

**Changes:**
1. Add `namespace string` parameter to function signature
2. Set `Namespace: namespace` in ObjectMeta

**Rationale:** Namespaced Kubernetes resources require namespace in ObjectMeta for Delete operations to target the correct namespace.

---

### Step 2: Update CleanUp Function Call

**File:** `internal/kinds/capp/resourcemanagers/dnsrecord.go`

**Location:** Line 93

**Current:**
```go
dnsRecord := rclient.GetBareDNSRecord(capp.Status.RouteStatus.DomainMappingObjectStatus.URL.Host)
```

**Target:**
```go
dnsRecord := rclient.GetBareDNSRecord(capp.Status.RouteStatus.DomainMappingObjectStatus.URL.Host, capp.Namespace)
```

**Changes:**
1. Add `capp.Namespace` as second argument to `GetBareDNSRecord()`

---

### Step 3: Update deletePreviousDNSRecords Function Call

**File:** `internal/kinds/capp/resourcemanagers/dnsrecord.go`

**Location:** Line 214

**Current:**
```go
recordset := rclient.GetBareDNSRecord(dnsRecord.Name)
```

**Target:**
```go
recordset := rclient.GetBareDNSRecord(dnsRecord.Name, dnsRecord.Namespace)
```

**Changes:**
1. Add `dnsRecord.Namespace` as second argument to `GetBareDNSRecord()`

---

### Step 4: Add Namespace Scoping to List Operations

**File:** `internal/kinds/capp/resourcemanagers/dnsrecord.go`

**Location:** Lines 201-203 in `getPreviousDNSRecords()` function

**Current:**
```go
listOptions := utils.GetListOptions(set)

if err := r.K8sclient.List(r.Ctx, &dnsRecords, &listOptions); err != nil {
```

**Target:**
```go
listOptions := utils.GetListOptions(set)
listOptions.Namespace = capp.Namespace

if err := r.K8sclient.List(r.Ctx, &dnsRecords, &listOptions); err != nil {
```

**Changes:**
1. Add `listOptions.Namespace = capp.Namespace` after creating list options

**Rationale:** Without explicit namespace scoping, List operations may attempt cluster-wide searches. For namespaced resources, we must scope to the specific namespace.

---

### Step 5: Build Verification

Verify changes compile successfully:

```bash
go fmt ./...
go vet ./...
go build -o bin/manager cmd/main.go
```

**Expected:** Binary created at `bin/manager` with no compilation errors.

---

## Validation Checklist

After completing all steps:

- [ ] `GetBareDNSRecord` function signature updated with namespace parameter
- [ ] All 2 callers of `GetBareDNSRecord` updated with namespace argument
- [ ] List operation explicitly sets `Namespace` field in options
- [ ] `go fmt ./...` passes
- [ ] `go vet ./...` passes  
- [ ] `go build -o bin/manager cmd/main.go` succeeds
- [ ] No compilation errors

## Phase 2 Scope

**Includes:**
- Helper function signature updated for namespace support
- All Delete operations namespace-aware
- List operations explicitly namespace-scoped

**Excludes (deferred to later phases):**
- RBAC configuration changes (Phase 3)
- Test updates (Phase 4)
- Deployment verification (Phase 5)


