#!/bin/bash

# LocalStack AWS Resource Initialization Script
# This script creates AWS resources based on otto-stack configuration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
AWS_ENDPOINT="http://localstack-core:4566"
AWS_REGION="us-east-1"
CONFIG_FILE="/config/otto-stack-config.yml"

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

# Wait for LocalStack to be ready
wait_for_localstack() {
    log_info "Waiting for LocalStack to be ready..."

    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if curl -s "${AWS_ENDPOINT}/health" > /dev/null 2>&1; then
            log_success "LocalStack is ready!"
            return 0
        fi

        log_info "Attempt $attempt/$max_attempts - LocalStack not ready yet..."
        sleep 2
        ((attempt++))
    done

    log_error "LocalStack failed to start within expected time"
    return 1
}

# Check if configuration file exists and has LocalStack services
check_config() {
    if [ ! -f "$CONFIG_FILE" ]; then
        log_warning "No otto-stack configuration found at $CONFIG_FILE"
        return 1
    fi

    # Check if any LocalStack services are enabled
    local enabled_services=$(yq '.stack.enabled[]' "$CONFIG_FILE" 2>/dev/null | grep -E "localstack|sqs|sns|s3" || true)
    
    if [ -z "$enabled_services" ]; then
        log_info "No LocalStack services enabled, skipping initialization"
        return 1
    fi

    return 0
}

