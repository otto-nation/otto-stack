package services

// ValidationContext provides context for service validation
type ValidationContext struct {
	AllowHidden     bool
	IsUserRequested bool
	IsDependency    bool
}

// NewUserValidationContext creates context for user-requested services
func NewUserValidationContext() ValidationContext {
	return ValidationContext{
		AllowHidden:     false,
		IsUserRequested: true,
		IsDependency:    false,
	}
}

// NewDependencyValidationContext creates context for dependency services
func NewDependencyValidationContext() ValidationContext {
	return ValidationContext{
		AllowHidden:     true,
		IsUserRequested: false,
		IsDependency:    true,
	}
}

// NewInternalValidationContext creates context for internal operations
func NewInternalValidationContext() ValidationContext {
	return ValidationContext{
		AllowHidden:     true,
		IsUserRequested: false,
		IsDependency:    false,
	}
}
