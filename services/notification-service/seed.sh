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

EMAIL="admin@gopher.com"
FIRST_NAME="GopherPay"
LAST_NAME="Admin"
PASSWORD="Admin@123"
TOKEN=""

ORGANIZATION_NAME="GopherPay"
ENV_ID=""

SMTP_HOST="smtp.gmail.com"
SMTP_PORT="587"
SMTP_USER=$SMTP_USER
SMTP_PASSWORD=$SMTP_PASS

NOTIFICATION_GROUP_ID=""

# ====- API URL =====

BASE_URL_V1="http://localhost:3000/v1"
BASE_URL_V2="http://localhost:3000/v2"
AUTH_REGISTER_API="auth/register"
AUTH_LOGIN_API="auth/login"
ORGANIZATION_API="organizations"
ENVIRONMENT_API="environments"
INTEGRATION_API="integrations"
NOTIFICATION_GROUP_API="notification-groups"
WORKFLOW_API="workflows"

# ===== UTIL FUNCTIONS =====

log_in() {
    echo "===== Start login process!"
    
    echo "1. Calling $AUTH_LOGIN_API API..."
    RESPONSE=$(curl -s -X POST "$BASE_URL_V1/$AUTH_LOGIN_API" \
        -H 'Content-Type: application/json' \
        -d "{
            \"email\": \"$EMAIL\",
            \"password\": \"$PASSWORD\"
        }"
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'

    echo "3. Retrieve token:"
    export TOKEN=$(echo "$RESPONSE" | jq -r '.data.token')
    if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
        echo "Success - Token captured globally - $TOKEN"
    else
        echo "Error - Login failed, token is empty."
        return 1
    fi
    
    echo "----- End login process!"
}

get_development_env() {
    echo "===== Start getting Development Environment process!"
    
    echo "1. Calling $ENVIRONMENT_API API..."
    RESPONSE=$(curl -s -X GET "$BASE_URL_V1/$ENVIRONMENT_API" \
        -H "Authorization: Bearer $TOKEN" \
        -H 'Content-Type: application/json'
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'

    echo "3. Retrieve ENV_ID:"
    ENV_ID=$(echo "$RESPONSE" | jq -r '.data[] | select(.name == "Development") | ._id')

    if [ "$ENV_ID" != "null" ] && [ -n "$ENV_ID" ]; then
        echo "Success - Found Development Environment ID: $ENV_ID"
        export DEV_ENV_ID=$ENV_ID
    else
        echo "Error - Could not find environment 'Development'."
        return 1
    fi
    
    echo "----- End get environment process!"
}

get_notification_group_id() {
    echo "===== Start getting Notification Group Id process!"
    
    echo "1. Calling $NOTIFICATION_GROUP_API API..."
    RESPONSE=$(curl -s -X GET "$BASE_URL_V1/$NOTIFICATION_GROUP_API" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Novu-Environment-Id: $ENV_ID" \
        -H 'Content-Type: application/json'
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'

    echo "3. Retrieve Notification Group ID (First element):"
    GROUP_ID=$(echo "$RESPONSE" | jq -r '.data[0]._id')

    if [ "$GROUP_ID" != "null" ] && [ -n "$GROUP_ID" ]; then
        echo "Success - Found Group ID: $GROUP_ID"
        export NOTIFICATION_GROUP_ID=$GROUP_ID
    else
        echo "Error - Could not find any Notification Group."
        return 1
    fi
    
    echo "----- End get Notification Group Id process!"
}

# ===== SEED FUNCTIONS =====

seed_admin() {
    echo "===== Start seed Admin account process!"
    
    echo "1. Calling $AUTH_REGISTER_API API..."
    RESPONSE=$(curl -s -X POST "$BASE_URL_V1/$AUTH_REGISTER_API" \
        -H 'Content-Type: application/json' \
        -d "{
            \"email\": \"$EMAIL\",
            \"password\": \"$PASSWORD\",
            \"firstName\": \"$FIRST_NAME\",
            \"lastName\": \"$LAST_NAME\"
        }"
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'
    
    echo "----- End seed account Admin process!"
}

seed_org() {
    echo "===== Start seed Organization process!"
    
    echo "1. Calling $ORGANIZATION_API API..."
    RESPONSE=$(curl -s -X POST "$BASE_URL_V1/$ORGANIZATION_API" \
        -H "Authorization: Bearer $TOKEN" \
        -H 'Content-Type: application/json' \
        -d "{
            \"name\": \"$ORGANIZATION_NAME\"
        }"
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'
    
    echo "----- End seed Organization process!"
}

seed_integration() {
    echo "===== Start seed Email Integration process!"
    
    echo "1. Calling $INTEGRATION_API API..."
    RESPONSE=$(curl -s -X POST "$BASE_URL_V1/$INTEGRATION_API" \
        -H "Authorization: Bearer $TOKEN" \
        -H 'Content-Type: application/json' \
        -d "{
            \"providerId\": \"nodemailer\",
            \"channel\": \"email\",
            \"name\": \"Gmailer\",
            \"_environmentId\": \"$ENV_ID\",
            \"credentials\": {
                \"host\": \"$SMTP_HOST\",
                \"port\": \"$SMTP_PORT\",
                \"secure\": false,
                \"password\": \"$SMTP_PASSWORD\",
                \"from\": \"$SMTP_USER\",
                \"user\": \"$SMTP_USER\"
            },
            \"active\": true,
            \"check\": false
        }"
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'

    echo "----- End seed Email Integration process!"
}

seed_workflow() {
    echo "===== Start seed Workflow process!"
    
    echo "1. Calling $WORKFLOW_API API..."
    RESPONSE=$(curl -s -X POST "$BASE_URL_V1/$WORKFLOW_API" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Novu-Environment-Id: $ENV_ID" \
        -H 'Content-Type: application/json' \
        -d "{
            \"name\": \"Email Verification\",
            \"description\": \"Email Verification\",
            \"active\": true,
            \"notificationGroupId\": \"$NOTIFICATION_GROUP_ID\",
            \"steps\": [
                {
                    \"name\": \"Send verification email\",
                    \"template\": {
                        \"type\": \"email\",
                        \"subject\": \"Confirm Email\",
                        \"content\": \"Hi {{subscriber.firstName}}, click here: {{verificationLink}}\",
                        \"contentType\": \"html\",
                        \"layoutIdentifier\": \"default\"
                    }
                }
            ]
        }"
    )

    echo "2. Receive response:"
    echo "$RESPONSE" | jq '.'
    
    echo "----- End seed Workflow process!"
}

# ===== SCRIPT =====
echo "***** START SEED SCRIPT *****"
seed_admin
sleep 1
log_in
sleep 1
seed_org
sleep 1
log_in
sleep 1
get_development_env
sleep 1
seed_integration
sleep 1
get_notification_group_id
sleep 1
seed_workflow
