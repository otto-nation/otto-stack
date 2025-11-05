package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

const (
	TemplateFilePath       = "cmd/generate-services/templates/services.tmpl"
	TypesTemplateFilePath  = "cmd/generate-services/templates/types.tmpl"
	GeneratedFilePath      = "internal/pkg/constants/services_generated.go"
	TypesGeneratedFilePath = "internal/pkg/constants/types_generated.go"
	ServicesDir            = "internal/config/services"
)

type constantData struct {
	Name  string
	Value string
}

type ServiceConstants struct {
	Categories      []constantData
	Clients         []constantData
	Ports           []constantData
	Names           []constantData
	Images          []constantData
	DefaultUsers    []constantData
	ConnectionFlags []constantData
	EnvVars         []constantData
	HealthEndpoints []constantData
	Tags            []constantData
	Capabilities    []constantData
	Networks        []constantData
	MemoryLimits    []constantData
	Protocols       []constantData
}

type collectors struct {
	categories      map[string]string
	clients         map[string]string
	ports           map[string]int
	names           map[string]string
	images          map[string]string
	defaultUsers    map[string]string
	connectionFlags map[string]string
	envVars         map[string]string
	healthEndpoints map[string]string
	tags            map[string]string
	capabilities    map[string]string
	networks        map[string]string
	memoryLimits    map[string]string
	protocols       map[string]string
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

	if err := generateTypes(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate types: %v\n", err)
		os.Exit(1)
	}

	printSummary(serviceConstants)
}

func extractServiceConstants() (*ServiceConstants, error) {
	c := newCollectors()

	err := filepath.Walk(ServicesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".yaml") {
			return err
		}
		return c.processFile(path)
	})

	return c.toConstants(), err
}

func newCollectors() *collectors {
	return &collectors{
		categories:      make(map[string]string),
		clients:         make(map[string]string),
		ports:           make(map[string]int),
		names:           make(map[string]string),
		images:          make(map[string]string),
		defaultUsers:    make(map[string]string),
		connectionFlags: make(map[string]string),
		envVars:         make(map[string]string),
		healthEndpoints: make(map[string]string),
		tags:            make(map[string]string),
		capabilities:    make(map[string]string),
		networks:        make(map[string]string),
		memoryLimits:    make(map[string]string),
		protocols:       make(map[string]string),
	}
}

func (c *collectors) processFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var service map[string]any
	if err := yaml.Unmarshal(data, &service); err != nil {
		return err
	}

	serviceName := strings.TrimSuffix(filepath.Base(path), ".yaml")

	c.addBasic(serviceName, path)
	c.addConnection(service, serviceName)
	c.addDocker(service, serviceName)
	c.addPorts(service, serviceName)
	c.addDependencies(service)
	c.addEnvironment(service, serviceName)
	c.addHealth(service, serviceName)
	c.addTags(service)

	return nil
}

func (c *collectors) addBasic(serviceName, path string) {
	c.names["Service"+toPascalCase(serviceName)] = serviceName

	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		cat := parts[len(parts)-2]
		c.categories["Category"+toPascalCase(cat)] = cat
	}
}

func (c *collectors) addConnection(service map[string]any, serviceName string) {
	conn, ok := service["connection"].(map[string]any)
	if !ok {
		return
	}

	if client, ok := conn["client"].(string); ok {
		c.clients["Client"+toPascalCase(client)] = client
	}

	if user, ok := conn["default_user"].(string); ok {
		c.defaultUsers["DefaultUser"+toPascalCase(serviceName)] = user
	}

	flags := []string{"user_flag", "host_flag", "port_flag", "database_flag"}
	for _, flag := range flags {
		if val, ok := conn[flag].(string); ok {
			key := toPascalCase(flag) + toPascalCase(serviceName)
			c.connectionFlags[key] = val
		}
	}
}

func (c *collectors) addDocker(service map[string]any, serviceName string) {
	docker, ok := service["docker"].(map[string]any)
	if !ok {
		return
	}

	if image, ok := docker["image"].(string); ok {
		c.images["Image"+toPascalCase(serviceName)] = image
	}

	if nets, ok := docker["networks"].([]any); ok {
		for _, net := range nets {
			if netStr, ok := net.(string); ok {
				c.networks["Network"+toPascalCase(netStr)] = netStr
			}
		}
	}

	if mem, ok := docker["memory_limit"].(string); ok {
		c.memoryLimits["MemoryLimit"+toPascalCase(mem)] = mem
	}
}

