package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// Timeout constants
	HTTPTimeout    = 5 * time.Second
	RetryDelay     = 2 * time.Second
	MaxRetries     = 30
	ServiceTimeout = 60 * time.Second

	// HTTP status codes
	HTTPStatusOK = 200

	// Environment variable names
	EnvServiceName = "INIT_SERVICE_NAME"
	EnvConfigDir   = "INIT_CONFIG_DIR"
	EnvEndpointURL = "SERVICE_ENDPOINT_URL"
	EnvRegion      = "AWS_DEFAULT_REGION"
)

// ProcessInit processes initialization for a service
func ProcessInit(serviceName, configDir, endpointURL, region string) error {
	ctx := context.Background()

	// Wait for service readiness
	if err := waitForService(ctx, serviceName, endpointURL); err != nil {
		return fmt.Errorf("service not ready: %w", err)
	}

	// Process configs
	return processConfigs(ctx, serviceName, configDir, endpointURL, region)
}

func waitForService(ctx context.Context, serviceName, endpointURL string) error {
	switch serviceName {
	case "localstack":
		return waitForHTTP(ctx, endpointURL+"/_localstack/health")
	case "postgres":
		return waitForPostgres(ctx)
	case "kafka":
		return waitForKafka(ctx, endpointURL)
	default:
		return nil
	}
}

func waitForHTTP(ctx context.Context, url string) error {
	client := &http.Client{Timeout: HTTPTimeout}
	for range MaxRetries {
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		if resp, err := client.Do(req); err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == HTTPStatusOK {
				return nil
			}
		}
		time.Sleep(RetryDelay)
	}
	return fmt.Errorf("service not ready after %v", ServiceTimeout)
}

