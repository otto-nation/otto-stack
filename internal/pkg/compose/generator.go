package compose

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/otto-nation/otto-stack/internal/pkg/services"
)

// ComposeService represents a Docker Compose service configuration
type ComposeService struct {
	Image       string            `yaml:"image"`
	Ports       []string          `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Command     []string          `yaml:"command,omitempty"`
	HealthCheck *HealthCheck      `yaml:"healthcheck,omitempty"`
	Restart     string            `yaml:"restart,omitempty"`
}

// HealthCheck represents Docker health check configuration
type HealthCheck struct {
	Test     []string `yaml:"test"`
	Interval string   `yaml:"interval,omitempty"`
	Timeout  string   `yaml:"timeout,omitempty"`
	Retries  int      `yaml:"retries,omitempty"`
}

// ComposeFile represents a complete Docker Compose file
type ComposeFile struct {
	Version  string                    `yaml:"version"`
	Services map[string]ComposeService `yaml:"services"`
	Volumes  map[string]interface{}    `yaml:"volumes,omitempty"`
	Networks map[string]interface{}    `yaml:"networks,omitempty"`
}

// Generator handles docker-compose file generation
type Generator struct {
	projectName string
	registry    *services.ServiceRegistry
}

// NewGenerator creates a new compose generator
func NewGenerator(projectName string, servicesPath string) (*Generator, error) {
	registry, err := services.NewServiceRegistry(servicesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create service registry: %w", err)
	}

	return &Generator{
		projectName: projectName,
		registry:    registry,
	}, nil
}

// Generate creates a docker-compose file for the specified services
func (g *Generator) Generate(serviceNames []string) (*ComposeFile, error) {
	compose := &ComposeFile{
		Version:  "3.8",
		Services: make(map[string]ComposeService),
		Volumes:  make(map[string]interface{}),
		Networks: make(map[string]interface{}),
	}

	// Add default network
	compose.Networks["default"] = map[string]interface{}{
		"name": fmt.Sprintf("%s-network", g.projectName),
	}

	for _, serviceName := range serviceNames {
		if err := g.addService(compose, serviceName); err != nil {
			return nil, fmt.Errorf("failed to add service %s: %w", serviceName, err)
		}
	}

	return compose, nil
}

// GenerateYAML generates the docker-compose YAML content
func (g *Generator) GenerateYAML(serviceNames []string) ([]byte, error) {
	compose, err := g.Generate(serviceNames)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(compose)
}

// addService adds a specific service to the compose file based on service definitions
func (g *Generator) addService(compose *ComposeFile, serviceName string) error {
	// Try to get service from registry first
	if g.registry != nil {
		if serviceDef, exists := g.registry.GetService(serviceName); exists {
			return g.addServiceFromDefinition(compose, serviceName, serviceDef)
		}
	}

	// Fallback to hardcoded services for backward compatibility
	switch serviceName {
	case "postgres":
		g.addPostgres(compose)
	case "redis":
		g.addRedis(compose)
	case "mysql":
		g.addMySQL(compose)
	case "kafka-broker":
		g.addKafka(compose)
	case "zookeeper":
		g.addZookeeper(compose)
	case "prometheus":
		g.addPrometheus(compose)
	case "jaeger":
		g.addJaeger(compose)
	case "localstack-core":
		g.addLocalstack(compose)
	default:
		return fmt.Errorf("unknown service: %s", serviceName)
	}
	return nil
}

// addServiceFromDefinition creates a compose service from a service definition
func (g *Generator) addServiceFromDefinition(compose *ComposeFile, serviceName string, def services.ServiceDefinition) error {
	service := ComposeService{
		Restart: "unless-stopped",
	}

	// Set image from defaults
	if def.Defaults.Image != "" {
		service.Image = def.Defaults.Image
	}

	// Configure environment variables
	if len(def.Environment) > 0 {
		service.Environment = make(map[string]string)
		for key, value := range def.Environment {
			service.Environment[key] = value
		}
	}

	// Add Docker-specific environment variables
	if len(def.Docker.Environment) > 0 {
		if service.Environment == nil {
			service.Environment = make(map[string]string)
		}
		for _, env := range def.Docker.Environment {
			// Parse KEY=VALUE format
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				service.Environment[parts[0]] = parts[1]
			}
		}
	}

	// Configure ports
	if def.Defaults.Port > 0 {
		service.Ports = []string{fmt.Sprintf("%d:%d", def.Defaults.Port, def.Defaults.Port)}
	}
	if len(def.Docker.Ports) > 0 {
		service.Ports = def.Docker.Ports
	}

	// Configure volumes
	if len(def.Docker.Volumes) > 0 {
		service.Volumes = def.Docker.Volumes
	}

	// Add named volumes if defined
	for _, vol := range def.Volumes {
		volumeName := fmt.Sprintf("%s-%s", g.projectName, vol.Name)
		compose.Volumes[volumeName] = nil

		// Add volume mount to service
		volumeMount := fmt.Sprintf("%s:%s", volumeName, vol.Mount)
		service.Volumes = append(service.Volumes, volumeMount)
	}

	// Configure health check
	if len(def.Docker.HealthCheck.Test) > 0 {
		service.HealthCheck = &HealthCheck{
			Test:     def.Docker.HealthCheck.Test,
			Interval: def.Docker.HealthCheck.Interval,
			Timeout:  def.Docker.HealthCheck.Timeout,
			Retries:  def.Docker.HealthCheck.Retries,
		}
	}

	// Configure restart policy
	if def.Docker.Restart != "" {
		service.Restart = def.Docker.Restart
	}

	// Configure dependencies
	if len(def.Dependencies.Required) > 0 {
		service.DependsOn = def.Dependencies.Required
	}
	if len(def.Docker.DependsOn) > 0 {
		service.DependsOn = append(service.DependsOn, def.Docker.DependsOn...)
	}

	compose.Services[serviceName] = service
	return nil
}

// Service-specific generators
func (g *Generator) addPostgres(compose *ComposeFile) {
	volumeName := fmt.Sprintf("%s-postgres-data", g.projectName)
	compose.Volumes[volumeName] = nil

	compose.Services["postgres"] = ComposeService{
		Image: "postgres:15-alpine",
		Environment: map[string]string{
			"POSTGRES_DB":       g.projectName,
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
		},
		Ports:   []string{"5432:5432"},
		Volumes: []string{fmt.Sprintf("%s:/var/lib/postgresql/data", volumeName)},
		HealthCheck: &HealthCheck{
			Test:     []string{"CMD-SHELL", "pg_isready -U postgres"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  5,
		},
		Restart: "unless-stopped",
	}
}

func (g *Generator) addRedis(compose *ComposeFile) {
	volumeName := fmt.Sprintf("%s-redis-data", g.projectName)
	compose.Volumes[volumeName] = nil

	compose.Services["redis"] = ComposeService{
		Image:   "redis:7-alpine",
		Ports:   []string{"6379:6379"},
		Volumes: []string{fmt.Sprintf("%s:/data", volumeName)},
		HealthCheck: &HealthCheck{
			Test:     []string{"CMD", "redis-cli", "ping"},
			Interval: "10s",
			Timeout:  "3s",
			Retries:  3,
		},
		Restart: "unless-stopped",
	}
}

func (g *Generator) addMySQL(compose *ComposeFile) {
	volumeName := fmt.Sprintf("%s-mysql-data", g.projectName)
	compose.Volumes[volumeName] = nil

	compose.Services["mysql"] = ComposeService{
		Image: "mysql:8.0",
		Environment: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      g.projectName,
			"MYSQL_USER":          "mysql",
			"MYSQL_PASSWORD":      "mysql",
		},
		Ports:   []string{"3306:3306"},
		Volumes: []string{fmt.Sprintf("%s:/var/lib/mysql", volumeName)},
		HealthCheck: &HealthCheck{
			Test:     []string{"CMD", "mysqladmin", "ping", "-h", "localhost"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  5,
		},
		Restart: "unless-stopped",
	}
}

func (g *Generator) addKafka(compose *ComposeFile) {
	kafkaVolume := fmt.Sprintf("%s-kafka-data", g.projectName)
	compose.Volumes[kafkaVolume] = nil

	// Add Zookeeper if not already present
	if _, exists := compose.Services["zookeeper"]; !exists {
		g.addZookeeper(compose)
	}

	compose.Services["kafka"] = ComposeService{
		Image: "confluentinc/cp-kafka:latest",
		Environment: map[string]string{
			"KAFKA_BROKER_ID":                        "1",
			"KAFKA_ZOOKEEPER_CONNECT":                "zookeeper:2181",
			"KAFKA_ADVERTISED_LISTENERS":             "PLAINTEXT://localhost:9092",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
		},
		Ports:     []string{"9092:9092"},
		DependsOn: []string{"zookeeper"},
		Volumes:   []string{fmt.Sprintf("%s:/var/lib/kafka/data", kafkaVolume)},
		Restart:   "unless-stopped",
	}

	// Add Kafka init container for topic creation
	g.addKafkaInit(compose)
}

func (g *Generator) addZookeeper(compose *ComposeFile) {
	zookeeperVolume := fmt.Sprintf("%s-zookeeper-data", g.projectName)
	compose.Volumes[zookeeperVolume] = nil

	compose.Services["zookeeper"] = ComposeService{
		Image: "confluentinc/cp-zookeeper:latest",
		Environment: map[string]string{
			"ZOOKEEPER_CLIENT_PORT": "2181",
			"ZOOKEEPER_TICK_TIME":   "2000",
		},
		Volumes: []string{fmt.Sprintf("%s:/var/lib/zookeeper/data", zookeeperVolume)},
		Restart: "unless-stopped",
	}
}

func (g *Generator) addKafkaInit(compose *ComposeFile) {
	initScript := `#!/bin/bash
set -e

# Wait for Kafka to be ready
echo "Waiting for Kafka to be ready..."
while ! kafka-topics --bootstrap-server kafka:9092 --list > /dev/null 2>&1; do
    echo "Kafka not ready yet, waiting..."
    sleep 2
done

echo "Creating default topics..."

# Create default topics
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic events --partitions 3 --replication-factor 1
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic notifications --partitions 1 --replication-factor 1
kafka-topics --bootstrap-server kafka:9092 --create --if-not-exists --topic logs --partitions 1 --replication-factor 1

echo "Kafka topics created successfully"
`

	compose.Services["kafka-init"] = ComposeService{
		Image:     "confluentinc/cp-kafka:latest",
		DependsOn: []string{"kafka"},
		Command: []string{
			"bash", "-c", initScript,
		},
		Restart: "no",
	}
}

func (g *Generator) addLocalstack(compose *ComposeFile) {
	volumeName := fmt.Sprintf("%s-localstack-data", g.projectName)
	compose.Volumes[volumeName] = nil

	compose.Services["localstack"] = ComposeService{
		Image: "localstack/localstack:latest",
		Environment: map[string]string{
			"SERVICES":    "s3,dynamodb,sqs,sns,lambda",
			"DEBUG":       "1",
			"DATA_DIR":    "/tmp/localstack/data",
			"DOCKER_HOST": "unix:///var/run/docker.sock",
		},
		Ports: []string{
			"4566:4566",           // LocalStack main port
			"4510-4559:4510-4559", // External service ports
		},
		Volumes: []string{
			fmt.Sprintf("%s:/tmp/localstack", volumeName),
			"/var/run/docker.sock:/var/run/docker.sock",
		},
		Restart: "unless-stopped",
	}

	// Add LocalStack init container for AWS resource creation
	g.addLocalstackInit(compose)
}

func (g *Generator) addLocalstackInit(compose *ComposeFile) {
	initScript := `#!/bin/bash
set -e

# Wait for LocalStack to be ready
echo "Waiting for LocalStack to be ready..."
while ! curl -s http://localstack:4566/_localstack/health > /dev/null; do
    echo "LocalStack not ready yet, waiting..."
    sleep 2
done

echo "Creating AWS resources..."

# Create default SQS queues
aws --endpoint-url=http://localstack:4566 sqs create-queue --queue-name events-queue --region us-east-1
aws --endpoint-url=http://localstack:4566 sqs create-queue --queue-name notifications-queue --region us-east-1

# Create default SNS topics
aws --endpoint-url=http://localstack:4566 sns create-topic --name events-topic --region us-east-1
aws --endpoint-url=http://localstack:4566 sns create-topic --name notifications-topic --region us-east-1

# Create default DynamoDB table
aws --endpoint-url=http://localstack:4566 dynamodb create-table \
    --table-name app-data \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1

# Create default S3 bucket
aws --endpoint-url=http://localstack:4566 s3 mb s3://app-bucket --region us-east-1

echo "AWS resources created successfully"
`

	compose.Services["localstack-init"] = ComposeService{
		Image:     "amazon/aws-cli:latest",
		DependsOn: []string{"localstack"},
		Command: []string{
			"bash", "-c", initScript,
		},
		Environment: map[string]string{
			"AWS_ACCESS_KEY_ID":     "test",
			"AWS_SECRET_ACCESS_KEY": "test",
			"AWS_DEFAULT_REGION":    "us-east-1",
		},
		Restart: "no",
	}
}

func (g *Generator) addPrometheus(compose *ComposeFile) {
	volumeName := fmt.Sprintf("%s-prometheus-data", g.projectName)
	compose.Volumes[volumeName] = nil

	compose.Services["prometheus"] = ComposeService{
		Image:   "prom/prometheus:latest",
		Ports:   []string{"9090:9090"},
		Volumes: []string{fmt.Sprintf("%s:/prometheus", volumeName)},
		Command: []string{
			"--config.file=/etc/prometheus/prometheus.yml",
			"--storage.tsdb.path=/prometheus",
			"--web.console.libraries=/etc/prometheus/console_libraries",
			"--web.console.templates=/etc/prometheus/consoles",
		},
		Restart: "unless-stopped",
	}
}

func (g *Generator) addJaeger(compose *ComposeFile) {
	compose.Services["jaeger"] = ComposeService{
		Image: "jaegertracing/all-in-one:latest",
		Environment: map[string]string{
			"COLLECTOR_OTLP_ENABLED": "true",
		},
		Ports: []string{
			"16686:16686", // Jaeger UI
			"14268:14268", // Jaeger collector
			"4317:4317",   // OTLP gRPC
			"4318:4318",   // OTLP HTTP
		},
		Restart: "unless-stopped",
	}
}
