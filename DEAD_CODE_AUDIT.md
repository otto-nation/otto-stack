# Dead Code Audit

This document tracks unreachable functions identified by the deadcode tool and explains why they should be kept.

## Summary

**Current Count:** 8 unreachable functions  
**Status:** All legitimate - no action needed

## Unreachable Functions

### 1. Dependency Injection Constructors (2)

These constructors are used for dependency injection in tests, allowing mock dependencies to be injected.

- `internal/core/docker/client.go:79` - **NewClientWithDependencies**
  - Used in unit tests to inject mock dependencies
  - Essential for testability and isolation
  
- `internal/pkg/services/service_manager.go:35` - **NewServiceWithDependencies**
  - Used in unit tests to inject mock dependencies
  - Essential for testability and isolation

### 2. Test Infrastructure (4)

Test helper functions that appear unreachable because deadcode only analyzes production code paths.

- `internal/pkg/cli/handlers/project/test_helpers.go:12` - **setupTestDir**
  - Helper function for setting up test directories
  - Used across multiple test files
  
- `internal/pkg/cli/handlers/project/test_helpers.go:29` - **createTestFile**
  - Helper function for creating test files
  - Used across multiple test files
  
- `internal/pkg/cli/handlers/project/test_helpers.go:35` - **createTestConfig**
  - Helper function for creating test configurations
  - Used across multiple test files

- `internal/core/docker/client.go:120` - **Client.GetComposeManager**
  - Getter method used in unit tests
  - Provides access to internal state for testing

### 3. E2E Test Usage (1)

Functions used in end-to-end tests which are not analyzed by deadcode.

- `internal/core/docker/labels.go:21` - **Client.ListProjectContainers**
  - Used in e2e tests to verify container state
  - Critical for integration testing

### 4. Code Generation (1)

Functions called by code generation tools, not by runtime code.

- `internal/pkg/config/config.go:76` - **LoadCommandConfig**
  - Used by `cmd/generate-cli/main.go` to load command definitions
  - Essential for CLI code generation from YAML

## Historical Progress

### Phase 1: Initial Cleanup
- **Before:** 52 unreachable functions
- **After:** 11 unreachable functions
- **Removed:** 41 functions (~79% reduction)

### Phase 2: Safe Cleanup
- **Before:** 11 unreachable functions
- **After:** 11 unreachable functions (test utilities restored)
- **Actions:** Removed 41 functions, ~750+ lines, 10 test files
- **Fixed:** CI validation for docs generation

### Phase 3: Validation Integration
- **Before:** 11 unreachable functions
- **After:** 9 unreachable functions
- **Actions:** Integrated CheckInitialization and ValidateUpFlags into handlers

### Phase 4: Service Layer Refactoring
- **Before:** 9 unreachable functions
- **After:** 8 unreachable functions
- **Actions:** Refactored logs and exec handlers to use Service.Logs() and Service.Exec()

## Conclusion

All 8 remaining unreachable functions serve legitimate purposes:
- **2** are DI constructors for testing
- **4** are test infrastructure helpers
- **1** is used in e2e tests
- **1** is used by code generation

**No further cleanup recommended.** These functions are essential for maintaining test quality, code generation, and architectural patterns.
