package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	DefaultVisibilityTimeout  = 30
	DefaultMessageRetention   = 1209600
	DefaultMaxReceiveCount    = 3
	DefaultRawMessageDelivery = true
)

// LocalStackConfig represents the configuration for LocalStack resource creation
type LocalStackConfig struct {
	SQSQueues      []SQSQueue      `json:"sqs_queues,omitempty"`
	SNSTopics      []SNSTopic      `json:"sns_topics,omitempty"`
	DynamoDBTables []DynamoDBTable `json:"dynamodb_tables,omitempty"`
}

// SQSQueue represents an SQS queue configuration
type SQSQueue struct {
	Name                          string `json:"name"`
	VisibilityTimeoutSeconds      int    `json:"visibility_timeout_seconds,omitempty"`
	MessageRetentionPeriod        int    `json:"message_retention_period,omitempty"`
	ReceiveMessageWaitTimeSeconds int    `json:"receive_message_wait_time_seconds,omitempty"`
	DelaySeconds                  int    `json:"delay_seconds,omitempty"`
	DeadLetterQueue               string `json:"dead_letter_queue,omitempty"`
	MaxReceiveCount               int    `json:"max_receive_count,omitempty"`
}

// SNSTopic represents an SNS topic configuration
type SNSTopic struct {
	Name          string         `json:"name"`
	DisplayName   string         `json:"display_name,omitempty"`
	Subscriptions []Subscription `json:"subscriptions,omitempty"`
}

// Subscription represents an SNS subscription
type Subscription struct {
	Protocol           string `json:"protocol"`
	Endpoint           string `json:"endpoint"`
	RawMessageDelivery bool   `json:"raw_message_delivery,omitempty"`
	FilterPolicy       string `json:"filter_policy,omitempty"`
}

// DynamoDBTable represents a DynamoDB table configuration
type DynamoDBTable struct {
	Name                   string                 `json:"name"`
	AttributeDefinitions   []AttributeDefinition  `json:"attribute_definitions"`
	KeySchema              []KeySchemaElement     `json:"key_schema"`
	ProvisionedThroughput  ProvisionedThroughput  `json:"provisioned_throughput"`
	GlobalSecondaryIndexes []GlobalSecondaryIndex `json:"global_secondary_indexes,omitempty"`
	TableClass             string                 `json:"table_class,omitempty"`
}

// AttributeDefinition represents a DynamoDB attribute definition
type AttributeDefinition struct {
	AttributeName string `json:"AttributeName"`
	AttributeType string `json:"AttributeType"`
}

// KeySchemaElement represents a DynamoDB key schema element
type KeySchemaElement struct {
	AttributeName string `json:"AttributeName"`
	KeyType       string `json:"KeyType"`
}

// ProvisionedThroughput represents DynamoDB provisioned throughput
type ProvisionedThroughput struct {
	ReadCapacityUnits  int `json:"ReadCapacityUnits"`
	WriteCapacityUnits int `json:"WriteCapacityUnits"`
}

// GlobalSecondaryIndex represents a DynamoDB GSI
type GlobalSecondaryIndex struct {
	IndexName             string                `json:"IndexName"`
	KeySchema             []KeySchemaElement    `json:"KeySchema"`
	Projection            Projection            `json:"Projection"`
	ProvisionedThroughput ProvisionedThroughput `json:"ProvisionedThroughput"`
}

// Projection represents a DynamoDB projection
type Projection struct {
	ProjectionType string `json:"ProjectionType"`
}

// ConvertToLocalStackConfig converts otto-stack YAML config to LocalStack JSON format
func ConvertToLocalStackConfig(yamlData []byte) ([]byte, error) {
	var config map[string]any
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	localStackConfig := LocalStackConfig{}

	// Extract service configuration
	serviceConfig, ok := config["service_configuration"].(map[string]any)
	if !ok {
		return json.Marshal(localStackConfig)
	}

	// Convert SQS queues
	localStackConfig.SQSQueues = convertSQSQueues(serviceConfig)

	// Convert SNS topics
	localStackConfig.SNSTopics = convertSNSTopics(serviceConfig)

	return json.Marshal(localStackConfig)
}

func convertSQSQueues(serviceConfig map[string]any) []SQSQueue {
	var queues []SQSQueue

	sqsConfig, exists := serviceConfig["localstack-sqs"].(map[string]any)
	if !exists {
		return queues
	}

	queueList, exists := sqsConfig["queues"].([]any)
	if !exists {
		return queues
	}

	for _, q := range queueList {
		queueMap, ok := q.(map[string]any)
		if !ok {
			continue
		}

		queue := SQSQueue{
			Name:                          getString(queueMap, "name"),
			VisibilityTimeoutSeconds:      getInt(queueMap, "visibility_timeout", DefaultVisibilityTimeout),
			MessageRetentionPeriod:        getInt(queueMap, "message_retention_period", DefaultMessageRetention),
			ReceiveMessageWaitTimeSeconds: getInt(queueMap, "receive_message_wait_time", 0),
			DelaySeconds:                  getInt(queueMap, "delay_seconds", 0),
			DeadLetterQueue:               getString(queueMap, "dead_letter_queue"),
			MaxReceiveCount:               getInt(queueMap, "max_receive_count", DefaultMaxReceiveCount),
		}
		queues = append(queues, queue)
	}

	return queues
}

func convertSNSTopics(serviceConfig map[string]any) []SNSTopic {
	var topics []SNSTopic

	snsConfig, exists := serviceConfig["localstack-sns"].(map[string]any)
	if !exists {
		return topics
	}

	topicList, exists := snsConfig["topics"].([]any)
	if !exists {
		return topics
	}

	for _, t := range topicList {
		topicMap, ok := t.(map[string]any)
		if !ok {
			continue
		}

		topic := SNSTopic{
			Name:        getString(topicMap, "name"),
			DisplayName: getString(topicMap, "display_name"),
		}

		// Convert subscriptions
		if subs, exists := topicMap["subscriptions"].([]any); exists {
			for _, s := range subs {
				if subMap, ok := s.(map[string]any); ok {
					sub := Subscription{
						Protocol:           getString(subMap, "protocol"),
						Endpoint:           getString(subMap, "endpoint"),
						RawMessageDelivery: getBool(subMap, "raw_message_delivery", DefaultRawMessageDelivery),
						FilterPolicy:       getString(subMap, "filter_policy"),
					}
					topic.Subscriptions = append(topic.Subscriptions, sub)
				}
			}
		}

		topics = append(topics, topic)
	}

	return topics
}

// Helper functions
func getString(m map[string]any, key string) string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getInt(m map[string]any, key string, defaultVal int) int {
	if val, exists := m[key]; exists {
		if i, ok := val.(int); ok {
			return i
		}
		if f, ok := val.(float64); ok {
			return int(f)
		}
	}
	return defaultVal
}

func getBool(m map[string]any, key string, defaultVal bool) bool {
	if val, exists := m[key]; exists {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultVal
}
