package errors

import (
	"fmt"
)

// Core error field types - keep only the most common ones
const (
	FieldFlags       = "flags"
	FieldProjectName = "project-name"
	FieldProjectPath = "project-path"
	FieldServiceName = "service-name"
)

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string
	Message string
	Cause   error
}

func (e *ValidationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("validation failed for %s: %s: %v", e.Field, e.Message, e.Cause)
	}
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string, cause error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

// NewValidationErrorf creates a new validation error with formatted message
func NewValidationErrorf(field, format string, args ...any) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf(format, args...),
		Cause:   nil,
	}
}

// ServiceError represents a service operation failure
type ServiceError struct {
	Service string
	Action  string
	Cause   error
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("service %s failed to %s: %v", e.Service, e.Action, e.Cause)
}

func (e *ServiceError) Unwrap() error {
	return e.Cause
}

// NewServiceError creates a new service error
func NewServiceError(service, action string, cause error) *ServiceError {
	return &ServiceError{
		Service: service,
		Action:  action,
		Cause:   cause,
	}
}

// NewServiceErrorf creates a new service error with formatted action
func NewServiceErrorf(service, action, format string, args ...any) *ServiceError {
	return &ServiceError{
		Service: service,
		Action:  fmt.Sprintf(format, args...),
		Cause:   nil,
	}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Path    string
	Message string
	Cause   error
}

func (e *ConfigError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("config error in %s: %s: %v", e.Path, e.Message, e.Cause)
	}
	return fmt.Sprintf("config error: %s: %v", e.Message, e.Cause)
}

func (e *ConfigError) Unwrap() error {
	return e.Cause
}

// NewConfigError creates a new config error
func NewConfigError(path, message string, cause error) *ConfigError {
	return &ConfigError{
		Path:    path,
		Message: message,
		Cause:   cause,
	}
}

// NewConfigErrorf creates a new config error with formatted message
func NewConfigErrorf(path, format string, args ...any) *ConfigError {
	return &ConfigError{
		Path:    path,
		Message: fmt.Sprintf(format, args...),
		Cause:   nil,
	}
}

// DockerError represents a Docker operation failure
type DockerError struct {
	Operation string
	Container string
	Cause     error
}

func (e *DockerError) Error() string {
	if e.Container != "" {
		return fmt.Sprintf("docker %s failed for container %s: %v", e.Operation, e.Container, e.Cause)
	}
	return fmt.Sprintf("docker %s failed: %v", e.Operation, e.Cause)
}

func (e *DockerError) Unwrap() error {
	return e.Cause
}

// NewDockerError creates a new docker error
func NewDockerError(operation, container string, cause error) *DockerError {
	return &DockerError{
		Operation: operation,
		Container: container,
		Cause:     cause,
	}
}

// NewDockerErrorf creates a new docker error with formatted operation
func NewDockerErrorf(operation, container, format string, args ...any) *DockerError {
	return &DockerError{
		Operation: fmt.Sprintf(format, args...),
		Container: container,
		Cause:     nil,
	}
}
