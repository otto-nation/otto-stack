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
	const (
		httpTimeout = 5 * time.Second
		retryDelay  = 2 * time.Second
		maxRetries  = 30
		httpOK      = 200
	)

	client := &http.Client{Timeout: httpTimeout}
	for range maxRetries {
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		if resp, err := client.Do(req); err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == httpOK {
				return nil
			}
		}
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("service not ready after 60s")
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
		if queues, ok := config["queues"].([]any); ok {
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
		}
	case "sns":
		if topics, ok := config["topics"].([]any); ok {
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
		}
	case "s3":
		if buckets, ok := config["buckets"].([]any); ok {
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
		}
	case "dynamodb":
		if tables, ok := config["tables"].([]any); ok {
			for _, t := range tables {
				if table, ok := t.(map[string]any); ok && table["name"] != nil {
					if name, ok := table["name"].(string); ok {
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

						cmd := exec.CommandContext(ctx, "aws", "--endpoint-url", endpointURL, "--region", region, "dynamodb", "create-table", "--table-name", name, "--key-schema", keySchema, "--attribute-definitions", attrDefs, "--billing-mode", billingMode)
						if err := cmd.Run(); err != nil {
							fmt.Printf("Warning: Failed to create DynamoDB table %s: %v\n", name, err)
						}
					}
				}
			}
		}
	}
	return nil
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
apk add --no-cache curl wget aws-cli postgresql-client > /dev/null 2>&1

# Install kafka tools if needed
if [ "$INIT_SERVICE_NAME" = "kafka" ]; then
    wget -q https://downloads.apache.org/kafka/2.8.0/kafka_2.13-2.8.0.tgz
    tar -xzf kafka_2.13-2.8.0.tgz
    mv kafka_2.13-2.8.0/bin/kafka-*.sh /usr/local/bin/
    rm -rf kafka_2.13-2.8.0*
fi

echo "Starting auto-discovery initialization for service: $INIT_SERVICE_NAME"
echo "Configuration directory: $CONFIG_DIR"
echo "Service endpoint: $SERVICE_ENDPOINT_URL"

# Simple shell-based processing for now
# This will be replaced with Go binary in the future
echo "Auto-discovery initialization complete for $INIT_SERVICE_NAME!"
`
