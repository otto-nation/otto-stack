# Shared Containers Architecture

## Components

```
CLI Commands → Context Detector → Registry Manager → Docker Compose
                                        ↓
                                  containers.yaml
```

### Registry Manager
`internal/pkg/registry/registry.go`

Tracks project → container relationships:
```go
type ContainerInfo struct {
    Name     string    // otto-stack-postgres
    Service  string    // postgres
    Projects []string  // [myapp, api-service]
}
```

Operations: `Register()`, `Unregister()`, `FindOrphans()`, `CleanOrphans()`

### Container Naming
- Shared: `otto-stack-{service}` (e.g., `otto-stack-postgres`)
- Project: `{project}-{service}` (e.g., `myapp-postgres`)

### Lifecycle Handlers

**up**: Register shared services → Start containers  
**down**: Prompt if shared by others → Stop → Unregister  
**cleanup**: Detect orphans → Suggest cleanup  
**cleanup --orphans**: Remove orphans from registry

## Key Flows

### Starting Shared Container
```
otto-stack up → Detect context → Filter shared services
→ Registry.Register(service, container, project) → docker compose up
```

### Stopping Shared Container
```
otto-stack down → Check registry → Prompt if used by others
→ docker compose down → Registry.Unregister(service, project)
```

### Orphan Detection
```
otto-stack cleanup → Registry.FindOrphans() → Report orphans
otto-stack cleanup --orphans → Registry.CleanOrphans()
```

## Configuration

### Project Config
```yaml
sharing:
  enabled: true
  services:
    postgres: true  # Optional overrides
```

### Registry File
`~/.otto-stack/shared/containers.yaml`
```yaml
shared_containers:
  postgres:
    name: otto-stack-postgres
    service: postgres
    projects: [myapp, api-service]
```

## Design Decisions

**Registry in user home**: Persists across project deletions  
**Track projects**: Prevent stopping containers used elsewhere  
**Prompt before stopping**: User awareness and control  
**Auto orphan detection**: Proactive cleanup suggestions

## See Also
- [User Guide](../SHARED_CONTAINERS.md)
- [CLI Architecture](CLI.md)
