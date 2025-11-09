# High-Level Design: Migrate to Namespaced CNAMERecord CR

## 1. Overview

### Current State
- Operator uses provider-dns v0.1.3 with cluster-scoped CNAMERecord CRs
- CNAMERecords created without namespace in ObjectMeta
- RBAC configured with ClusterRole for cluster-wide access
- API Group: `record.dns.crossplane.io/v1alpha1`

### Target State
- Operator uses provider-dns-v2 v1.0.1 with namespaced CNAMERecord CRs
- CNAMERecords scoped to Capp namespace
- RBAC updated for v2 API Group (potentially remains ClusterRole for multi-namespace support)
- API Group: `record.dns-v2.crossplane.io/v1alpha1`

### Migration Objectives
- Enable namespace isolation for DNS records
- Maintain backward compatibility in operator behavior
- Ensure proper RBAC permissions for namespaced resources
- Validate through comprehensive testing strategy

## 2. Architecture Components Affected

### 2.1 Dependency Layer
**Component:** Go module dependencies and API imports

**Changes:**
- Replace provider-dns v0.1.3 dependency with provider-dns-v2 v1.0.1
- Update import paths across all modules (core operator, tests, utilities)
- Update Kubernetes scheme registration

**Impact Analysis:**
- 9 source files require import path updates
- API type compatibility must be verified (CNAMERecord, CNAMERecordSpec, CNAMERecordStatus, CNAMERecordList)
- Potential struct field changes in v2 API need investigation

### 2.2 Resource Management Layer
**Component:** DNSRecord manager reconciliation logic

**Key Design Decision:** Namespace Ownership Model
- Each Capp CR owns CNAMERecord in its own namespace
- CNAMERecord lifecycle tied to parent Capp namespace
- Namespace becomes part of resource identity

**Changes Required:**
- Resource preparation: Include namespace in CNAMERecord ObjectMeta
- Resource lookup: Ensure Get/List operations are namespace-scoped
- Resource cleanup: Delete operations target correct namespace
- Helper functions: Accept namespace as parameter where applicable

**Side Benefit:** Fixes existing deletion bug (incorrect name comparison in cleanup logic)

### 2.3 Authorization Layer
**Component:** RBAC permissions model

**Key Design Decision:** ClusterRole vs Role
- **Option A (Recommended):** Keep ClusterRole
  - Rationale: Operator manages Capps across multiple namespaces
  - Trade-off: Broader permissions but operational flexibility
- **Option B:** Use namespace-scoped Role
  - Rationale: If operator deployed per-namespace
  - Trade-off: Enhanced security but requires role-per-namespace deployment

**Changes Required:**
- Update API Group reference from `record.dns.crossplane.io` to `record.dns-v2.crossplane.io`
- Regenerate RBAC manifests via controller-gen
- Update Helm chart RBAC templates

### 2.4 Status Reporting Layer
**Component:** Capp status synchronization

**Changes:**
- Capp status embeds CNAMERecordStatus from provider-dns-v2
- Verify v2 status structure compatibility
- Update status sync logic if v2 introduces new fields

## 3. Testing Strategy Architecture

### 3.1 Test Pyramid Approach

**Layer 1: Unit Tests (Adjust Existing)**
- Focus: Resource preparation logic, helper functions
- Validation: Namespace correctly set in all resource operations
- Coverage: Existing unit tests adjusted for namespaced API

**Layer 2: Integration Tests (Adjust Existing)**
- Focus: Operator reconciliation with v2 CRDs
- Validation: CNAMERecord created/updated/deleted in correct namespace
- Coverage: Existing E2E tests updated for import changes

**Layer 3: New Validation Tests (Add If Gaps Exist)**
- Focus: Namespace isolation, RBAC enforcement, multi-namespace scenarios
- Validation: No cross-namespace resource access, proper permission enforcement
- Coverage: New tests only if existing coverage insufficient

### 3.2 Validation Phases

**Phase 1: Regression Prevention**
- Run existing test suite with updated imports
- Ensure no behavioral changes except namespace scoping
- Identify and fix breaking changes in v2 API

**Phase 2: Coverage Analysis**
- Review test coverage for namespace-specific scenarios
- Determine if new tests required (namespace isolation, RBAC)
- Document coverage gaps

**Phase 3: Pre-Production Validation**
- Manual testing in test cluster
- Real DNS provider integration
- Monitoring and observability validation

## 4. Implementation Phases

### Phase 1: Dependency Migration
**Objective:** Update all references to use provider-dns-v2

**Deliverables:**
- Updated go.mod with v2 dependency
- All import paths migrated to v2 API package
- Build passes without compilation errors

**Validation:** `make build` succeeds

### Phase 2: Resource Scoping Implementation
**Objective:** Add namespace awareness to all CNAMERecord operations

**Deliverables:**
- Resource creation includes namespace
- Get/List/Delete operations namespace-scoped
- Helper functions accept namespace parameter

**Validation:** Code review confirms namespace in all resource operations

### Phase 3: RBAC Update
**Objective:** Update permissions for v2 API Group

**Deliverables:**
- RBAC manifests reference new API Group
- Generated CRDs and manifests updated
- Helm charts reflect RBAC changes

**Validation:** `make manifests` generates correct RBAC

### Phase 4: Test Validation
**Objective:** Ensure operator behavior validated via tests

**Deliverables:**
- Existing tests adjusted for v2 API
- New tests added if coverage gaps identified
- All tests passing

**Validation:** `make test && make test-e2e` passes

### Phase 5: Deployment Verification
**Objective:** Validate in live environment

**Deliverables:**
- Operator deployed to test cluster
- CNAMERecords created in correct namespaces
- DNS resolution functional
- Cleanup operations successful

**Validation:** Manual verification checklist completed

## 5. Success Criteria

**Functional Requirements:**
- Operator creates CNAMERecord in same namespace as Capp
- DNS records functional and resolvable
- Cleanup deletes CNAMERecord from correct namespace
- No cross-namespace resource access

**Non-Functional Requirements:**
- All existing tests pass with updated API
- Build pipeline succeeds
- RBAC follows principle of least privilege
- Documentation updated for v2 migration

**Operational Requirements:**
- Operator deploys successfully with v2 API
- Monitoring/logging unchanged or improved
- No performance degradation
- Rollback procedure documented

## 6. Key Artifacts for LLD

For subsequent Low-Level Design documentation, the following areas require detailed specifications:

1. **API Type Mapping:** Complete struct-level mapping between v1 and v2 API types
2. **Code Changes Specification:** File-by-file, function-by-function change list with line numbers
3. **Test Case Specifications:** Detailed test scenarios with assertions and expected outcomes
4. **RBAC Policy Details:** Exact permissions required per namespace/cluster scope
5. **Deployment Procedures:** Step-by-step upgrade and rollback procedures
6. **Migration Scripts:** Tooling for cleaning up legacy cluster-scoped resources

