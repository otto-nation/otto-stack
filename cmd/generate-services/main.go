package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-playground/validator/v10"
	"github.com/otto-nation/otto-stack/cmd/codegen"
	pkgerrors "github.com/otto-nation/otto-stack/internal/pkg/errors"
	"gopkg.in/yaml.v3"
)

var validate = validator.New()
var serviceSchema *ServiceSchema // Loaded at startup

// Schema-derived constants for validation
const (
	ServiceTypeContainer     = "container"
	ServiceTypeComposite     = "composite"
	ServiceTypeConfiguration = "configuration"
)

const (
	TemplateFilePath            = "cmd/generate-services/templates/services.tmpl"
	MainConfigTemplateFilePath  = "cmd/generate-services/templates/main_config.tmpl"
	SchemaFilePath              = "cmd/generate-services/service-schema.json"
	GeneratedFilePath           = "internal/pkg/services/services_generated.go"
	GeneratedConfigsDir         = "internal/pkg/types"
	MainConfigGeneratedFilePath = GeneratedConfigsDir + "/service_config_generated.go"
	ServicesDir                 = "internal/config/services"
	YAMLExtension               = ".yaml"
	DefaultProtocol             = "tcp"
)

// Keys defines all YAML structure keys
var Keys = struct {
	Connection struct {
		Root         string
		Client       string
		DefaultUser  string
		DefaultPort  string
		UserFlag     string
		HostFlag     string
		PortFlag     string
		DatabaseFlag string
	}
	Container struct {
		Root        string
		Image       string
		Ports       string
		Networks    string
		MemoryLimit string
		External    string
		Protocol    string
	}
	Service struct {
		Root            string
		Characteristics string
	}
	Environment  string
	Dependencies struct {
		Root     string
		Provides string
	}
	HealthCheck struct {
		Root     string
		Endpoint string
	}
	Tags         string
	ConfigSchema string
}{
	Connection: struct {
		Root         string
		Client       string
		DefaultUser  string
		DefaultPort  string
		UserFlag     string
		HostFlag     string
		PortFlag     string
		DatabaseFlag string
	}{"connection", "client", "default_user", "default_port", "user_flag", "host_flag", "port_flag", "database_flag"},
	Container: struct {
		Root        string
		Image       string
		Ports       string
		Networks    string
		MemoryLimit string
		External    string
		Protocol    string
	}{"container", "image", "ports", "networks", "memory_limit", "external", "protocol"},
	Service: struct {
		Root            string
		Characteristics string
	}{"service", "characteristics"},
	Environment: "environment",
	Dependencies: struct {
		Root     string
		Provides string
	}{"dependencies", "provides"},
	HealthCheck: struct {
		Root     string
		Endpoint string
	}{"health_check", "endpoint"},
	Tags:         "tags",
	ConfigSchema: "configuration_schema",
}

// Prefix defines constant name prefixes
var Prefix = struct {
	Service, Category, Client, DefaultUser, Image, MemoryLimit            string
	Network, Port, Protocol, Env, EnvKey, Capability, HealthEndpoint, Tag string
}{
	"Service", "Category", "Client", "DefaultUser", "Image", "MemoryLimit",
	"Network", "Port", "Protocol", "Env", "EnvKey", "Capability", "HealthEndpoint", "Tag",
}

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
	EnvKeys         []constantData
	HealthEndpoints []constantData
	Tags            []constantData
	Capabilities    []constantData
	Networks        []constantData
	MemoryLimits    []constantData
	Protocols       []constantData
	Services        []ServiceData
	ConfigSchemas   []ServiceConfigSchema
}

type ServiceData struct {
	Name            string
	Characteristics []string
}

type ServiceConfigSchema struct {
	ServiceName string
	Schema      map[string]any
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
	envKeys         map[string]string
	healthEndpoints map[string]string
	tags            map[string]string
	capabilities    map[string]string
	networks        map[string]string
	memoryLimits    map[string]string
	protocols       map[string]string
	services        []ServiceData
	configSchemas   []ServiceConfigSchema
}

