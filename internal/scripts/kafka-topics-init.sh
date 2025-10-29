#!/bin/bash

# Kafka Topic Initialization Script
# This script creates Kafka topics based on otto-stack configuration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
KAFKA_BROKER="kafka-broker:29092"
CONFIG_FILE="/config/otto-stack-config.yml"
DEFAULT_PARTITIONS=3
DEFAULT_REPLICATION_FACTOR=1

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Wait for Kafka to be ready
wait_for_kafka() {
    log_info "Waiting for Kafka to be ready..."

    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if kafka-broker-api-versions --bootstrap-server "$KAFKA_BROKER" > /dev/null 2>&1; then
            log_success "Kafka is ready!"
            return 0
        fi

        log_info "Attempt $attempt/$max_attempts - Kafka not ready yet..."
        sleep 2
        ((attempt++))
    done

    log_error "Kafka failed to start within expected time"
    return 1
}

# Check if configuration file exists and has Kafka services
check_config() {
    if [ ! -f "$CONFIG_FILE" ]; then
        log_warning "No otto-stack configuration found at $CONFIG_FILE"
        return 1
    fi

    # Check if Kafka services are enabled
    local kafka_enabled=$(yq '.stack.enabled[]' "$CONFIG_FILE" 2>/dev/null | grep -E "kafka" || true)
    
    if [ -z "$kafka_enabled" ]; then
        log_info "No Kafka services enabled, skipping topic creation"
        return 1
    fi

    return 0
}

