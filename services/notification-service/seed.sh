#!/bin/sh

# Author: Vu Ngoc Thuan
# Description: This is is used to seed data for Novu System via API

# ===== VARIABLES & INFO =====

EMAIL="admin@gopher.com"
FIRST_NAME="GopherPay"
LAST_NAME="Admin"
PASSWORD="Admin@123"
#"Admin@123"
TOKEN=""

ORGANIZATION_NAME="GopherPay"
ENV_ID=""

# ====- API URL =====

BASE_URL_V1="http://localhost:3000/v1"
BASE_URL_V2="http://localhost:3000/v2"
AUTH_REGISTER_API="auth/register"
AUTH_LOGIN_API="auth/login"
ORGANIZATION_API="organizations"
ENVIRONMENT_API="environments"

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

    echo "2. Retrieve ENV_ID:"
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

# ===== SCRIPT =====
seed_admin
log_in
seed_org
get_development_env