func main() {
	// Load schema first
	var err error
	serviceSchema, err = loadServiceSchema(SchemaFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load schema: %v\n", err)
		os.Exit(1)
	}

	serviceConstants, configSchemas, err := extractServiceConstants()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to extract service constants: %v\n", err)
		os.Exit(1)
	}

	if err := generateConstants(serviceConstants); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate constants: %v\n", err)
		os.Exit(1)
	}

	if err := generateMultiFileSchema(configSchemas); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate multi-file schema: %v\n", err)
		os.Exit(1)
	}

	if err := generateDockerTypes(serviceSchema); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate docker types: %v\n", err)
		os.Exit(1)
	}

	if err := generateYAMLKeys(serviceSchema); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate YAML keys: %v\n", err)
		os.Exit(1)
	}

	printSummary(serviceConstants)
}

func extractServiceConstants() (*ServiceConstants, []ServiceConfigSchema, error) {
	collectors := initCollectors()
	var configSchemas []ServiceConfigSchema

	err := filepath.Walk(ServicesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, YAMLExtension) {
			return err
		}
		return processServiceFile(path, collectors, &configSchemas)
	})

	if err != nil {
		return nil, nil, err
	}

	return buildServiceConstants(collectors, configSchemas), configSchemas, nil
}

func initCollectors() *collectors {
	return &collectors{
		categories:      make(map[string]string),
		clients:         make(map[string]string),
		ports:           make(map[string]int),
		names:           make(map[string]string),
		images:          make(map[string]string),
		defaultUsers:    make(map[string]string),
		connectionFlags: make(map[string]string),
		envVars:         make(map[string]string),
		envKeys:         make(map[string]string),
		healthEndpoints: make(map[string]string),
		tags:            make(map[string]string),
		capabilities:    make(map[string]string),
		networks:        make(map[string]string),
		memoryLimits:    make(map[string]string),
		protocols:       make(map[string]string),
	}
}

func processServiceFile(path string, collectors *collectors, configSchemas *[]ServiceConfigSchema) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var service map[string]any
	if err := yaml.Unmarshal(data, &service); err != nil {
		return err
	}

	// Validate required fields
	if err := validateServiceYAML(service, path); err != nil {
		return fmt.Errorf("validation failed for %s: %w", path, err)
	}

	serviceName := strings.TrimSuffix(filepath.Base(path), YAMLExtension)
	processService(service, serviceName, path, collectors, configSchemas)
	return nil
}

func validateServiceYAML(service map[string]any, path string) error {
	// Unmarshal into validation struct
	data, err := yaml.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service for validation: %w", err)
	}

	var svc serviceValidation
	if err := yaml.Unmarshal(data, &svc); err != nil {
		return fmt.Errorf("failed to unmarshal service for validation: %w", err)
	}

	// Use validator
	validate := validator.New()
	if err := validate.Struct(svc); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Validate enums against schema
	if svc.ServiceType != "" {
		if err := validateEnumField("service_type", svc.ServiceType); err != nil {
			return err
		}
	}
	if svc.InitService != nil && svc.InitService.Mode != "" {
		if err := validateEnumField("init_service_mode", svc.InitService.Mode); err != nil {
			return err
		}
	}

	// Custom validation: container type requires container.image
	if svc.ServiceType == ServiceTypeContainer {
		if svc.Container == nil || svc.Container.Image == "" {
			return fmt.Errorf("service_type '%s' requires container.image", ServiceTypeContainer)
		}
	}

	// Custom validation: init_service.enabled requires mode
	if svc.InitService != nil && svc.InitService.Enabled {
		if svc.InitService.Mode == "" {
			return fmt.Errorf("init_service.mode is required when init_service.enabled is true")
		}
	}

	return nil
}

// serviceValidation defines the validation schema for service YAMLs
type serviceValidation struct {
	Name        string                 `yaml:"name" validate:"required"`
	Description string                 `yaml:"description" validate:"required"`
	ServiceType string                 `yaml:"service_type"` // validated dynamically
	Container   *containerValidation   `yaml:"container,omitempty"`
	InitService *initServiceValidation `yaml:"init_service,omitempty"`
}

type containerValidation struct {
	Image string `yaml:"image"`
}

type initServiceValidation struct {
	Enabled bool   `yaml:"enabled"`
	Mode    string `yaml:"mode"` // validated dynamically
}

// validateEnumField validates a field against schema enum
func validateEnumField(fieldName, value string) error {
	enumName := strings.ReplaceAll(fieldName, ".", "_")
	if enumDef, exists := serviceSchema.Enums[enumName]; exists {
		for _, validValue := range enumDef.Values {
			if value == validValue {
				return nil
			}
		}
		return fmt.Errorf("%s must be one of: %s", fieldName, strings.Join(enumDef.Values, ", "))
	}
	return nil
}

