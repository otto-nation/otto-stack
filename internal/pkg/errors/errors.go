package errors

import (
	"fmt"
)

// Error codes
const (
	ErrCodeUnknown       = "UNKNOWN"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeAlreadyExists = "ALREADY_EXISTS"
	ErrCodeInvalid       = "INVALID"
	ErrCodePermission    = "PERMISSION_DENIED"
	ErrCodeTimeout       = "TIMEOUT"
	ErrCodeUnavailable   = "UNAVAILABLE"
	ErrCodeInternal      = "INTERNAL"
	ErrCodeOperationFail = "OPERATION_FAILED"
)

// Exit codes
const (
	ExitSuccess         = 0
	ExitGeneralError    = 1
	ExitNotFound        = 2
	ExitInvalidInput    = 3
	ExitPermissionError = 4
	ExitAlreadyExists   = 5
	ExitTimeout         = 6
	ExitUnavailable     = 7
)

// Common components
const (
	ComponentDocker   = "docker"
	ComponentService  = "service"
	ComponentServices = "services"
	ComponentConfig   = "config"
	ComponentStack    = "stack"
	ComponentRegistry = "registry"
	ComponentProject  = "project"
)

// Common fields
const (
	FieldFlags       = "flags"
	FieldProjectName = "project-name"
	FieldProjectPath = "project-path"
	FieldServiceName = "service-name"
)

// Error Context Guidelines
//
// Use the appropriate constructor based on the error type:
//
// 1. USER INPUT ERRORS - Use NewValidationError with field name constants:
//    NewValidationError(code, FieldProjectName, message, nil)
//    NewValidationError(code, FieldServiceName, message, nil)
//    NewValidationErrorf(code, FieldFlags, format, args...)
//
// 2. SYSTEM/ENVIRONMENT ERRORS - Use NewSystemError (no context parameter):
//    NewSystemError(code, message, nil)
//    NewSystemErrorf(code, format, args...)
//
// Examples of system errors:
//    - Docker not available
//    - File conflicts
//    - Permission issues
//    - Network errors
//    - Version parsing failures
//    - Config validation failures
//
// Examples:
//   Good: NewValidationError(ErrCodeInvalid, FieldProjectName, "project name too short", nil)
//   Good: NewSystemError(ErrCodeInvalid, "Docker not available", nil)
//   Bad:  NewValidationError(ErrCodeInvalid, "", "Docker not available", nil)  // Use NewSystemError instead
//   Bad:  NewValidationError(ErrCodeInvalid, "my-project", "invalid name", nil)  // Don't use actual values

// Error represents any application error
type Error struct {
	Code    string
	Context string // service name, file path, field name, etc.
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a validation error for user input fields
// Use when validating user-provided data like project names, service names, flags
func NewValidationError(code, field, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Context: field,
		Message: message,
		Cause:   cause,
	}
}

// NewValidationErrorf creates a validation error with formatted message
func NewValidationErrorf(code, field, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Context: field,
		Message: fmt.Sprintf(format, args...),
		Cause:   nil,
	}
}

// NewSystemError creates an error for system/environment issues
// Use for Docker availability, file conflicts, version parsing, etc.
func NewSystemError(code, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Context: "",
		Message: message,
		Cause:   cause,
	}
}

// NewSystemErrorf creates a system error with formatted message
func NewSystemErrorf(code, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Context: "",
		Message: fmt.Sprintf(format, args...),
		Cause:   nil,
	}
}

// NewServiceError creates an error for service operations
func NewServiceError(code, serviceName, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Context: serviceName,
		Message: message,
		Cause:   cause,
	}
}

// NewConfigError creates an error for configuration issues
func NewConfigError(code, configPath, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Context: configPath,
		Message: message,
		Cause:   cause,
	}
}

// NewConfigErrorf creates a config error with formatted message
func NewConfigErrorf(code, configPath, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Context: configPath,
		Message: fmt.Sprintf(format, args...),
		Cause:   nil,
	}
}

// NewDockerError creates an error for Docker operations
func NewDockerError(code, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Context: ComponentDocker,
		Message: message,
		Cause:   cause,
	}
}
