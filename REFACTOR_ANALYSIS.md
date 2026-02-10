# Issue #112: Go-Idiomatic Interface-Based Type Discrimination Analysis

## Executive Summary

The codebase uses **7 string-based type enums** across 21,618 lines of Go code. The primary anti-pattern is `ExecutionContext.Type` which creates redundant state and prevents compile-time safety. This analysis provides a complete refactoring strategy organized by priority and impact.

---

## Current Anti-Patterns Identified

### 1. **ExecutionContext.Type** (HIGH PRIORITY - Core Issue)
**Location:** `internal/pkg/cli/context/types.go`

**Problem:**
```go
type ContextType string
const (
    Project ContextType = "project"
    Shared  ContextType = "shared"
)

type ExecutionContext struct {
    Type             ContextType
    Project          *ProjectInfo // nil if Shared
    SharedContainers *SharedInfo
}
```

**Issues:**
- Redundant state: `Type == Shared` duplicates `Project == nil`
- Not type-safe: Can create invalid states (Type=Project but Project=nil)
- Runtime checks instead of compile-time guarantees
- Used in **7 handler files**

**Usage Pattern:**
```go
if execCtx.Type == clicontext.Shared {
    return h.handleGlobalContext(...)
}
return h.handleProjectContext(...)
```

**Files Affected:**
- `internal/pkg/cli/handlers/lifecycle/up.go`
- `internal/pkg/cli/handlers/lifecycle/down.go`
- `internal/pkg/cli/handlers/lifecycle/restart.go`
- `internal/pkg/cli/handlers/operations/status.go`
- `internal/pkg/cli/handlers/operations/connect.go`
- `internal/pkg/cli/handlers/operations/exec.go`
- `internal/pkg/cli/handlers/operations/logs.go`

---

### 2. **ServiceType** (MEDIUM PRIORITY)
**Location:** `internal/pkg/types/service_types.go`

**Current:**
```go
type ServiceType string
const (
    ServiceTypeContainer     ServiceType = "container"
    ServiceTypeComposite     ServiceType = "composite"
    ServiceTypeConfiguration ServiceType = "configuration"
)
```

**Usage:** 2 locations
- `cmd/generate-services/main.go:305` - Type check for container services
- `internal/pkg/compose/generator.go:96` - Skip configuration services

**Impact:** Medium - Used in service generation and compose file creation

---

### 3. **ResourceType** (MEDIUM PRIORITY)
**Location:** `internal/core/docker/client.go`

**Current:**
```go
type ResourceType string
const (
    ResourceContainer ResourceType = "container"
    ResourceVolume    ResourceType = "volume"
    ResourceNetwork   ResourceType = "network"
    ResourceImage     ResourceType = "image"
)
```

**Usage:** `internal/core/docker/resource_manager.go`
- Two switch statements for List() and Remove() operations
- Each resource type has different Docker API calls

**Impact:** Medium - Core Docker operations, but well-isolated

---

### 4. **State/Health Enums** (LOW PRIORITY - Technical Debt)

#### ServiceState & HealthStatus
**Location:** `internal/pkg/types/types.go`
```go
type ServiceState string
const (
    ServiceStateRunning  ServiceState = "running"
    ServiceStateStopped  ServiceState = "stopped"
    ServiceStateStarting ServiceState = "starting"
)

type HealthStatus string
const (
    HealthStatusHealthy   HealthStatus = "healthy"
    HealthStatusUnhealthy HealthStatus = "unhealthy"
    HealthStatusStarting  HealthStatus = "starting"
    HealthStatusNone      HealthStatus = "none"
)
```

#### DockerServiceState & DockerHealthStatus
**Location:** `internal/core/docker/types.go`
```go
type DockerServiceState string
const (
    DockerServiceStateRunning DockerServiceState = "running"
    DockerServiceStateStopped DockerServiceState = "stopped"
    DockerServiceStateCreated DockerServiceState = "created"
)

type DockerHealthStatus string
const (
    DockerHealthStatusHealthy   DockerHealthStatus = "healthy"
    DockerHealthStatusUnhealthy DockerHealthStatus = "unhealthy"
    DockerHealthStatusStarting  DockerHealthStatus = "starting"
    DockerHealthStatusNone      DockerHealthStatus = "none"
)
```

**Usage:** Primarily in `internal/pkg/display/status_formatter.go` for icon selection

**Impact:** Low - These are display/status values, less critical for type safety

---