func buildServiceConstants(collectors *collectors, configSchemas []ServiceConfigSchema) *ServiceConstants {
	return &ServiceConstants{
		Categories:      stringMapToConstants(collectors.categories),
		Clients:         stringMapToConstants(collectors.clients),
		Ports:           intMapToConstants(collectors.ports),
		Names:           stringMapToConstants(collectors.names),
		Images:          stringMapToConstants(collectors.images),
		DefaultUsers:    stringMapToConstants(collectors.defaultUsers),
		ConnectionFlags: stringMapToConstants(collectors.connectionFlags),
		EnvVars:         stringMapToConstants(collectors.envVars),
		EnvKeys:         stringMapToConstants(collectors.envKeys),
		HealthEndpoints: stringMapToConstants(collectors.healthEndpoints),
		Tags:            stringMapToConstants(collectors.tags),
		Capabilities:    stringMapToConstants(collectors.capabilities),
		Networks:        stringMapToConstants(collectors.networks),
		MemoryLimits:    stringMapToConstants(collectors.memoryLimits),
		Protocols:       stringMapToConstants(collectors.protocols),
		Services:        collectors.services,
		ConfigSchemas:   configSchemas,
	}
}

func processService(service map[string]any, serviceName, path string, collectors *collectors, configSchemas *[]ServiceConfigSchema) {
	processBasicInfo(serviceName, path, collectors)
	processConnection(service, serviceName, collectors)
	processContainer(service, serviceName, collectors)
	processEnvironment(service, serviceName, collectors)
	processDependencies(service, collectors)
	processHealth(service, serviceName, collectors)
	processTags(service, collectors)
	processCharacteristics(service, serviceName, collectors)
	processConfigSchema(service, serviceName, configSchemas)
}

func processBasicInfo(serviceName, path string, collectors *collectors) {
	collectors.names[Prefix.Service+toPascalCase(serviceName)] = serviceName

	// Extract category from directory name (parent of the service file)
	dir := filepath.Dir(path)
	cat := filepath.Base(dir)
	if cat != "." && cat != "" {
		collectors.categories[Prefix.Category+toPascalCase(cat)] = cat
	}
}

func processConnection(service map[string]any, serviceName string, collectors *collectors) {
	serviceSection, ok := service[Keys.Service.Root].(map[string]any)
	if !ok {
		return
	}
	conn, ok := serviceSection[Keys.Connection.Root].(map[string]any)
	if !ok {
		return
	}

	if client, ok := conn[Keys.Connection.Client].(string); ok {
		collectors.clients[Prefix.Client+toPascalCase(client)] = client
	}
	if user, ok := conn[Keys.Connection.DefaultUser].(string); ok {
		collectors.defaultUsers[Prefix.DefaultUser+toPascalCase(serviceName)] = user
	}
	if port, ok := conn[Keys.Connection.DefaultPort].(int); ok {
		collectors.ports[Prefix.Port+toPascalCase(serviceName)] = port
	} else if portStr, ok := conn[Keys.Connection.DefaultPort].(string); ok {
		if portNum, err := strconv.Atoi(portStr); err == nil {
			collectors.ports[Prefix.Port+toPascalCase(serviceName)] = portNum
		}
	}

	processConnectionFlags(conn, serviceName, collectors)
}

func processConnectionFlags(conn map[string]any, serviceName string, collectors *collectors) {
	flags := []string{Keys.Connection.UserFlag, Keys.Connection.HostFlag, Keys.Connection.PortFlag, Keys.Connection.DatabaseFlag}
	for _, flag := range flags {
		if val, ok := conn[flag].(string); ok {
			key := toPascalCase(flag) + toPascalCase(serviceName)
			collectors.connectionFlags[key] = val
		}
	}
}

func processContainer(service map[string]any, serviceName string, collectors *collectors) {
	container, ok := service[Keys.Container.Root].(map[string]any)
	if !ok {
		return
	}

	processContainerBasics(container, serviceName, collectors)
	processContainerPorts(container, serviceName, collectors)
	processContainerNetworks(container, collectors)
}

