package handlers

// Package names from YAML
const (
	PackageCore     = "core"
	PackageServices = "services"
	PackageInit     = "init"
	PackageStandard = "standard"
)

// Command mappings
var CoreCommands = map[string]func() any{
	"up":      func() any { return "NewUpHandler" },
	"down":    func() any { return "NewDownHandler" },
	"restart": func() any { return "NewRestartHandler" },
	"status":  func() any { return "NewStatusHandler" },
	"logs":    func() any { return "NewLogsHandler" },
	"exec":    func() any { return "NewExecHandler" },
	"connect": func() any { return "NewConnectHandler" },
	"cleanup": func() any { return "NewCleanupHandler" },
}

var ServicesCommands = map[string]func() any{
	"services":  func() any { return "NewServicesHandler" },
	"deps":      func() any { return "NewDepsHandler" },
	"conflicts": func() any { return "NewConflictsHandler" },
}

var StandardCommands = map[string]func() any{
	"doctor":     func() any { return "NewDoctorHandler" },
	"completion": func() any { return "NewCompletionHandler" },
	"validate":   func() any { return "NewValidateHandler" },
	"version":    func() any { return "NewVersionHandler" },
}