# Create SQS queues
create_sqs_queues() {
    # Check if SQS services are enabled
    local sqs_enabled=$(yq '.stack.enabled[]' "$CONFIG_FILE" | grep -E "localstack-sqs|sqs" || true)
    
    if [ -z "$sqs_enabled" ]; then
        log_info "SQS not enabled, skipping queue creation"
        return 0
    fi

    log_info "Creating SQS queues..."

    # Get project name for prefixing
    local project_name=$(yq '.project.name' "$CONFIG_FILE")
    
    # Get queue names from service configuration or use defaults
    local queue_names=$(yq '.service-configuration.localstack-sqs.queue_names // ["default-queue", "test-queue"]' "$CONFIG_FILE")
    
    # Convert YAML array to bash array
    local queues=()
    while IFS= read -r queue; do
        # Remove quotes and add to array
        queue=$(echo "$queue" | sed 's/^"//;s/"$//')
        queues+=("$queue")
    done < <(echo "$queue_names" | yq '.[]')

    if [ ${#queues[@]} -eq 0 ]; then
        log_info "No SQS queues configured"
        return 0
    fi

    log_info "Found ${#queues[@]} SQS queue(s) to create"

    # Create each queue
    for queue_name in "${queues[@]}"; do
        local full_queue_name="${project_name}-${queue_name}"
        
        log_info "Creating SQS queue: $full_queue_name"
        
        if aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
           sqs create-queue --queue-name "$full_queue_name" > /dev/null 2>&1; then
            log_success "Created SQS queue: $full_queue_name"
        else
            log_error "Failed to create SQS queue: $full_queue_name"
        fi
    done
}

# Create SNS topics
create_sns_topics() {
    # Check if SNS services are enabled
    local sns_enabled=$(yq '.stack.enabled[]' "$CONFIG_FILE" | grep -E "localstack-sns|sns" || true)
    
    if [ -z "$sns_enabled" ]; then
        log_info "SNS not enabled, skipping topic creation"
        return 0
    fi

    log_info "Creating SNS topics..."

    # Get project name for prefixing
    local project_name=$(yq '.project.name' "$CONFIG_FILE")
    
    # Get topic names from service configuration or use defaults
    local topic_names=$(yq '.service-configuration.localstack-sns.topic_names // ["default-topic"]' "$CONFIG_FILE")
    
    # Convert YAML array to bash array
    local topics=()
    while IFS= read -r topic; do
        topic=$(echo "$topic" | sed 's/^"//;s/"$//')
        topics+=("$topic")
    done < <(echo "$topic_names" | yq '.[]')

    if [ ${#topics[@]} -eq 0 ]; then
        log_info "No SNS topics configured"
        return 0
    fi

    log_info "Found ${#topics[@]} SNS topic(s) to create"

    # Create each topic
    for topic_name in "${topics[@]}"; do
        local full_topic_name="${project_name}-${topic_name}"
        
        log_info "Creating SNS topic: $full_topic_name"
        
        if aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
           sns create-topic --name "$full_topic_name" > /dev/null 2>&1; then
            log_success "Created SNS topic: $full_topic_name"
        else
            log_error "Failed to create SNS topic: $full_topic_name"
        fi
    done
}

# Create S3 buckets
create_s3_buckets() {
    # Check if S3 services are enabled
    local s3_enabled=$(yq '.stack.enabled[]' "$CONFIG_FILE" | grep -E "localstack-s3|s3" || true)
    
    if [ -z "$s3_enabled" ]; then
        log_info "S3 not enabled, skipping bucket creation"
        return 0
    fi

    log_info "Creating S3 buckets..."

    # Get project name for prefixing
    local project_name=$(yq '.project.name' "$CONFIG_FILE")
    
    # Get bucket names from service configuration or use defaults
    local bucket_names=$(yq '.service-configuration.localstack-s3.bucket_names // ["uploads", "static-assets"]' "$CONFIG_FILE")
    
    # Convert YAML array to bash array
    local buckets=()
    while IFS= read -r bucket; do
        bucket=$(echo "$bucket" | sed 's/^"//;s/"$//')
        buckets+=("$bucket")
    done < <(echo "$bucket_names" | yq '.[]')

    if [ ${#buckets[@]} -eq 0 ]; then
        log_info "No S3 buckets configured"
        return 0
    fi

    log_info "Found ${#buckets[@]} S3 bucket(s) to create"

    # Create each bucket
    for bucket_name in "${buckets[@]}"; do
        local full_bucket_name="${project_name}-${bucket_name}"
        
        log_info "Creating S3 bucket: $full_bucket_name"
        
        if aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
           s3 mb "s3://$full_bucket_name" > /dev/null 2>&1; then
            log_success "Created S3 bucket: $full_bucket_name"
        else
            log_error "Failed to create S3 bucket: $full_bucket_name"
        fi
    done
}

# Main execution
main() {
    log_info "Starting LocalStack AWS resource initialization..."

    # Set AWS credentials for LocalStack
    export AWS_ACCESS_KEY_ID=test
    export AWS_SECRET_ACCESS_KEY=test
    export AWS_DEFAULT_REGION="$AWS_REGION"

    # Check configuration
    if ! check_config; then
        log_info "Exiting - no LocalStack services to initialize"
        exit 0
    fi

    # Wait for LocalStack to be ready
    if ! wait_for_localstack; then
        log_error "LocalStack is not ready, exiting"
        exit 1
    fi

    # Create resources based on enabled services
    create_sqs_queues
    create_sns_topics  
    create_s3_buckets

    log_success "LocalStack AWS resource initialization completed!"
}

# Run main function
main "$@"
                }" > /dev/null

            # Get DLQ ARN for redrive policy
            local dlq_url=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
                sqs get-queue-url --queue-name "$dead_letter_queue" --query 'QueueUrl' --output text)
            local dlq_arn=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
                sqs get-queue-attributes --queue-url "$dlq_url" --attribute-names QueueArn --query 'Attributes.QueueArn' --output text)

            log_success "Created dead letter queue: $dead_letter_queue"
        fi

        # Prepare attributes
        local attributes="{
            \"VisibilityTimeoutSeconds\": \"$visibility_timeout\",
            \"MessageRetentionPeriod\": \"$message_retention_period\",
            \"ReceiveMessageWaitTimeSeconds\": \"$receive_message_wait_time\",
            \"DelaySeconds\": \"$delay_seconds\""

        # Add redrive policy if DLQ is specified
        if [ -n "$dead_letter_queue" ]; then
            attributes="$attributes,
            \"RedrivePolicy\": \"{\\\"deadLetterTargetArn\\\":\\\"$dlq_arn\\\",\\\"maxReceiveCount\\\":$max_receive_count}\""
        fi

        attributes="$attributes}"

        # Create main queue
        aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
            sqs create-queue \
            --queue-name "$queue_name" \
            --attributes "$attributes" > /dev/null

        log_success "Created queue: $queue_name"
    done
}

# Create SNS topics and subscriptions
create_sns_topics() {
    log_info "Creating SNS topics..."

    # Check if sns_topics exists in config
    if ! jq -e '.sns_topics' "$CONFIG_FILE" > /dev/null 2>&1; then
        log_info "No SNS topics configured"
        return 0
    fi

    # Get topic count
    local topic_count=$(jq '.sns_topics | length' "$CONFIG_FILE")

    if [ "$topic_count" -eq 0 ]; then
        log_info "No SNS topics configured"
        return 0
    fi

    log_info "Found $topic_count SNS topic(s) to create"

    # Create each topic
    for i in $(seq 0 $((topic_count - 1))); do
        local topic_config=$(jq ".sns_topics[$i]" "$CONFIG_FILE")
        local topic_name=$(echo "$topic_config" | jq -r '.name')
        local display_name=$(echo "$topic_config" | jq -r '.display_name // empty')

        log_info "Creating topic: $topic_name"

        # Create topic
        local topic_arn=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
            sns create-topic \
            --name "$topic_name" \
            --query 'TopicArn' --output text)

        # Set display name if provided
        if [ "$display_name" != "null" ] && [ "$display_name" != "" ]; then
            aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
                sns set-topic-attributes \
                --topic-arn "$topic_arn" \
                --attribute-name DisplayName \
                --attribute-value "$display_name"
        fi

        log_success "Created topic: $topic_name (ARN: $topic_arn)"

        # Create subscriptions
        local subscriptions=$(echo "$topic_config" | jq -r '.subscriptions // []')
        local subscription_count=$(echo "$subscriptions" | jq 'length')

        if [ "$subscription_count" -gt 0 ]; then
            log_info "Creating $subscription_count subscription(s) for topic: $topic_name"

            for j in $(seq 0 $((subscription_count - 1))); do
                local sub_config=$(echo "$subscriptions" | jq ".[$j]")
                local protocol=$(echo "$sub_config" | jq -r '.protocol')
                local endpoint=$(echo "$sub_config" | jq -r '.endpoint')
                local raw_message_delivery=$(echo "$sub_config" | jq -r '.raw_message_delivery // true')
                local filter_policy=$(echo "$sub_config" | jq -r '.filter_policy // {}')

                # Handle SQS endpoints (convert queue name to ARN)
                if [ "$protocol" = "sqs" ]; then
                    local queue_url=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
                        sqs get-queue-url --queue-name "$endpoint" --query 'QueueUrl' --output text 2>/dev/null || echo "")

                    if [ -z "$queue_url" ]; then
                        log_warning "SQS queue '$endpoint' not found for subscription. Skipping."
                        continue
                    fi

                    local queue_arn=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
                        sqs get-queue-attributes --queue-url "$queue_url" --attribute-names QueueArn --query 'Attributes.QueueArn' --output text)
                    endpoint="$queue_arn"
                fi

                log_info "Creating subscription: $protocol -> $endpoint"

                # Create subscription
                local subscription_arn=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
                    sns subscribe \
                    --topic-arn "$topic_arn" \
                    --protocol "$protocol" \
                    --notification-endpoint "$endpoint" \
                    --query 'SubscriptionArn' --output text)

                # Set raw message delivery if specified
                if [ "$protocol" = "sqs" ]; then
                    aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
                        sns set-subscription-attributes \
                        --subscription-arn "$subscription_arn" \
                        --attribute-name RawMessageDelivery \
                        --attribute-value "$raw_message_delivery"
                fi

                # Set filter policy if specified
                if [ "$filter_policy" != "{}" ] && [ "$filter_policy" != "null" ]; then
                    aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
                        sns set-subscription-attributes \
                        --subscription-arn "$subscription_arn" \
                        --attribute-name FilterPolicy \
                        --attribute-value "$filter_policy"
                fi

                log_success "Created subscription: $protocol -> $(basename "$endpoint")"
            done
        fi
    done
}

