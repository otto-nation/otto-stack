# Phase 1 Implementation Plan: ExecutionContext Refactor

## Overview
Refactor `ExecutionContext` from string-based enum to interface-based type discrimination.

---

## Step-by-Step Implementation

### Step 1: Create New Interface Types (No Breaking Changes)

**File:** `internal/pkg/cli/context/types.go`

**Action:** Add new types alongside existing ones

```go
// NEW: Context interface for type-safe discrimination
type Context interface {
    SharedRoot() string
    isContext() // unexported marker prevents external implementation
}

// NEW: ProjectContext for project-scoped operations
type ProjectContext struct {
    Project *ProjectInfo
    Shared  *SharedInfo
}

func (p *ProjectContext) SharedRoot() string { return p.Shared.Root }
func (p *ProjectContext) isContext()         {}

// NEW: SharedContext for global shared container operations
type SharedContext struct {
    Shared *SharedInfo
}

func (s *SharedContext) SharedRoot() string { return s.Shared.Root }
func (s *SharedContext) isContext()         {}

// KEEP OLD: ExecutionContext (for backward compatibility during migration)
type ExecutionContext struct {
    Type             ContextType
    Project          *ProjectInfo
    SharedContainers *SharedInfo
}
```

**Tests:** Add `context_interface_test.go`
```go
func TestProjectContext_ImplementsContext(t *testing.T) {
    var _ Context = (*ProjectContext)(nil)
}

func TestSharedContext_ImplementsContext(t *testing.T) {
    var _ Context = (*SharedContext)(nil)
}

func TestProjectContext_SharedRoot(t *testing.T) {
    ctx := &ProjectContext{
        Shared: &SharedInfo{Root: "/test/shared"},
    }
    assert.Equal(t, "/test/shared", ctx.SharedRoot())
}
```

---

### Step 2: Add New Detector Method

**File:** `internal/pkg/cli/context/detector.go`

**Action:** Add `DetectContext()` method alongside existing `Detect()`

```go
// NEW: DetectContext returns interface-based context
func (d *Detector) DetectContext() (Context, error) {
    sharedRoot := filepath.Join(d.homeDir, core.OttoStackDir, core.SharedDir)
    if err := os.MkdirAll(sharedRoot, core.PermReadWriteExec); err != nil {
        return nil, err
    }

    sharedInfo := &SharedInfo{Root: sharedRoot}
    project, err := d.findProjectRoot()
    if err != nil {
        return nil, err
    }

    if project != nil {
        return &ProjectContext{
            Project: project,
            Shared:  sharedInfo,
        }, nil
    }

    return &SharedContext{
        Shared: sharedInfo,
    }, nil
}

// KEEP: Detect() for backward compatibility
func (d *Detector) Detect() (*ExecutionContext, error) {
    // Existing implementation unchanged
}
```

**Tests:** Add to `detector_test.go`
```go
func TestDetector_DetectContext_ProjectContext(t *testing.T) {
    // Setup test directory with .otto-stack
    // Call DetectContext()
    // Assert returns *ProjectContext
}

func TestDetector_DetectContext_SharedContext(t *testing.T) {
    // Setup test directory without .otto-stack
    // Call DetectContext()
    // Assert returns *SharedContext
}
```

---

### Step 3: Update Handlers One-by-One

#### 3.1: Update UpHandler

**File:** `internal/pkg/cli/handlers/lifecycle/up.go`

**Before:**
```go
func (h *UpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
    detector, err := clicontext.NewDetector()
    if err != nil {
        return err
    }

    execCtx, err := detector.Detect()
    if err != nil {
        return err
    }

    if execCtx.Type == clicontext.Shared {
        return h.handleGlobalContext(ctx, cmd, args, base, execCtx)
    }

    return h.handleProjectContext(ctx, cmd, args, base, execCtx)
}
```

**After:**
```go
func (h *UpHandler) Handle(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand) error {
    detector, err := clicontext.NewDetector()
    if err != nil {
        return err
    }

    execCtx, err := detector.DetectContext()
    if err != nil {
        return err
    }

    switch c := execCtx.(type) {
    case *clicontext.ProjectContext:
        return h.handleProjectContext(ctx, cmd, args, base, c)
    case *clicontext.SharedContext:
        return h.handleGlobalContext(ctx, cmd, args, base, c)
    default:
        return fmt.Errorf("unknown context type: %T", execCtx)
    }
}
```

**Update method signatures:**
```go
// Before
func (h *UpHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error

// After
func (h *UpHandler) handleProjectContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ProjectContext) error

// Before
func (h *UpHandler) handleGlobalContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.ExecutionContext) error

// After
func (h *UpHandler) handleGlobalContext(ctx context.Context, cmd *cobra.Command, args []string, base *base.BaseCommand, execCtx *clicontext.SharedContext) error
```

**Update internal usage:**
```go
// Before
execCtx.SharedContainers.Root

// After (ProjectContext)
execCtx.Shared.Root

// After (SharedContext)
execCtx.Shared.Root
```

**Tests:** Update `up_test.go`
```go
func TestUpHandler_Handle_ProjectContext(t *testing.T) {
    // Test with ProjectContext
}

func TestUpHandler_Handle_SharedContext(t *testing.T) {
    // Test with SharedContext
}
```

---

#### 3.2: Update DownHandler

**File:** `internal/pkg/cli/handlers/lifecycle/down.go`

**Changes:** Same pattern as UpHandler
- Replace `detector.Detect()` with `detector.DetectContext()`
- Replace `if execCtx.Type == clicontext.Shared` with type switch
- Update method signatures to accept concrete types
- Update tests

---

#### 3.3: Update RestartHandler

