# CLI Architecture

## Layers

```
CLI Entry → Commands → Middleware → Handlers → Business Logic → Infrastructure
```

## Command Flow

```
User Input → Cobra Parser → Middleware Chain → Handler → Output
                              ├─ Validation
                              ├─ Context Detection
                              └─ Config Loading
```

## Components

### Commands (`internal/config/commands.yaml`)
Declarative command definitions with flags and help text.

### Middleware (`internal/pkg/cli/middleware/`)
Pre-processes requests: validation, context detection, config loading.

### Handlers (`internal/pkg/cli/handlers/`)
- **Lifecycle**: up, down, cleanup
- **Operations**: status, logs, exec  
- **Project**: init, validate

### Context Detection (`internal/pkg/cli/context/`)
Determines project vs global context by checking for `.otto-stack/config.yaml`.

### Registry (`internal/pkg/registry/`)
Tracks shared container usage across projects.

### Config System (`internal/pkg/config/`)
Multi-source config: CLI flags → Project config → Service configs → Defaults

## Design Patterns

### Simple Conditionals (Not Strategy)
Use if/else for 2-case scenarios (project/global context). Strategy pattern only when 3+ variations exist.

```go
func (h *UpHandler) Handle(ctx *context.ExecutionContext) error {
    if ctx.Type == context.ProjectContext {
        return h.handleProjectContext(ctx)
    }
    return h.handleGlobalContext(ctx)
}
```

### Function Decomposition
Break handlers into focused functions with single responsibilities.

### Middleware Chain
Pre-process requests before handlers for validation and setup.

## Error Handling

Centralized messages in `internal/config/messages.yaml`:
```yaml
errors:
  no_project_found:
    message: "No otto-stack project found"
    hint: "Run 'otto-stack init' to create a project"
```

## Code Organization

```
internal/
├── core/              # Types, constants
├── pkg/
│   ├── cli/          # Commands, handlers, middleware
│   ├── config/       # Configuration
│   ├── registry/     # Registry management
│   └── compose/      # Docker Compose wrapper
└── config/           # YAML definitions
```

## See Also
- [Shared Containers Architecture](SHARED_CONTAINERS.md)