# Create DynamoDB tables
create_dynamodb_tables() {
    log_info "Creating DynamoDB tables..."

    # Check if dynamodb_tables exists in config
    if ! jq -e '.dynamodb_tables' "$CONFIG_FILE" > /dev/null 2>&1; then
        log_info "No DynamoDB tables configured"
        return 0
    fi

    # Get table count
    local table_count=$(jq '.dynamodb_tables | length' "$CONFIG_FILE")

    if [ "$table_count" -eq 0 ]; then
        log_info "No DynamoDB tables configured"
        return 0
    fi

    log_info "Found $table_count DynamoDB table(s) to create"

    # List and delete existing tables first
    local existing_tables=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        dynamodb list-tables --query "TableNames[]" --output text 2>/dev/null || echo "")

    if [ -n "$existing_tables" ]; then
        log_info "Deleting existing DynamoDB tables..."
        for table in $existing_tables; do
            if [ -n "$table" ]; then
                log_info "Deleting table: $table"
                aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
                    dynamodb delete-table --table-name "$table" > /dev/null 2>&1 || true
            fi
        done
    fi

    # Create each table
    for i in $(seq 0 $((table_count - 1))); do
        local table_config=$(jq ".dynamodb_tables[$i]" "$CONFIG_FILE")
        local table_name=$(echo "$table_config" | jq -r '.name')

        log_info "Creating DynamoDB table: $table_name"

        # Build AWS CLI command
        local cmd="aws --endpoint-url=$AWS_ENDPOINT --region=$AWS_REGION dynamodb create-table --table-name $table_name"

        # Add attribute definitions
        local attr_defs=$(echo "$table_config" | jq -c '.attribute_definitions // []')
        if [ "$attr_defs" != "[]" ]; then
            cmd="$cmd --attribute-definitions '$attr_defs'"
        fi

        # Add key schema
        local key_schema=$(echo "$table_config" | jq -c '.key_schema // []')
        if [ "$key_schema" != "[]" ]; then
            cmd="$cmd --key-schema '$key_schema'"
        fi

        # Add provisioned throughput
        local throughput=$(echo "$table_config" | jq -c '.provisioned_throughput // {"ReadCapacityUnits": 5, "WriteCapacityUnits": 5}')
        cmd="$cmd --provisioned-throughput '$throughput'"

        # Add global secondary indexes if present
        local gsi=$(echo "$table_config" | jq -c '.global_secondary_indexes // []')
        if [ "$gsi" != "[]" ]; then
            cmd="$cmd --global-secondary-indexes '$gsi'"
        fi

        # Add table class
        local table_class=$(echo "$table_config" | jq -r '.table_class // "STANDARD"')
        cmd="$cmd --table-class $table_class"

        # Execute the command
        eval "$cmd" > /dev/null 2>&1

        if [ $? -eq 0 ]; then
            log_success "Created DynamoDB table: $table_name"
        else
            log_error "Failed to create DynamoDB table: $table_name"
        fi
    done
}