# Create Kafka topics
create_kafka_topics() {
    log_info "Creating Kafka topics..."

    # Get project name for prefixing
    local project_name=$(yq '.project.name' "$CONFIG_FILE")
    
    # Get topic names from service configuration or use defaults
    local topic_names=$(yq '.service-configuration.kafka-topics.topic_names // ["events", "notifications", "audit-logs"]' "$CONFIG_FILE")
    
    # Get topic configuration
    local partitions=$(yq '.service-configuration.kafka-topics.partitions // 3' "$CONFIG_FILE")
    local replication_factor=$(yq '.service-configuration.kafka-topics.replication_factor // 1' "$CONFIG_FILE")
    
    # Convert YAML array to bash array
    local topics=()
    while IFS= read -r topic; do
        topic=$(echo "$topic" | sed 's/^"//;s/"$//')
        topics+=("$topic")
    done < <(echo "$topic_names" | yq '.[]')

    if [ ${#topics[@]} -eq 0 ]; then
        log_info "No Kafka topics configured"
        return 0
    fi

    log_info "Found ${#topics[@]} Kafka topic(s) to create"

    # Create each topic
    for topic_name in "${topics[@]}"; do
        local full_topic_name="${project_name}-${topic_name}"
        
        log_info "Creating Kafka topic: $full_topic_name (partitions: $partitions, replication: $replication_factor)"
        
        if kafka-topics --bootstrap-server "$KAFKA_BROKER" \
           --create --topic "$full_topic_name" \
           --partitions "$partitions" \
           --replication-factor "$replication_factor" \
           --if-not-exists > /dev/null 2>&1; then
            log_success "Created Kafka topic: $full_topic_name"
        else
            log_error "Failed to create Kafka topic: $full_topic_name"
        fi
    done
}

# Main execution
main() {
    log_info "Starting Kafka topic initialization..."

    # Check configuration
    if ! check_config; then
        log_info "Exiting - no Kafka services to initialize"
        exit 0
    fi

    # Wait for Kafka to be ready
    if ! wait_for_kafka; then
        log_error "Kafka is not ready, exiting"
        exit 1
    fi

    # Create topics
    create_kafka_topics

    log_success "Kafka topic initialization completed!"
}

# Run main function
main "$@"

        log_info "Attempt $attempt/$max_attempts - Kafka not ready yet..."
        sleep 2
        ((attempt++))
    done

    log_error "Kafka failed to start within expected time"
    return 1
}

# Check if configuration file exists
check_config() {
    if [ ! -f "$CONFIG_FILE" ]; then
        log_warning "No topics configuration found at $CONFIG_FILE"
        log_info "Creating default topics only"
        return 1
    fi

    if [ ! -s "$CONFIG_FILE" ]; then
        log_warning "Topics configuration file is empty"
        return 1
    fi

    return 0
}

# Create default topics
create_default_topics() {
    log_info "Creating default topics..."

    local default_topics=(
        "test:3:1"
        "events:6:1"
        "user-events:3:1"
        "notifications:3:1"
    )

    for topic_spec in "${default_topics[@]}"; do
        IFS=':' read -r topic_name partitions replication_factor <<< "$topic_spec"

        log_info "Creating default topic: $topic_name (partitions: $partitions, replication-factor: $replication_factor)"

        kafka-topics --create --if-not-exists \
            --bootstrap-server "$KAFKA_BROKER" \
            --topic "$topic_name" \
            --partitions "$partitions" \
            --replication-factor "$replication_factor"

        log_success "Created default topic: $topic_name"
    done
}

# Create custom topics from configuration
create_custom_topics() {
    log_info "Creating custom topics from configuration..."

    # Check if topics array exists in config
    if ! jq -e '.topics' "$CONFIG_FILE" > /dev/null 2>&1; then
        log_info "No custom topics configured"
        return 0
    fi

    # Get topic count
    local topic_count=$(jq '.topics | length' "$CONFIG_FILE")

    if [ "$topic_count" -eq 0 ]; then
        log_info "No custom topics configured"
        return 0
    fi

    log_info "Found $topic_count custom topic(s) to create"

    # Create each topic
    for i in $(seq 0 $((topic_count - 1))); do
        local topic_config=$(jq ".topics[$i]" "$CONFIG_FILE")
        local topic_name=$(echo "$topic_config" | jq -r '.name')
        local partitions=$(echo "$topic_config" | jq -r ".partitions // $DEFAULT_PARTITIONS")
        local replication_factor=$(echo "$topic_config" | jq -r ".replication_factor // $DEFAULT_REPLICATION_FACTOR")
        local cleanup_policy=$(echo "$topic_config" | jq -r '.cleanup_policy // "delete"')
        local retention_ms=$(echo "$topic_config" | jq -r '.retention_ms // empty')

        log_info "Creating custom topic: $topic_name (partitions: $partitions, replication-factor: $replication_factor)"

        # Create the topic
        kafka-topics --create --if-not-exists \
            --bootstrap-server "$KAFKA_BROKER" \
            --topic "$topic_name" \
            --partitions "$partitions" \
            --replication-factor "$replication_factor"

        # Apply additional configurations if specified
        local configs=()

        if [ "$cleanup_policy" != "delete" ] && [ "$cleanup_policy" != "null" ]; then
            configs+=("cleanup.policy=$cleanup_policy")
        fi

        if [ "$retention_ms" != "null" ] && [ "$retention_ms" != "" ]; then
            configs+=("retention.ms=$retention_ms")
        fi

        # Apply configurations if any
        if [ ${#configs[@]} -gt 0 ]; then
            local config_string=$(IFS=','; echo "${configs[*]}")
            log_info "Applying configurations to topic $topic_name: $config_string"

            kafka-configs --bootstrap-server "$KAFKA_BROKER" \
                --entity-type topics \
                --entity-name "$topic_name" \
                --alter \
                --add-config "$config_string"
        fi

        log_success "Created custom topic: $topic_name"
    done
}

# List all topics
list_topics() {
    log_info "Listing all created topics..."
    echo ""

    local topics=$(kafka-topics --list --bootstrap-server "$KAFKA_BROKER")

    if [ -n "$topics" ]; then
        log_info "📋 Available Topics:"
        echo "$topics" | while read -r topic; do
            if [ -n "$topic" ]; then
                # Get topic details
                local details=$(kafka-topics --describe --bootstrap-server "$KAFKA_BROKER" --topic "$topic" 2>/dev/null | head -1)
                echo "  • $topic"
                if [[ "$details" =~ PartitionCount:([0-9]+) ]]; then
                    local partition_count="${BASH_REMATCH[1]}"
                    echo "    Partitions: $partition_count"
                fi
                if [[ "$details" =~ ReplicationFactor:([0-9]+) ]]; then
                    local repl_factor="${BASH_REMATCH[1]}"
                    echo "    Replication Factor: $repl_factor"
                fi
                echo ""
            fi
        done
    else
        log_warning "No topics found"
    fi
}

# Verify topic creation
verify_topics() {
    log_info "Verifying topic creation..."

    if check_config; then
        local topic_count=$(jq '.topics | length' "$CONFIG_FILE")

        for i in $(seq 0 $((topic_count - 1))); do
            local topic_name=$(jq -r ".topics[$i].name" "$CONFIG_FILE")

            if kafka-topics --list --bootstrap-server "$KAFKA_BROKER" | grep -q "^$topic_name$"; then
                log_success "✓ Topic '$topic_name' created successfully"
            else
                log_error "✗ Topic '$topic_name' was not created"
            fi
        done
    fi
}

# Main execution
main() {
    log_info "🚀 Starting Kafka topic initialization"

    # Wait for Kafka to be ready
    if ! wait_for_kafka; then
        exit 1
    fi

    # Check if we have custom configuration
    if check_config; then
        log_info "Using custom topics configuration from: $CONFIG_FILE"
        create_custom_topics
    else
        # Create default topics if no configuration
        create_default_topics
    fi

    # List all created topics
    list_topics

    # Verify custom topics were created
    if check_config; then
        verify_topics
    fi

    log_success "🎉 Kafka topic initialization completed!"
}

# Run main function
main "$@"
