package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/otto-nation/otto-stack/internal/pkg/constants"
	"gopkg.in/yaml.v3"
)

const (
	TemplateFilePath  = "cmd/generate-services/templates/services.tmpl"
	GeneratedFilePath = "internal/pkg/constants/services_generated.go"
	ServicesDir       = "internal/config/services"
)

type constantData struct {
	Name  string
	Value string
}

type ServiceYAML struct {
	Name         string           `yaml:"name"`
	Description  string           `yaml:"description"`
	Category     string           `yaml:"category"`
	Type         string           `yaml:"type"`
	Tags         []string         `yaml:"tags"`
	Environment  map[string]any   `yaml:"environment"`
	Connection   ConnectionYAML   `yaml:"connection"`
	Dependencies DependenciesYAML `yaml:"dependencies"`
	Docker       map[string]any   `yaml:"docker"`
	HealthCheck  map[string]any   `yaml:"health_check"`
}

type ConnectionYAML struct {
	Client       string `yaml:"client"`
	DefaultUser  string `yaml:"default_user"`
	DefaultPort  int    `yaml:"default_port"`
	UserFlag     string `yaml:"user_flag"`
	HostFlag     string `yaml:"host_flag"`
	PortFlag     string `yaml:"port_flag"`
	DatabaseFlag string `yaml:"database_flag"`
}

type DependenciesYAML struct {
	Required  []string `yaml:"required"`
	Soft      []string `yaml:"soft"`
	Conflicts []string `yaml:"conflicts"`
	Provides  []string `yaml:"provides"`
}

type DockerYAML struct {
	Image   string   `yaml:"image"`
	Ports   []string `yaml:"ports"`
	Volumes []string `yaml:"volumes"`
}

type HealthCheckYAML struct {
	Endpoint string `yaml:"endpoint"`
	Enabled  bool   `yaml:"enabled"`
}

type ServiceConstants struct {
	Categories      []constantData
	ServiceTypes    []constantData
	Clients         []constantData
	Ports           []constantData
	Names           []constantData
	Images          []constantData
	DefaultUsers    []constantData
	ConnectionFlags []constantData
	EnvVars         []constantData
	HealthEndpoints []constantData
	Tags            []constantData
	Dependencies    []constantData
	Conflicts       []constantData
	Provides        []constantData
}

type collectors struct {
	categories      map[string]bool
	serviceTypes    map[string]bool
	clients         map[string]bool
	ports           map[string]int
	names           map[string]bool
	images          map[string]string
	defaultUsers    map[string]string
	connectionFlags map[string]string
	envVars         map[string]string
	healthEndpoints map[string]string
	tags            map[string]bool
	dependencies    map[string][]string
	conflicts       map[string][]string
	provides        map[string][]string
}

func main() {
	serviceConstants, err := extractServiceConstants()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to extract service constants: %v\n", err)
		os.Exit(1)
	}

	if err := generateConstants(serviceConstants); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate constants: %v\n", err)
		os.Exit(1)
	}

	printSummary(serviceConstants)
}