**File:** `internal/pkg/cli/handlers/lifecycle/restart.go`

**Changes:** Same pattern as UpHandler

---

#### 3.4: Update StatusHandler

**File:** `internal/pkg/cli/handlers/operations/status.go`

**Special case:** Has additional flag logic

**Before:**
```go
if execCtx.Type == clicontext.Shared || showAll || showShared {
    return h.handleSharedStatus(ctx, cmd, args, base, execCtx)
}
```

**After:**
```go
switch c := execCtx.(type) {
case *clicontext.ProjectContext:
    if showAll || showShared {
        return h.handleSharedStatus(ctx, cmd, args, base, c.Shared)
    }
    return h.handleProjectStatus(ctx, cmd, args, base, c)
case *clicontext.SharedContext:
    return h.handleSharedStatus(ctx, cmd, args, base, c.Shared)
default:
    return fmt.Errorf("unknown context type: %T", execCtx)
}
```

---

#### 3.5: Update ConnectHandler

**File:** `internal/pkg/cli/handlers/operations/connect.go`

**Changes:** Same pattern as UpHandler

---

#### 3.6: Update ExecHandler

**File:** `internal/pkg/cli/handlers/operations/exec.go`

**Changes:** Same pattern as UpHandler

---

#### 3.7: Update LogsHandler

**File:** `internal/pkg/cli/handlers/operations/logs.go`

**Changes:** Same pattern as UpHandler

---

### Step 4: Update All Tests

**Files to update:**
- `internal/pkg/cli/handlers/lifecycle/up_test.go`
- `internal/pkg/cli/handlers/lifecycle/down_test.go`
- `internal/pkg/cli/handlers/lifecycle/restart_test.go`
- `internal/pkg/cli/handlers/operations/status_test.go`
- `internal/pkg/cli/handlers/operations/connect_test.go`
- `internal/pkg/cli/handlers/operations/exec_test.go`
- `internal/pkg/cli/handlers/operations/logs_test.go`

**Pattern:**
```go
// Before
execCtx := &clicontext.ExecutionContext{
    Type: clicontext.Project,
    Project: &clicontext.ProjectInfo{...},
    SharedContainers: &clicontext.SharedInfo{...},
}

// After
execCtx := &clicontext.ProjectContext{
    Project: &clicontext.ProjectInfo{...},
    Shared: &clicontext.SharedInfo{...},
}
```

---

### Step 5: Remove Old Code (Breaking Change)

**File:** `internal/pkg/cli/context/types.go`

**Action:** Remove deprecated types

```go
// REMOVE: ContextType enum
// type ContextType string
// const (
//     Project ContextType = "project"
//     Shared  ContextType = "shared"
// )

// REMOVE: ExecutionContext struct
// type ExecutionContext struct {
//     Type             ContextType
//     Project          *ProjectInfo
//     SharedContainers *SharedInfo
// }
```

**File:** `internal/pkg/cli/context/detector.go`

**Action:** Remove old Detect() method

```go
// REMOVE: Old Detect() method
// func (d *Detector) Detect() (*ExecutionContext, error) { ... }

// RENAME: DetectContext() -> Detect()
func (d *Detector) Detect() (Context, error) {
    // Implementation from DetectContext()
}
```

---

### Step 6: Run Full Test Suite

```bash
# Unit tests
task test-unit

# Integration tests
task test-integration

# E2E tests
task test-e2e

# Build verification
go build ./...

# Generate code
task generate
```

---

### Step 7: Update Documentation

**Files to update:**
- `README.md` - If it mentions context types
- `docs/architecture.md` - Update architecture diagrams
- `docs/development.md` - Update development guide

---

## Verification Checklist

- [ ] All new types implement Context interface
- [ ] All 7 handlers updated to use type switches
- [ ] All handler tests updated
- [ ] Context detector tests updated
- [ ] No references to old ContextType enum
- [ ] No references to ExecutionContext.Type field
- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] Build succeeds
- [ ] Code generation succeeds
- [ ] No breaking changes in CLI interface
- [ ] Documentation updated

---

## Rollback Plan

If issues are discovered:

1. **Immediate:** Revert the PR
2. **Investigation:** Identify root cause
3. **Fix:** Address issues in feature branch
4. **Re-test:** Full test suite
5. **Re-deploy:** Create new PR

**Rollback commit message:**
```
Revert "feat: refactor ExecutionContext to interface-based types"

This reverts commit <hash>.

Reason: <specific issue>
```

---

## Performance Considerations

**Expected Impact:** Negligible
- Type switches are compile-time optimized
- No additional allocations
- Same memory footprint

**Benchmark:** Add benchmarks to verify
```go
func BenchmarkContextTypeSwitch(b *testing.B) {
    ctx := &ProjectContext{...}
    for i := 0; i < b.N; i++ {
        switch ctx.(type) {
        case *ProjectContext:
            // no-op
        case *SharedContext:
            // no-op
        }
    }
}
```

---

## Migration Timeline

**Day 1-2:** Steps 1-2 (New types, no breaking changes)
**Day 3-5:** Step 3 (Update handlers)
**Day 6:** Step 4 (Update tests)
**Day 7:** Steps 5-6 (Remove old code, verification)
**Day 8:** Step 7 (Documentation)

**Total:** ~8 working days

---

## Success Criteria

1. ✅ Zero runtime type assertion failures
2. ✅ All tests passing (unit, integration, e2e)
3. ✅ No performance degradation
4. ✅ Code coverage maintained or improved
5. ✅ No breaking changes in CLI interface
6. ✅ Clean git history with atomic commits
7. ✅ Documentation updated
8. ✅ Team review approved