### 5. **Other String Enums** (LOW PRIORITY)

#### RestartPolicy
```go
type RestartPolicy string
const (
    RestartPolicyNo            RestartPolicy = "no"
    RestartPolicyAlways        RestartPolicy = "always"
    RestartPolicyOnFailure     RestartPolicy = "on-failure"
    RestartPolicyUnlessStopped RestartPolicy = "unless-stopped"
)
```

#### ConnectionType
```go
type ConnectionType string
const (
    ConnectionTypeCLI ConnectionType = "cli"
)
```

#### ParameterType
```go
type ParameterType string
const (
    ParameterTypeString  ParameterType = "string"
    ParameterTypeInteger ParameterType = "integer"
)
```

#### ShellType
```go
type ShellType string
const (
    ShellTypeBash       ShellType = "bash"
    ShellTypeZsh        ShellType = "zsh"
    ShellTypeFish       ShellType = "fish"
    ShellTypePowershell ShellType = "powershell"
)
```

**Impact:** Very Low - Configuration values, rarely used in type discrimination

---

## Proposed Go-Idiomatic Solutions

### Phase 1: ExecutionContext Refactor (HIGH PRIORITY)

#### New Design
```go
// Context interface for type discrimination
type Context interface {
    SharedRoot() string
    isContext() // unexported marker method
}

// ProjectContext for project-scoped operations
type ProjectContext struct {
    Project *ProjectInfo
    Shared  *SharedInfo
}

func (p *ProjectContext) SharedRoot() string { return p.Shared.Root }
func (p *ProjectContext) isContext()         {}

// SharedContext for global shared container operations
type SharedContext struct {
    Shared *SharedInfo
}

func (s *SharedContext) SharedRoot() string { return s.Shared.Root }
func (s *SharedContext) isContext()         {}
```

#### Updated Detector
```go
func (d *Detector) Detect() (Context, error) {
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
```

#### Handler Usage (Type-Safe)
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

    // Type-safe discrimination
    switch ctx := execCtx.(type) {
    case *clicontext.ProjectContext:
        return h.handleProjectContext(ctx, cmd, args, base, ctx)
    case *clicontext.SharedContext:
        return h.handleGlobalContext(ctx, cmd, args, base, ctx)
    default:
        return fmt.Errorf("unknown context type")
    }
}
```

**Benefits:**
- Compile-time type safety
- Impossible to create invalid states
- Clear intent through type switches
- Follows Go stdlib patterns (io.Reader, error, context.Context)

---

### Phase 2: ServiceType Refactor (MEDIUM PRIORITY)

#### New Design
```go
// Service interface for type discrimination
type Service interface {
    Name() string
    isService() // unexported marker
}

// ContainerService represents a Docker container service
type ContainerService struct {
    ServiceConfig
}

func (c *ContainerService) isService() {}

// CompositeService represents a group of services
type CompositeService struct {
    ServiceConfig
    Components []string
}

func (c *CompositeService) isService() {}

// ConfigurationService represents configuration-only (no container)
type ConfigurationService struct {
    ServiceConfig
}

func (c *ConfigurationService) isService() {}
```

#### Usage in Generator
```go
func (g *Generator) processServiceConfigAndDependencies(config *types.ServiceConfig, ...) error {
    // Type-safe check
    if _, isConfig := config.(*types.ConfigurationService); isConfig {
        return nil // Skip configuration services
    }
    
    serviceConfig := g.buildService(config)
    if serviceConfig != nil {
        serviceList[config.Name] = serviceConfig
    }
    return nil
}
```

---

### Phase 3: ResourceType Refactor (MEDIUM PRIORITY)

#### New Design
```go
// Resource interface for Docker resources
type Resource interface {
    List(ctx context.Context, filter filters.Args) ([]string, error)
    Remove(ctx context.Context, names []string) error
    isResource() // unexported marker
}

// ContainerResource handles container operations
type ContainerResource struct {
    client *Client
}

func (r *ContainerResource) List(ctx context.Context, filter filters.Args) ([]string, error) {
    containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filter})
    if err != nil {
        return nil, err
    }
    names := make([]string, len(containers))
    for i, c := range containers {
        names[i] = c.ID
    }
    return names, nil
}

func (r *ContainerResource) Remove(ctx context.Context, names []string) error {
    // Implementation
}

func (r *ContainerResource) isResource() {}