func extractServiceConstants() (*ServiceConstants, error) {
	collectors := newCollectors()

	err := filepath.Walk(ServicesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !constants.IsYAMLFile(path) {
			return err
		}

		service, err := parseServiceFile(path)
		if err != nil {
			return err
		}

		collectors.collect(service)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return collectors.toServiceConstants(), nil
}

func parseServiceFile(path string) (*ServiceYAML, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var service ServiceYAML
	if err := yaml.Unmarshal(data, &service); err != nil {
		return nil, err
	}

	return &service, nil
}

func newCollectors() *collectors {
	return &collectors{
		categories:      make(map[string]bool),
		serviceTypes:    make(map[string]bool),
		clients:         make(map[string]bool),
		ports:           make(map[string]int),
		names:           make(map[string]bool),
		images:          make(map[string]string),
		defaultUsers:    make(map[string]string),
		connectionFlags: make(map[string]string),
		envVars:         make(map[string]string),
		healthEndpoints: make(map[string]string),
		tags:            make(map[string]bool),
		dependencies:    make(map[string][]string),
		conflicts:       make(map[string][]string),
		provides:        make(map[string][]string),
	}
}

func (c *collectors) collect(service *ServiceYAML) {
	c.collectBasicInfo(service)
	c.collectConnectionInfo(service)
	c.collectDockerInfo(service)
	c.collectEnvironmentVars(service)
	c.collectHealthEndpoints(service)
	c.collectTags(service)
	c.collectDependencies(service)
}

func (c *collectors) collectBasicInfo(service *ServiceYAML) {
	if service.Name != "" {
		c.names[service.Name] = true
	}
	if service.Category != "" {
		c.categories[service.Category] = true
	}
	if service.Type != "" {
		c.serviceTypes[service.Type] = true
	}
}

func (c *collectors) collectConnectionInfo(service *ServiceYAML) {
	conn := service.Connection
	serviceName := strings.ToUpper(service.Name)

	if conn.Client != "" {
		c.clients[conn.Client] = true
	}
	if conn.DefaultPort > 0 {
		c.ports[serviceName+"_PORT"] = conn.DefaultPort
	}
	if conn.DefaultUser != "" {
		c.defaultUsers[serviceName+"_USER"] = conn.DefaultUser
	}

	flags := map[string]string{
		"USER_FLAG": conn.UserFlag,
		"HOST_FLAG": conn.HostFlag,
		"PORT_FLAG": conn.PortFlag,
		"DB_FLAG":   conn.DatabaseFlag,
	}

	for suffix, value := range flags {
		if value != "" {
			c.connectionFlags[serviceName+"_"+suffix] = value
		}
	}
}

func (c *collectors) collectDockerInfo(service *ServiceYAML) {
	if image, ok := service.Docker["image"].(string); ok && image != "" {
		key := strings.ToUpper(service.Name) + "_IMAGE"
		c.images[key] = image
	}
}

func (c *collectors) collectEnvironmentVars(service *ServiceYAML) {
	serviceName := strings.ToUpper(service.Name)
	for key, value := range service.Environment {
		if strValue, ok := value.(string); ok {
			constKey := serviceName + "_" + strings.ToUpper(key)
			c.envVars[constKey] = strValue
		}
	}
}

func (c *collectors) collectHealthEndpoints(service *ServiceYAML) {
	if endpoint, ok := service.HealthCheck["endpoint"].(string); ok && endpoint != "" {
		key := strings.ToUpper(service.Name) + "_HEALTH_ENDPOINT"
		c.healthEndpoints[key] = endpoint
	}
}

func (c *collectors) collectTags(service *ServiceYAML) {
	for _, tag := range service.Tags {
		c.tags[tag] = true
	}
}

func (c *collectors) collectDependencies(service *ServiceYAML) {
	serviceName := strings.ToUpper(service.Name)
	deps := service.Dependencies

	if len(deps.Required) > 0 {
		c.dependencies[serviceName+"_DEPENDENCIES"] = deps.Required
	}
	if len(deps.Conflicts) > 0 {
		c.conflicts[serviceName+"_CONFLICTS"] = deps.Conflicts
	}
	if len(deps.Provides) > 0 {
		c.provides[serviceName+"_PROVIDES"] = deps.Provides
	}
}

func (c *collectors) toServiceConstants() *ServiceConstants {
	return &ServiceConstants{
		Categories:      c.mapToSortedConstants(c.categories, "Category"),
		ServiceTypes:    c.mapToSortedConstants(c.serviceTypes, "ServiceType"),
		Clients:         c.mapToSortedConstants(c.clients, "Client"),
		Names:           c.mapToSortedConstants(c.names, "Service"),
		Tags:            c.mapToSortedConstants(c.tags, "Tag"),
		Ports:           c.intMapToConstants(c.ports, "DefaultPort"),
		Images:          c.stringMapToConstants(c.images, "DefaultImage"),
		DefaultUsers:    c.stringMapToConstants(c.defaultUsers, "DefaultUser"),
		ConnectionFlags: c.stringMapToConstants(c.connectionFlags, "Flag"),
		EnvVars:         c.stringMapToConstants(c.envVars, "Env"),
		HealthEndpoints: c.stringMapToConstants(c.healthEndpoints, "HealthEndpoint"),
		Dependencies:    c.arrayMapToConstants(c.dependencies, "Deps"),
		Conflicts:       c.arrayMapToConstants(c.conflicts, "Conflicts"),
		Provides:        c.arrayMapToConstants(c.provides, "Provides"),
	}
}

func (c *collectors) mapToSortedConstants(m map[string]bool, prefix string) []constantData {
	var result []constantData
	for key := range m {
		result = append(result, constantData{
			Name:  prefix + toPascalCase(key),
			Value: key,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func (c *collectors) intMapToConstants(m map[string]int, prefix string) []constantData {
	var result []constantData
	for key, value := range m {
		result = append(result, constantData{
			Name:  prefix + toPascalCase(key),
			Value: fmt.Sprintf("%d", value),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func (c *collectors) stringMapToConstants(m map[string]string, prefix string) []constantData {
	var result []constantData
	for key, value := range m {
		result = append(result, constantData{
			Name:  prefix + toPascalCase(key),
			Value: value,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func (c *collectors) arrayMapToConstants(m map[string][]string, prefix string) []constantData {
	var result []constantData
	for key, value := range m {
		result = append(result, constantData{
			Name:  prefix + toPascalCase(key),
			Value: fmt.Sprintf("[]string{%s}", formatStringArray(value)),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func generateConstants(serviceConstants *ServiceConstants) error {
	tmpl, err := template.ParseFiles(TemplateFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(GeneratedFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() { _ = file.Close() }()

	return tmpl.Execute(file, serviceConstants)
}

func printSummary(sc *ServiceConstants) {
	fmt.Printf("Generated comprehensive service constants:\n")
	fmt.Printf("  Categories: %d, Types: %d, Clients: %d\n",
		len(sc.Categories), len(sc.ServiceTypes), len(sc.Clients))
	fmt.Printf("  Ports: %d, Images: %d, Env Vars: %d\n",
		len(sc.Ports), len(sc.Images), len(sc.EnvVars))
	fmt.Printf("  Dependencies: %d, Conflicts: %d, Tags: %d\n",
		len(sc.Dependencies), len(sc.Conflicts), len(sc.Tags))
}

func formatStringArray(arr []string) string {
	var quoted []string
	for _, s := range arr {
		quoted = append(quoted, fmt.Sprintf(`"%s"`, s))
	}
	return strings.Join(quoted, ", ")
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			switch strings.ToLower(part) {
			case "json":
				parts[i] = "JSON"
			case "yaml":
				parts[i] = "YAML"
			case "xml":
				parts[i] = "XML"
			case "http":
				parts[i] = "HTTP"
			case "https":
				parts[i] = "HTTPS"
			case "tty":
				parts[i] = "TTY"
			case "url":
				parts[i] = "URL"
			case "api":
				parts[i] = "API"
			default:
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}
	}
	return strings.Join(parts, "")
}
