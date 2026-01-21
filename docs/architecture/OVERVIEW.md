# Architecture Overview

Architecture documentation for otto-stack developers.

## Documents

### [CLI Architecture](CLI.md)
Overview of the CLI architecture, including:
- Command flow and layers
- Handler pattern
- Middleware chain
- Context detection
- Configuration system
- Error handling
- Testing strategy

### [Shared Containers Architecture](SHARED_CONTAINERS.md)
Deep dive into shared containers feature:
- Component architecture
- Registry management
- Container naming strategy
- Lifecycle integration
- Orphan detection
- Data flow diagrams
- Implementation details

## Quick Reference

### Key Architectural Decisions

1. **Simple Conditionals Over Strategy Pattern**
   - Use if/else for 2-case scenarios (project/global context)
   - Strategy pattern only when 3+ variations exist
   - YAGNI principle: don't over-engineer

2. **Function Decomposition**
   - Break handlers into focused functions
   - Each function has single responsibility
   - Improves readability and testability

3. **Declarative Configuration**
   - Commands defined in YAML (commands.yaml)
   - Messages centralized (messages.yaml)
   - Schema-driven validation (schema.yaml)

4. **Context-Aware Behavior**
   - Detect project vs global context
   - Adapt command behavior accordingly
   - Consistent user experience

5. **Registry for Safety**
   - Track shared container usage
   - Prevent accidental stops
   - Enable orphan detection

### Architecture Principles

- **Separation of Concerns**: Clear layer boundaries
- **Single Responsibility**: Each component has one job
- **YAGNI**: Don't add complexity until needed
- **Fail Fast**: Validate early, fail with clear messages
- **User Safety**: Prompt before destructive operations

### Code Organization

```
internal/
├── core/              # Types, constants
├── pkg/
│   ├── cli/          # CLI layer
│   ├── config/       # Configuration
│   ├── registry/     # Registry management
│   └── compose/      # Docker Compose wrapper
└── config/           # YAML definitions
```

### Testing Approach

- **Unit Tests**: Test components in isolation
- **Integration Tests**: Test component interactions
- **E2E Tests**: Test full user workflows
- **Manual Tests**: Real-world scenarios

## Contributing

When adding new features:

1. **Start Simple**: Use straightforward implementations
2. **Add Complexity Only When Needed**: Don't anticipate future requirements
3. **Document Decisions**: Explain why, not just what
4. **Test Thoroughly**: Unit, integration, and E2E tests
5. **Update Architecture Docs**: Keep documentation current

## See Also

- [User Guide](../SHARED_CONTAINERS.md)
- [Configuration Guide](../../docs-site/content/configuration.md)
- [CLI Reference](../../docs-site/content/cli-reference.md)
- [Contributing Guide](../CONTRIBUTING.md)