func waitForPostgres(ctx context.Context) error {
	const (
		retryDelay = 2 * time.Second
		maxRetries = 30
	)

	for range maxRetries {
		cmd := exec.CommandContext(ctx, "pg_isready", "-h", "localhost", "-p", "5432")
		if err := cmd.Run(); err == nil {
			return nil
		}
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("postgres not ready after 60s")
}

func waitForKafka(ctx context.Context, endpoint string) error {
	const (
		retryDelay = 2 * time.Second
		maxRetries = 30
	)

	for range maxRetries {
		cmd := exec.CommandContext(ctx, "kafka-topics.sh", "--bootstrap-server", endpoint, "--list")
		if err := cmd.Run(); err == nil {
			return nil
		}
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("kafka not ready after 60s")
}

func processConfigs(ctx context.Context, serviceName, configDir, endpointURL, region string) error {
	pattern := filepath.Join(configDir, serviceName+"-*.yml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		filename := filepath.Base(match)
		configType := strings.TrimSuffix(strings.TrimPrefix(filename, serviceName+"-"), ".yml")

		if err := processConfig(ctx, serviceName, match, configType, endpointURL, region); err != nil {
			fmt.Printf("Warning: Failed to process %s: %v\n", match, err)
		}
	}
	return nil
}

func processConfig(ctx context.Context, serviceName, configFile, configType, endpointURL, region string) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	switch serviceName {
	case "localstack":
		return processLocalStack(ctx, configType, config, endpointURL, region)
	case "postgres":
		return processPostgres(ctx, configType, config)
	case "kafka":
		return processKafka(ctx, configType, config, endpointURL)
	default:
		return fmt.Errorf("unsupported service: %s", serviceName)
	}
}

func processLocalStack(ctx context.Context, configType string, config map[string]any, endpointURL, region string) error {
	switch configType {
	case "sqs":
		return processSQSQueues(ctx, config, endpointURL, region)
	case "sns":
		return processSNSTopics(ctx, config, endpointURL, region)
	case "s3":
		return processS3Buckets(ctx, config, endpointURL, region)
	case "dynamodb":
		return processDynamoDBTables(ctx, config, endpointURL, region)
	}
	return nil
}

func processSQSQueues(ctx context.Context, config map[string]any, endpointURL, region string) error {
	queues, ok := config["queues"].([]any)
	if !ok {
		return nil
	}

	for _, q := range queues {
		if queue, ok := q.(map[string]any); ok && queue["name"] != nil {
			if name, ok := queue["name"].(string); ok {
				cmd := exec.CommandContext(ctx, "aws", "--endpoint-url", endpointURL, "--region", region, "sqs", "create-queue", "--queue-name", name)
				if err := cmd.Run(); err != nil {
					fmt.Printf("Warning: Failed to create SQS queue %s: %v\n", name, err)
				}
			}
		}
	}
	return nil
}

func processSNSTopics(ctx context.Context, config map[string]any, endpointURL, region string) error {
	topics, ok := config["topics"].([]any)
	if !ok {
		return nil
	}

	for _, t := range topics {
		if topic, ok := t.(map[string]any); ok && topic["name"] != nil {
			if name, ok := topic["name"].(string); ok {
				cmd := exec.CommandContext(ctx, "aws", "--endpoint-url", endpointURL, "--region", region, "sns", "create-topic", "--name", name)
				if err := cmd.Run(); err != nil {
					fmt.Printf("Warning: Failed to create SNS topic %s: %v\n", name, err)
				}
			}
		}
	}
	return nil
}

func processS3Buckets(ctx context.Context, config map[string]any, endpointURL, region string) error {
	buckets, ok := config["buckets"].([]any)
	if !ok {
		return nil
	}

	for _, b := range buckets {
		if bucket, ok := b.(map[string]any); ok && bucket["name"] != nil {
			if name, ok := bucket["name"].(string); ok {
				cmd := exec.CommandContext(ctx, "aws", "--endpoint-url", endpointURL, "--region", region, "s3", "mb", "s3://"+name)
				if err := cmd.Run(); err != nil {
					fmt.Printf("Warning: Failed to create S3 bucket %s: %v\n", name, err)
				}
			}
		}
	}
	return nil
}

func processDynamoDBTables(ctx context.Context, config map[string]any, endpointURL, region string) error {
	tables, ok := config["tables"].([]any)
	if !ok {
		return nil
	}

	for _, t := range tables {
		if table, ok := t.(map[string]any); ok && table["name"] != nil {
			if name, ok := table["name"].(string); ok {
				if err := createDynamoDBTable(ctx, table, name, endpointURL, region); err != nil {
					fmt.Printf("Warning: Failed to create DynamoDB table %s: %v\n", name, err)
				}
			}
		}
	}
	return nil
}

func createDynamoDBTable(ctx context.Context, table map[string]any, name, endpointURL, region string) error {
	keySchema := `[{"AttributeName": "id", "KeyType": "HASH"}]`
	attrDefs := `[{"AttributeName": "id", "AttributeType": "S"}]`
	billingMode := "PAY_PER_REQUEST"

	if ks, exists := table["key_schema"]; exists {
		if ksBytes, err := json.Marshal(ks); err == nil {
			keySchema = string(ksBytes)
		}
	}
	if ad, exists := table["attribute_definitions"]; exists {
		if adBytes, err := json.Marshal(ad); err == nil {
			attrDefs = string(adBytes)
		}
	}
	if bm, exists := table["billing_mode"].(string); exists {
		billingMode = bm
	}

	cmd := exec.CommandContext(ctx, "aws", "--endpoint-url", endpointURL, "--region", region,
		"dynamodb", "create-table", "--table-name", name, "--key-schema", keySchema,
		"--attribute-definitions", attrDefs, "--billing-mode", billingMode)
	return cmd.Run()
}

func processPostgres(ctx context.Context, configType string, config map[string]any) error {
	switch configType {
	case "schemas":
		if schemas, ok := config["schemas"].([]any); ok {
			for _, s := range schemas {
				if schema, ok := s.(map[string]any); ok && schema["name"] != nil {
					if name, ok := schema["name"].(string); ok {
						cmd := exec.CommandContext(ctx, "psql", "-h", "localhost", "-U", "postgres", "-c", fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", name))
						if err := cmd.Run(); err != nil {
							fmt.Printf("Warning: Failed to create schema %s: %v\n", name, err)
						}
					}
				}
			}
		}
	case "databases":
		if databases, ok := config["databases"].([]any); ok {
			for _, d := range databases {
				if database, ok := d.(map[string]any); ok && database["name"] != nil {
					if name, ok := database["name"].(string); ok {
						cmd := exec.CommandContext(ctx, "createdb", "-h", "localhost", "-U", "postgres", name)
						if err := cmd.Run(); err != nil {
							fmt.Printf("Warning: Failed to create database %s: %v\n", name, err)
						}
					}
				}
			}
		}
	}
	return nil
}

func processKafka(ctx context.Context, configType string, config map[string]any, endpoint string) error {
	if configType == "topics" {
		if topics, ok := config["topics"].([]any); ok {
			for _, t := range topics {
				if topic, ok := t.(map[string]any); ok && topic["name"] != nil {
					if name, ok := topic["name"].(string); ok {
						partitions := "3"
						replicationFactor := "1"

						if p, exists := topic["partitions"]; exists {
							partitions = fmt.Sprintf("%v", p)
						}
						if rf, exists := topic["replication_factor"]; exists {
							replicationFactor = fmt.Sprintf("%v", rf)
						}

						cmd := exec.CommandContext(ctx, "kafka-topics.sh", "--bootstrap-server", endpoint, "--create", "--topic", name, "--partitions", partitions, "--replication-factor", replicationFactor, "--if-not-exists")
						if err := cmd.Run(); err != nil {
							fmt.Printf("Warning: Failed to create Kafka topic %s: %v\n", name, err)
						}
					}
				}
			}
		}
	}
	return nil
}

// GenericInitScript provides a shell script that calls the Go processor
var GenericInitScript = `#!/bin/sh
set -e

# Install required tools
apk add --no-cache curl wget aws-cli postgresql-client yq > /dev/null 2>&1

# Install kafka tools if needed
if [ "$INIT_SERVICE_NAME" = "kafka" ]; then
    wget -q https://downloads.apache.org/kafka/2.8.0/kafka_2.13-2.8.0.tgz
    tar -xzf kafka_2.13-2.8.0.tgz
    mv kafka_2.13-2.8.0/bin/kafka-*.sh /usr/local/bin/
    rm -rf kafka_2.13-2.8.0*
fi

echo "Starting initialization for service: $INIT_SERVICE_NAME"
echo "Configuration directory: $INIT_CONFIG_DIR"
echo "Service endpoint: $SERVICE_ENDPOINT_URL"

# Wait for service readiness
wait_for_service() {
    case "$INIT_SERVICE_NAME" in
        "localstack")
            echo "Waiting for LocalStack to be ready..."
            for i in $(seq 1 30); do
                if curl -f "$SERVICE_ENDPOINT_URL/_localstack/health" >/dev/null 2>&1; then
                    echo "LocalStack is ready"
                    return 0
                fi
                sleep 2
            done
            echo "Warning: LocalStack not ready after 60s"
            ;;
        "postgres")
            echo "Waiting for PostgreSQL to be ready..."
            for i in $(seq 1 30); do
                if pg_isready -h localhost -p 5432 >/dev/null 2>&1; then
                    echo "PostgreSQL is ready"
                    return 0
                fi
                sleep 2
            done
            echo "Warning: PostgreSQL not ready after 60s"
            ;;
        "kafka")
            echo "Waiting for Kafka to be ready..."
            for i in $(seq 1 30); do
                if kafka-topics.sh --bootstrap-server "$SERVICE_ENDPOINT_URL" --list >/dev/null 2>&1; then
                    echo "Kafka is ready"
                    return 0
                fi
                sleep 2
            done
            echo "Warning: Kafka not ready after 60s"
            ;;
    esac
}

# Process configuration files
process_configs() {
    cd "$INIT_CONFIG_DIR" || return 1
    
    for config_file in ${INIT_SERVICE_NAME}-*.yml; do
        [ -f "$config_file" ] || continue
        
        config_type=$(echo "$config_file" | sed "s/${INIT_SERVICE_NAME}-//;s/.yml//")
        echo "Processing $config_file (type: $config_type)"
        
        case "$INIT_SERVICE_NAME" in
            "localstack")
                process_localstack_config "$config_file" "$config_type"
                ;;
            "postgres")
                process_postgres_config "$config_file" "$config_type"
                ;;
            "kafka")
                process_kafka_config "$config_file" "$config_type"
                ;;
        esac
    done
}

# Process LocalStack configurations
process_localstack_config() {
    local file="$1"
    local type="$2"
    
    case "$type" in
        "s3")
            yq eval '.buckets[]?.name' "$file" | while read -r bucket; do
                [ -n "$bucket" ] && [ "$bucket" != "null" ] || continue
                echo "Creating S3 bucket: $bucket"
                aws --endpoint-url "$SERVICE_ENDPOINT_URL" --region "${AWS_DEFAULT_REGION:-us-east-1}" s3 mb "s3://$bucket" 2>/dev/null || echo "Warning: Failed to create bucket $bucket"
            done
            ;;
        "sqs")
            yq eval '.queues[]?.name' "$file" | while read -r queue; do
                [ -n "$queue" ] && [ "$queue" != "null" ] || continue
                echo "Creating SQS queue: $queue"
                aws --endpoint-url "$SERVICE_ENDPOINT_URL" --region "${AWS_DEFAULT_REGION:-us-east-1}" sqs create-queue --queue-name "$queue" >/dev/null 2>&1 || echo "Warning: Failed to create queue $queue"
            done
            ;;
        "sns")
            yq eval '.topics[]?.name' "$file" | while read -r topic; do
                [ -n "$topic" ] && [ "$topic" != "null" ] || continue
                echo "Creating SNS topic: $topic"
                aws --endpoint-url "$SERVICE_ENDPOINT_URL" --region "${AWS_DEFAULT_REGION:-us-east-1}" sns create-topic --name "$topic" >/dev/null 2>&1 || echo "Warning: Failed to create topic $topic"
            done
            ;;
        "dynamodb")
            yq eval '.tables[]?.name' "$file" | while read -r table; do
                [ -n "$table" ] && [ "$table" != "null" ] || continue
                echo "Creating DynamoDB table: $table"
                # Simple table creation - complex schemas would need Go implementation
                aws --endpoint-url "$SERVICE_ENDPOINT_URL" --region "${AWS_DEFAULT_REGION:-us-east-1}" dynamodb create-table \
                    --table-name "$table" \
                    --attribute-definitions AttributeName=id,AttributeType=S \
                    --key-schema AttributeName=id,KeyType=HASH \
                    --billing-mode PAY_PER_REQUEST >/dev/null 2>&1 || echo "Warning: Failed to create table $table"
            done
            ;;
    esac
}

# Process PostgreSQL configurations
process_postgres_config() {
    local file="$1"
    local type="$2"
    
    case "$type" in
        "databases")
            yq eval '.databases[]?.name' "$file" | while read -r db; do
                [ -n "$db" ] && [ "$db" != "null" ] || continue
                echo "Creating database: $db"
                psql -h localhost -U postgres -c "CREATE DATABASE \"$db\";" 2>/dev/null || echo "Warning: Failed to create database $db"
            done
            ;;
    esac
}

# Process Kafka configurations
process_kafka_config() {
    local file="$1"
    local type="$2"
    
    case "$type" in
        "topics")
            yq eval '.topics[]?.name' "$file" | while read -r topic; do
                [ -n "$topic" ] && [ "$topic" != "null" ] || continue
                partitions=$(yq eval ".topics[] | select(.name == \"$topic\") | .partitions // 1" "$file")
                echo "Creating Kafka topic: $topic (partitions: $partitions)"
                kafka-topics.sh --bootstrap-server "$SERVICE_ENDPOINT_URL" --create --topic "$topic" --partitions "$partitions" --replication-factor 1 2>/dev/null || echo "Warning: Failed to create topic $topic"
            done
            ;;
    esac
}

# Main execution
wait_for_service
process_configs

echo "Initialization complete for $INIT_SERVICE_NAME!"
`
