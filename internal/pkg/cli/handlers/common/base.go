package common

// BaseHandler provides common implementations for handler interface methods
type BaseHandler struct{}

// ValidateArgs provides default validation (accepts any arguments)
func (h *BaseHandler) ValidateArgs(args []string) error {
	return nil
}

// GetRequiredFlags provides default implementation (no required flags)
func (h *BaseHandler) GetRequiredFlags() []string {
	return []string{}
}
