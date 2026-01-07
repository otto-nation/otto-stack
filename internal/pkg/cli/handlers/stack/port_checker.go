package stack

// PortConflict represents a port conflict
type PortConflict struct {
	Port        string
	ServiceName string
	ProcessName string
	PID         string
}

// Legacy port checking functions were removed after command pattern migration
// Port conflict checking logic should be moved to UpCommand.Execute if needed