func processContainerBasics(container map[string]any, serviceName string, collectors *collectors) {
	if image, ok := container[Keys.Container.Image].(string); ok {
		collectors.images[Prefix.Image+toPascalCase(serviceName)] = image
	}
	if mem, ok := container[Keys.Container.MemoryLimit].(string); ok {
		collectors.memoryLimits[Prefix.MemoryLimit+toPascalCase(mem)] = mem
	}
}

func processContainerPorts(container map[string]any, serviceName string, collectors *collectors) {
	ports, ok := container[Keys.Container.Ports].([]any)
	if !ok || len(ports) == 0 {
		return
	}

	portMap, ok := ports[0].(map[string]any)
	if !ok {
		return
	}

	if external, ok := portMap[Keys.Container.External].(string); ok {
		if portNum, err := strconv.Atoi(external); err == nil {
			collectors.ports[Prefix.Port+toPascalCase(serviceName)] = portNum
		}
	}

	protocol := DefaultProtocol
	if p, ok := portMap[Keys.Container.Protocol].(string); ok {
		protocol = p
	}
	collectors.protocols[Prefix.Protocol+toPascalCase(protocol)] = protocol
}

func processContainerNetworks(container map[string]any, collectors *collectors) {
	nets, ok := container[Keys.Container.Networks].([]any)
	if !ok {
		return
	}
	for _, net := range nets {
		if netStr, ok := net.(string); ok {
			collectors.networks["Network"+toPascalCase(netStr)] = netStr
		}
	}
}

func processEnvironment(service map[string]any, serviceName string, collectors *collectors) {
	env, ok := service[Keys.Environment].(map[string]any)
	if !ok {
		return
	}

	for key, value := range env {
		if strValue, ok := value.(string); ok {
			envKey := Prefix.Env + toPascalCase(serviceName) + toPascalCase(key)
			collectors.envVars[envKey] = strValue

			envKeyConstant := Prefix.EnvKey + toPascalCase(key)
			if _, exists := collectors.envKeys[envKeyConstant]; !exists {
				collectors.envKeys[envKeyConstant] = key
			}
		}
	}
}

func processDependencies(service map[string]any, collectors *collectors) {
	deps, ok := service[Keys.Dependencies.Root].(map[string]any)
	if !ok {
		return
	}

	provides, ok := deps[Keys.Dependencies.Provides].([]any)
	if !ok {
		return
	}

	for _, cap := range provides {
		if capStr, ok := cap.(string); ok {
			collectors.capabilities[Prefix.Capability+toPascalCase(capStr)] = capStr
		}
	}
}

func processHealth(service map[string]any, serviceName string, collectors *collectors) {
	health, ok := service[Keys.HealthCheck.Root].(map[string]any)
	if !ok {
		return
	}

	if endpoint, ok := health[Keys.HealthCheck.Endpoint].(string); ok {
		collectors.healthEndpoints[Prefix.HealthEndpoint+toPascalCase(serviceName)] = endpoint
	}
}

func processTags(service map[string]any, collectors *collectors) {
	tags, ok := service[Keys.Tags].([]any)
	if !ok {
		return
	}

	for _, tag := range tags {
		if tagStr, ok := tag.(string); ok {
			collectors.tags[Prefix.Tag+toPascalCase(tagStr)] = tagStr
		}
	}
}

func processCharacteristics(service map[string]any, serviceName string, collectors *collectors) {
	serviceSection, ok := service[Keys.Service.Root].(map[string]any)
	if !ok {
		return
	}

	characteristics, ok := serviceSection[Keys.Service.Characteristics].([]any)
	if !ok {
		return
	}

	var charStrings []string
	for _, char := range characteristics {
		if charStr, ok := char.(string); ok {
			charStrings = append(charStrings, charStr)
		}
	}

	if len(charStrings) > 0 {
		collectors.services = append(collectors.services, ServiceData{
			Name:            serviceName,
			Characteristics: charStrings,
		})
	}
}

func processConfigSchema(service map[string]any, serviceName string, configSchemas *[]ServiceConfigSchema) {
	schema, ok := service["configuration_schema"].(map[string]any)
	if !ok {
		return
	}

	*configSchemas = append(*configSchemas, ServiceConfigSchema{
		ServiceName: serviceName,
		Schema:      schema,
	})
}

const (
	minPathParts = 3
	pathParts    = 2
)

