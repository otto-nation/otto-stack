package cli

// ShellType represents supported shell types for completion
type ShellType string

const (
	ShellTypeBash       ShellType = "bash"
	ShellTypeZsh        ShellType = "zsh"
	ShellTypeFish       ShellType = "fish"
	ShellTypePowerShell ShellType = "powershell"
)
