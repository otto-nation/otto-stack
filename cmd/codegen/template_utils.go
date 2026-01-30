package codegen

import (
	"os"
	"path/filepath"
	"text/template"

	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
)

const DirPermissions = 0750

// ParseTemplate parses a template file with common functions
func ParseTemplate(templatePath, templateName string) (*template.Template, error) {
	tmpl, err := template.New(templateName).Funcs(template.FuncMap{
		"toPascalCase": ToPascalCase,
	}).ParseFiles(templatePath)
	if err != nil {
		return nil, pkgerrors.NewServiceError("generator", "parse template", err)
	}

	return tmpl.Lookup(filepath.Base(templatePath)), nil
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, DirPermissions)
}

// TemplateExecutor handles template parsing and execution
type TemplateExecutor struct {
	templatePath string
	outputPath   string
}

// NewTemplateExecutor creates a new template executor
func NewTemplateExecutor(templatePath, outputPath string) *TemplateExecutor {
	return &TemplateExecutor{
		templatePath: templatePath,
		outputPath:   outputPath,
	}
}

// ExecuteTemplateWithFuncs parses and executes a template with custom functions
func (te *TemplateExecutor) ExecuteTemplateWithFuncs(data any, funcMap template.FuncMap) error {
	tmpl, err := template.New(filepath.Base(te.templatePath)).Funcs(funcMap).ParseFiles(te.templatePath)
	if err != nil {
		return pkgerrors.NewConfigError("template", "failed to parse template with functions", err)
	}

	file, err := os.Create(te.outputPath)
	if err != nil {
		return pkgerrors.NewConfigError("output", "failed to create output file", err)
	}
	defer func() { _ = file.Close() }()

	if err := tmpl.Execute(file, data); err != nil {
		return pkgerrors.NewConfigError("execution", "failed to execute template", err)
	}

	return nil
}

// ServiceFileData represents data for generating a service-specific file
type ServiceFileData struct {
	ServiceName string
	StructName  string
	FileName    string
	Schema      map[string]any
}