func (c *collectors) addPorts(service map[string]any, serviceName string) {
	ports, ok := service["ports"].([]any)
	if !ok {
		return
	}

	for _, port := range ports {
		portStr, ok := port.(string)
		if !ok {
			continue
		}

		parts := strings.Split(portStr, ":")
		if len(parts) == 2 {
			if portNum, err := strconv.Atoi(parts[1]); err == nil {
				c.ports["Port"+toPascalCase(serviceName)] = portNum
				c.protocols["ProtocolTcp"] = "tcp"
				return
			}
		}
	}
}

func (c *collectors) addDependencies(service map[string]any) {
	deps, ok := service["dependencies"].(map[string]any)
	if !ok {
		return
	}

	provides, ok := deps["provides"].([]any)
	if !ok {
		return
	}

	for _, cap := range provides {
		if capStr, ok := cap.(string); ok {
			c.capabilities["Capability"+toPascalCase(capStr)] = capStr
		}
	}
}

func (c *collectors) addEnvironment(service map[string]any, serviceName string) {
	env, ok := service["environment"].(map[string]any)
	if !ok {
		return
	}

	for key, value := range env {
		if strValue, ok := value.(string); ok {
			envKey := "Env" + toPascalCase(serviceName) + toPascalCase(key)
			c.envVars[envKey] = strValue
		}
	}
}

func (c *collectors) addHealth(service map[string]any, serviceName string) {
	health, ok := service["health_check"].(map[string]any)
	if !ok {
		return
	}

	if endpoint, ok := health["endpoint"].(string); ok {
		c.healthEndpoints["HealthEndpoint"+toPascalCase(serviceName)] = endpoint
	}
}

func (c *collectors) addTags(service map[string]any) {
	tags, ok := service["tags"].([]any)
	if !ok {
		return
	}

	for _, tag := range tags {
		if tagStr, ok := tag.(string); ok {
			c.tags["Tag"+toPascalCase(tagStr)] = tagStr
		}
	}
}

func (c *collectors) toConstants() *ServiceConstants {
	return &ServiceConstants{
		Categories:      stringMapToConstants(c.categories),
		Clients:         stringMapToConstants(c.clients),
		Ports:           intMapToConstants(c.ports),
		Names:           stringMapToConstants(c.names),
		Images:          stringMapToConstants(c.images),
		DefaultUsers:    stringMapToConstants(c.defaultUsers),
		ConnectionFlags: stringMapToConstants(c.connectionFlags),
		EnvVars:         stringMapToConstants(c.envVars),
		HealthEndpoints: stringMapToConstants(c.healthEndpoints),
		Tags:            stringMapToConstants(c.tags),
		Capabilities:    stringMapToConstants(c.capabilities),
		Networks:        stringMapToConstants(c.networks),
		MemoryLimits:    stringMapToConstants(c.memoryLimits),
		Protocols:       stringMapToConstants(c.protocols),
	}
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
	defer file.Close()

	return tmpl.Execute(file, serviceConstants)
}

func generateTypes() error {
	tmpl, err := template.ParseFiles(TypesTemplateFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse types template: %w", err)
	}

	file, err := os.Create(TypesGeneratedFilePath)
	if err != nil {
		return fmt.Errorf("failed to create types file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, nil)
}

func printSummary(sc *ServiceConstants) {
	fmt.Printf("Generated service constants:\n")
	fmt.Printf("  Categories: %d, Images: %d, Ports: %d\n",
		len(sc.Categories), len(sc.Images), len(sc.Ports))
}

func stringMapToConstants(m map[string]string) []constantData {
	var result []constantData
	for key, value := range m {
		result = append(result, constantData{Name: key, Value: value})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func intMapToConstants(m map[string]int) []constantData {
	var result []constantData
	for key, value := range m {
		result = append(result, constantData{Name: key, Value: fmt.Sprintf("%d", value)})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}