func generateConstants(serviceConstants *ServiceConstants) error {
	tmpl, err := template.ParseFiles(TemplateFilePath)
	if err != nil {
		return pkgerrors.NewServiceError("generator", "parse template", err)
	}

	file, err := os.Create(GeneratedFilePath)
	if err != nil {
		// Try creating the directory and retry
		if err := codegen.EnsureDir(filepath.Dir(GeneratedFilePath)); err != nil {
			return pkgerrors.NewServiceError("generator", "create directory", err)
		}
		file, err = os.Create(GeneratedFilePath)
		if err != nil {
			return pkgerrors.NewServiceError("generator", "create file", err)
		}
	}
	defer func() { _ = file.Close() }()

	return tmpl.Execute(file, serviceConstants)
}

func generateMultiFileSchema(configSchemas []ServiceConfigSchema) error {
	// Generate service config files
	var serviceFiles []codegen.ServiceFileData

	for _, schema := range configSchemas {
		if len(schema.Schema) == 0 {
			continue
		}

		structName := codegen.ToPascalCase(schema.ServiceName) + "Config"
		fileName := strings.ToLower(codegen.ToPascalCase(schema.ServiceName)) + "_config_generated.go"

		serviceFiles = append(serviceFiles, codegen.ServiceFileData{
			ServiceName: schema.ServiceName,
			StructName:  structName,
			FileName:    fileName,
			Schema:      schema.Schema,
		})

		// Generate individual config file
		filePath := filepath.Join(GeneratedConfigsDir, fileName)
		if err := generateIndividualConfigFile(structName, filePath); err != nil {
			return fmt.Errorf("failed to generate %s: %w", fileName, err)
		}
	}

	// Generate main ServiceConfig with embedded structs
	mainConfigData := struct {
		Services []codegen.ServiceFileData
	}{
		Services: serviceFiles,
	}

	executor := codegen.NewTemplateExecutor(MainConfigTemplateFilePath, MainConfigGeneratedFilePath)
	funcMap := template.FuncMap{
		"toPascalCase": codegen.ToPascalCase,
	}

	return executor.ExecuteTemplateWithFuncs(mainConfigData, funcMap)
}

func generateIndividualConfigFile(structName, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write simple empty struct - schema details not needed for now
	fmt.Fprintln(file, "// Code generated by generate-services. DO NOT EDIT.")
	fmt.Fprintln(file, "package types")
	fmt.Fprintln(file)
	fmt.Fprintf(file, "type %s struct {}\n", structName)

	return nil
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

func generateDockerTypes(schema *ServiceSchema) error {
	const (
		DockerTemplateFilePath  = "cmd/generate-services/templates/docker.tmpl"
		DockerGeneratedFilePath = "internal/core/docker/types_generated.go"
	)

	tmpl, err := template.ParseFiles(DockerTemplateFilePath)
	if err != nil {
		return pkgerrors.NewServiceError("generator", "parse docker template", err)
	}

	file, err := os.Create(DockerGeneratedFilePath)
	if err != nil {
		if err := codegen.EnsureDir(filepath.Dir(DockerGeneratedFilePath)); err != nil {
			return pkgerrors.NewServiceError("generator", "create directory", err)
		}
		file, err = os.Create(DockerGeneratedFilePath)
		if err != nil {
			return pkgerrors.NewServiceError("generator", "create docker types file", err)
		}
	}
	defer func() { _ = file.Close() }()

	// Generate enum constants from schema
	enumGroups := schema.generateEnumConstants()

	return tmpl.Execute(file, map[string]any{
		"EnumGroups": enumGroups,
	})
}

func generateYAMLKeys(schema *ServiceSchema) error {
	const (
		KeysTemplateFilePath  = "cmd/generate-services/templates/keys.tmpl"
		KeysGeneratedFilePath = "internal/pkg/services/keys_generated.go"
	)

	tmpl, err := template.ParseFiles(KeysTemplateFilePath)
	if err != nil {
		return pkgerrors.NewServiceError("generator", "parse keys template", err)
	}

	file, err := os.Create(KeysGeneratedFilePath)
	if err != nil {
		return pkgerrors.NewServiceError("generator", "create keys file", err)
	}
	defer func() { _ = file.Close() }()

	// Generate YAML key constants from schema
	yamlKeys := schema.generateYAMLKeys()

	return tmpl.Execute(file, yamlKeys)
}
