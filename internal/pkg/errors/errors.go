package errors

import (
	"errors"
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

// Error represents any application error
type Error struct {
	Code    string
	Context string // service name, file path, field name, etc.
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Context != "" {
		if e.Cause != nil {
			return fmt.Sprintf("%s: %s: %v", e.Context, e.Message, e.Cause)
		}
		return fmt.Sprintf("%s: %s", e.Context, e.Message)
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// New creates a new error
func New(code, context, message string, cause error) *Error {
	return &Error{Code: code, Context: context, Message: message, Cause: cause}
}

// Newf creates a new error with formatted message
func Newf(code, context, format string, args ...any) *Error {
	return &Error{Code: code, Context: context, Message: fmt.Sprintf(format, args...), Cause: nil}
}

// Legacy aliases
type ValidationError = Error
type ServiceError = Error
type ConfigError = Error

func NewValidationError(code, field, message string, cause error) *Error {
	return New(code, field, message, cause)
}

func NewValidationErrorf(code, field, format string, args ...any) *Error {
	return Newf(code, field, format, args...)
}

func NewServiceError(code, service, action string, cause error) *Error {
	return New(code, service, action, cause)
}

func NewConfigError(code, path, message string, cause error) *Error {
	return New(code, path, message, cause)
}

func NewConfigErrorf(code, path, format string, args ...any) *Error {
	return Newf(code, path, format, args...)
}

// Helpers
func GetErrorCode(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ErrCodeUnknown
}

func IsNotFound(err error) bool {
	return GetErrorCode(err) == ErrCodeNotFound
}

func IsAlreadyExists(err error) bool {
	return GetErrorCode(err) == ErrCodeAlreadyExists
}

func IsRetryable(err error) bool {
	code := GetErrorCode(err)
	return code == ErrCodeTimeout || code == ErrCodeUnavailable
}

func IsPermissionDenied(err error) bool {
	return GetErrorCode(err) == ErrCodePermission
}

func IsInvalid(err error) bool {
	return GetErrorCode(err) == ErrCodeInvalid
}

func FormatForUser(err error) string {
	if err == nil {
		return ""
	}
	code := GetErrorCode(err)
	switch code {
	case ErrCodeNotFound:
		return fmt.Sprintf("%s\nHint: Check the name and try again", err.Error())
	case ErrCodeUnavailable:
		return fmt.Sprintf("%s\nHint: Ensure required services are running", err.Error())
	case ErrCodePermission:
		return fmt.Sprintf("%s\nHint: Check file permissions or run with appropriate privileges", err.Error())
	default:
		return err.Error()
	}
}

func GetExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}
	switch GetErrorCode(err) {
	case ErrCodeNotFound:
		return ExitNotFound
	case ErrCodeInvalid:
		return ExitInvalidInput
	case ErrCodePermission:
		return ExitPermissionError
	case ErrCodeAlreadyExists:
		return ExitAlreadyExists
	case ErrCodeTimeout:
		return ExitTimeout
	case ErrCodeUnavailable:
		return ExitUnavailable
	default:
		return ExitGeneralError
	}
}
