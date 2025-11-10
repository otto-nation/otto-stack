#!/bin/bash

# Configuration Bridge Script
# Converts otto-stack YAML config to LocalStack JSON format

set -e

CONFIG_FILE="${1:-/config/otto-stack.yaml}"
OUTPUT_FILE="${2:-/config/localstack-config.json}"

if [ ! -f "$CONFIG_FILE" ]; then
    echo "Configuration file not found: $CONFIG_FILE"
    exit 1
fi

# Extract service configuration and convert to JSON
yq eval '
{
  "sqs_queues": (.service_configuration."localstack-sqs".queues // [] | map({
    "name": .name,
    "visibility_timeout_seconds": (.visibility_timeout // 30),
    "dead_letter_queue": .dead_letter_queue,
    "max_receive_count": (.max_receive_count // 3)
  })),
  "sns_topics": (.service_configuration."localstack-sns".topics // [] | map({
    "name": .name,
    "subscriptions": (.subscriptions // [] | map({
      "protocol": .protocol,
      "endpoint": .endpoint,
      "raw_message_delivery": true
    }))
  }))
}
' "$CONFIG_FILE" > "$OUTPUT_FILE"

echo "Configuration converted: $CONFIG_FILE -> $OUTPUT_FILE"
