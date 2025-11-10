package cli

// ShellType represents supported shell types for completion
type ShellType string

const (
	ShellTypeBash       ShellType = "bash"
	ShellTypeZsh        ShellType = "zsh"
	ShellTypeFish       ShellType = "fish"
	ShellTypePowerShell ShellType = "powershell"
)

// IsValid returns true if the shell type is supported
func (s ShellType) IsValid() bool {
	switch s {
	case ShellTypeBash, ShellTypeZsh, ShellTypeFish, ShellTypePowerShell:
		return true
	default:
		return false
	}
}

// AllShellTypeStrings returns all supported shell types as strings
func AllShellTypeStrings() []string {
	return []string{
		string(ShellTypeBash),
		string(ShellTypeZsh),
		string(ShellTypeFish),
		string(ShellTypePowerShell),
	}
}
