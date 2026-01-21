# Dead Code Audit Results

## Summary
- **Total unreachable functions**: 172
- **Focus**: Production code in `internal/pkg` (excluding test helpers and e2e framework)

## Dead Code by Package

### High Priority - Display Package (21 functions)
**Status**: Entire old formatter system unused after go-pretty migration

#### formatter.go (9 functions) - REMOVE ENTIRE FILE
- `New()` - Old formatter constructor
- `Formatter.initHandlers()` 
- `Formatter.FormatStatus()`
- `Formatter.FormatServiceCatalog()`
- `Formatter.FormatValidation()`
- `Formatter.FormatVersion()`
- `Formatter.FormatHealth()`
- `JSONHandler.Handle()`
- `YAMLHandler.Handle()`

#### catalog_formatter.go (4 functions) - KEEP (used in tests)
- `NewCatalogFormatter()` - Used in Formatter struct (which is dead)
- `CatalogFormatter.FormatTable()` - Unused
- `CatalogFormatter.FormatGrouped()` - Unused  
- `FilterCatalogByCategory()` - Used in tests

**Decision**: Remove formatter.go entirely. Keep catalog_formatter.go for now (used in tests).

#### validation_formatter.go (8 functions) - REMOVE ENTIRE FILE
- `NewValidationFormatter()`
- `ValidationFormatter.FormatTable()`
- `NewHealthFormatter()`
- `HealthFormatter.FormatTable()`
- `HealthFormatter.formatHealthSummary()`
- `HealthFormatter.getHealthIcon()`
- `NewVersionFormatter()`
- `VersionFormatter.FormatTable()`

### Medium Priority - Compose Package (6 functions)
#### naming.go (6 functions) - KEEP (may be needed for shared containers)
- `NewNamingStrategy()`
- `NamingStrategy.ContainerName()`
- `NamingStrategy.VolumeName()`
- `NamingStrategy.NetworkName()`
- `NamingStrategy.IsShared()`
- `NamingStrategy.isShared()`

**Decision**: Keep - these are for shared container feature, tested but not yet used in production.

### Medium Priority - Config Package (5 functions)
#### config.go (5 functions)
- `LoadServiceConfig()` - Unused
- `LoadCommandConfig()` - Unused (we use LoadCommandConfigStruct)
- `getServiceConfigDir()` - Unused
- `loadServiceConfigFile()` - Unused
- `mergeServiceConfigs()` - Unused

**Decision**: Remove unused config loaders.

### Low Priority - Other Packages

#### CLI Context (6 functions) - KEEP
- Context detection and types - may be needed for future features

#### Services (6 functions)
- `NewServiceWithDependencies()` - Used in tests
- `Service.Logs()` - Implemented but not exposed via CLI yet
- `Service.Exec()` - Implemented but not exposed via CLI yet
- `ExtractVisibleServiceNames()` - Unused
- `IsYAMLFile()` - Unused
- `TrimYAMLExt()` - Unused

**Decision**: Remove utility functions. Keep Logs/Exec for future CLI commands.

#### Version (7 functions) - KEEP
- Version parsing and constraints - may be needed for update checker

#### Validation (2 functions) - KEEP
- `CheckInitialization()` - Used by handlers
- `ValidateUpFlags()` - Used by handlers

#### Errors (2 functions)
- `NewServiceErrorf()` - Unused variant
- `NewDockerErrorf()` - Unused variant

**Decision**: Remove unused error constructors.

#### Filesystem (2 functions)
- `CopyFile()` - Unused
- `ExpandPath()` - Unused

**Decision**: Remove unused filesystem utilities.

#### Logger (2 functions)
- `Warn()` - Unused
- `With()` - Unused

**Decision**: Keep - standard logger methods.

## Cleanup Plan

### Phase 1: Display Package Cleanup (Issue #62)
1. Remove `internal/pkg/display/formatter.go` (old formatter system)
2. Remove `internal/pkg/display/validation_formatter.go` (old validation formatters)
3. Update tests to remove references
4. Verify handlers still work with go-pretty tables

### Phase 2: Config Package Cleanup
1. Remove unused config loaders from `config.go`
2. Verify no hidden dependencies

### Phase 3: Utilities Cleanup
1. Remove unused service utilities
2. Remove unused error constructors
3. Remove unused filesystem functions

### Phase 4: Test/E2E Cleanup (Separate PR)
- Clean up test helpers and e2e framework (56 functions)
- Keep for now as they may be useful for future tests

## Files to Delete
- `internal/pkg/display/formatter.go` (entire file)
- `internal/pkg/display/validation_formatter.go` (entire file)

## Functions to Remove
- Config: 5 functions from config.go
- Services: 3 utility functions from utils.go
- Errors: 2 error constructors from types.go
- Filesystem: 2 functions from operations.go

## Expected Impact
- **Lines removed**: ~500-600 lines
- **Files deleted**: 2
- **Test updates needed**: Yes (remove formatter tests)
- **Risk**: Low (all dead code confirmed by static analysis)