# Create default resources if no config provided
create_default_resources() {
    log_info "Creating default AWS resources..."

    # Create default SQS queues
    log_info "Creating default SQS queues..."

    aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        sqs create-queue --queue-name "default-queue" > /dev/null
    log_success "Created default queue: default-queue"

    aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        sqs create-queue --queue-name "test-queue" > /dev/null
    log_success "Created test queue: test-queue"

    # Create default SNS topic
    log_info "Creating default SNS topics..."

    local topic_arn=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        sns create-topic --name "default-topic" --query 'TopicArn' --output text)
    log_success "Created default topic: default-topic"

    # Subscribe default queue to default topic
    local queue_arn=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        sqs get-queue-attributes \
        --queue-url "http://localhost:4566/000000000000/default-queue" \
        --attribute-names QueueArn --query 'Attributes.QueueArn' --output text)

    local subscription_arn=$(aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        sns subscribe \
        --topic-arn "$topic_arn" \
        --protocol sqs \
        --notification-endpoint "$queue_arn" \
        --query 'SubscriptionArn' --output text)

    aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        sns set-subscription-attributes \
        --subscription-arn "$subscription_arn" \
        --attribute-name RawMessageDelivery \
        --attribute-value true

    log_success "Subscribed default-queue to default-topic with raw message delivery"
}

# List created resources
list_resources() {
    log_info "Listing created AWS resources..."

    echo ""
    log_info "📋 SQS Queues:"
    aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        sqs list-queues --query 'QueueUrls' --output table 2>/dev/null || log_warning "No SQS queues found"

    echo ""
    log_info "📋 SNS Topics:"
    aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        sns list-topics --query 'Topics[].TopicArn' --output table 2>/dev/null || log_warning "No SNS topics found"

    echo ""
    log_info "📋 SNS Subscriptions:"
    aws --endpoint-url="$AWS_ENDPOINT" --region="$AWS_REGION" \
        sns list-subscriptions --query 'Subscriptions[].{Protocol:Protocol,Endpoint:Endpoint,TopicArn:TopicArn}' --output table 2>/dev/null || log_warning "No SNS subscriptions found"

    echo ""
}

# Main execution
main() {
    log_info "🚀 Starting AWS resource initialization for LocalStack"

    # Wait for LocalStack to be ready
    if ! wait_for_localstack; then
        exit 1
    fi

    # Check if we have configuration
    if check_config; then
        log_info "Using configuration from: $CONFIG_FILE"

        # Create resources from configuration
        create_sqs_queues
        create_sns_topics
        create_dynamodb_tables
    else
        # Create default resources
        create_default_resources
    fi

    # List all created resources
    list_resources

    log_success "🎉 AWS resource initialization completed!"
}

# Run main function
main "$@"