// Similar for VolumeResource, NetworkResource, ImageResource
```

#### Usage
```go
func (rm *ResourceManager) GetResource(resourceType string) (Resource, error) {
    switch resourceType {
    case "container":
        return &ContainerResource{client: rm.client}, nil
    case "volume":
        return &VolumeResource{client: rm.client}, nil
    case "network":
        return &NetworkResource{client: rm.client}, nil
    case "image":
        return &ImageResource{client: rm.client}, nil
    default:
        return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
    }
}
```

---

## Implementation Strategy

### Phase 1: ExecutionContext (Week 1-2)
**Priority:** HIGH - Core functionality, affects 7 handlers

1. Create new `Context` interface and concrete types
2. Update `Detector.Detect()` to return interface
3. Update 7 handlers to use type switches
4. Update all tests
5. Remove old `ContextType` enum
6. Verify all integration tests pass

**Estimated Impact:**
- Files changed: ~15
- Lines changed: ~200-300
- Test files: ~7

---

### Phase 2: ServiceType (Week 3)
**Priority:** MEDIUM - Service generation and compose

1. Create `Service` interface hierarchy
2. Update service config loading to return concrete types
3. Update compose generator
4. Update service generation tool
5. Update tests

**Estimated Impact:**
- Files changed: ~8
- Lines changed: ~150-200
- Test files: ~4

---

### Phase 3: ResourceType (Week 4)
**Priority:** MEDIUM - Docker operations

1. Create `Resource` interface
2. Implement concrete resource types
3. Update ResourceManager
4. Update tests

**Estimated Impact:**
- Files changed: ~5
- Lines changed: ~200-250
- Test files: ~3

---

### Phase 4: State/Health Enums (Future - Technical Debt)
**Priority:** LOW - Display/status values

These are less critical as they're primarily used for display purposes and don't involve complex type discrimination logic. Consider refactoring only if:
- Adding new state types becomes frequent
- Type safety issues emerge
- Display logic becomes more complex

---

## Risk Assessment

### High Risk Areas
1. **ExecutionContext** - Core to all CLI operations
   - Mitigation: Comprehensive integration tests
   - Rollback plan: Keep old code commented for 1 release

2. **ServiceType** - Affects service generation
   - Mitigation: Test with all service types
   - Validation: Compare generated compose files before/after

### Medium Risk Areas
3. **ResourceType** - Docker operations
   - Mitigation: Well-isolated, good test coverage
   - Validation: Test all resource operations

### Low Risk Areas
4. **State/Health enums** - Display only
   - Minimal risk, can be done incrementally

---

## Testing Strategy

### Unit Tests
- Test each concrete type implements interface
- Test type switches handle all cases
- Test invalid state prevention

### Integration Tests
- Test all CLI commands with both context types
- Test service generation with all service types
- Test Docker operations with all resource types

### Regression Tests
- Compare behavior before/after refactor
- Ensure no breaking changes in CLI interface
- Validate generated files are identical

---

## Success Metrics

1. **Type Safety:** Zero runtime type assertion failures
2. **Code Quality:** Reduced cyclomatic complexity in handlers
3. **Maintainability:** Easier to add new context/service/resource types
4. **Performance:** No measurable performance degradation
5. **Test Coverage:** Maintain or improve current coverage (90%+)

---

## Dependencies & Blockers

### Prerequisites
- Issue #102 (terminology standardization) - **Should be completed first**
- All current PRs merged to avoid conflicts

### Blockers
- None identified - refactoring is self-contained

---

## Rollout Plan

### Stage 1: Feature Branch Development
- Create `feat/interface-based-types` branch
- Implement Phase 1 (ExecutionContext)
- Full test suite passing

### Stage 2: Alpha Testing
- Deploy to internal test environment
- Run full integration test suite
- Manual testing of all CLI commands

### Stage 3: Code Review
- Detailed review by team
- Address feedback
- Performance benchmarking

### Stage 4: Merge & Release
- Merge to main
- Include in next minor version release
- Monitor for issues

### Stage 5: Follow-up Phases
- Implement Phase 2 & 3 in subsequent releases
- Phase 4 as technical debt cleanup

---

## Conclusion

The refactoring to interface-based type discrimination will:
- ✅ Eliminate redundant state
- ✅ Provide compile-time type safety
- ✅ Follow Go idioms and best practices
- ✅ Make the codebase more maintainable
- ✅ Prevent entire classes of bugs

**Recommendation:** Proceed with Phase 1 (ExecutionContext) immediately as it provides the highest value and addresses the core anti-pattern identified in issue #112.
