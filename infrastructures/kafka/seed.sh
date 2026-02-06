#!/bin/sh

# Author: Vu Ngoc Thuan
# Description: This is is used to seed data for Novu System via API

load_env() {
    if [ -f .seed.env ]; then
        ENV_PATH=".seed.env"
    elif [ -f /tmp/.seed.env ]; then
        ENV_PATH="/tmp/.seed.env"
    fi

    if [ -n "$ENV_PATH" ]; then
        echo "===== Loading variables from $ENV_PATH"
        export $(grep -v '^#' "$ENV_PATH" | xargs)
    else
        echo "===== Error: .seed.env file not found!"
        exit 1
    fi
}

load_env

# ===== VARIABLES & INFO =====

API_KEY=$API_KEY

CLUSTER_NAME="Fintech_Event_Bus"
SAMPLE_TOPIC="email-verification"
SAMPLE_TOPIC_PARTITIONS="10"
CONFLUENT_COMMAND_TOPIC_PARTITIONS="5"
CONNECTOR_NAME="Novu-Listerner"

# ====- API URL =====

TOPIC_URL="http://kafka_dashboard:8080/api/clusters/$CLUSTER_NAME/topics"
CONNECTOR_URL="http://kafka_connect:8083/connectors"
API_URL="http://notification-service-api-1:3000/v1/events/trigger"

# ===== SEED FUNCTIONS =====

seed_sample_topic() {
    echo "===== Start seed Sample Topic process!"
    
    echo "1. Calling Create Topics API..."
    RESPONSE=$(curl -s -X POST "$TOPIC_URL" \
        -H 'Content-Type: application/json' \
        -d "{
            \"name\": \"$SAMPLE_TOPIC\", 
            \"partitions\": $SAMPLE_TOPIC_PARTITIONS, 
            \"configs\": {
                \"cleanup.policy\": \"delete\"
            }
        }"
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'
    
    echo "----- End seed seed Sample Topic process!"
}

seed_confluent_command_topic() {
    echo "===== Start seed _confluent-command Topic process!"
    
    echo "1. Calling Create Topics API..."
    RESPONSE=$(curl -s -X POST "$TOPIC_URL" \
        -H 'Content-Type: application/json' \
        -d "{
            \"name\": \"_confluent-command\", 
            \"partitions\": $CONFLUENT_COMMAND_TOPIC_PARTITIONS, 
            \"configs\": {
                \"cleanup.policy\": \"delete\"
            }
        }"
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'
    
    echo "----- End seed seed _confluent-command Topic process!"
}

seed_http_sink_connector() {
    echo "===== Start seed HTTP Sink Connector process!"
    
    echo "1. Calling Create Connector API..."
    RESPONSE=$(curl -s -X POST "$CONNECTOR_URL" \
        -H 'Content-Type: application/json' \
        -d "{
            \"name\": \"$CONNECTOR_NAME\",
            \"config\": {
                \"connector.class\": \"io.confluent.connect.http.HttpSinkConnector\",
                \"http.api.url\": \"$API_URL\",
                \"headers\": \"Content-Type:application/json|Authorization:ApiKey $API_KEY\",
                \"topics\": \"$SAMPLE_TOPIC\",
                \"tasks.max\": \"1\",
                \"batch.max.size\": \"1\",
                \"batch.json.as.array\": \"false\",
                \"concurrency.limit\": \"5\",
                \"request.method\": \"POST\",
                \"request.body.format\": \"json\",
                \"value.converter\": \"org.apache.kafka.connect.json.JsonConverter\",
                \"value.converter.schemas.enable\": \"false\",
                \"key.converter\": \"org.apache.kafka.connect.storage.StringConverter\",
                \"confluent.topic.bootstrap.servers\": \"kafka:29092\",
                \"reporter.bootstrap.servers\": \"kafka:29092\",
                \"reporter.error.topic.name\": \"http-error-topic\",
                \"reporter.result.topic.name\": \"http-success-topic\",
                \"reporter.error.topic.replication.factor\": \"1\",
                \"reporter.result.topic.replication.factor\": \"1\",
                \"errors.tolerance\": \"all\",
                \"errors.log.enable\": \"true\",
                \"errors.log.include.messages\": \"true\",
                \"confluent.license\": \"\"
            }
        }"
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'
    
    echo "----- End seed HTTP Sink Connector process!"
}

# ===== SCRIPT =====
echo "***** START SEED SCRIPT *****"
seed_sample_topic
sleep 1
seed_confluent_command_topic
sleep 1 
seed_http_sink_connector