#!/bin/bash

TOKEN_FILE=".user_token"

BASE_URL="http://localhost:8080/api/v1"

function save_token() {
    echo "$1" > $TOKEN_FILE
    echo "🔐 Token saved!"
}

function get_token() {
    if [ ! -f $TOKEN_FILE ]; then
        echo "❌ No token found. Please login first."
        exit 1
    fi
    cat $TOKEN_FILE
}

function register_user() {
    curl -s -X POST "$BASE_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\":\"$1\",
            \"username\":\"$2\",
            \"password\":\"$3\",
            \"first_name\":\"$4\",
            \"last_name\":\"$5\",
            \"phone\":\"$6\"
        }"
    echo
}

function login_user() {
    RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"email_or_username\":\"$1\",
            \"password\":\"$2\"
        }")

    echo "🔎 Login Response:"
    echo "$RESPONSE"

    TOKEN=$(echo "$RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d '"' -f4)

    if [ -n "$TOKEN" ]; then
        save_token "$TOKEN"
        echo "✅ Login successful!"
    else
        echo "❌ Login failed"
    fi
}

function get_me() {
    TOKEN=$(get_token)
    curl -s -X GET "$BASE_URL/users/me" \
        -H "Authorization: Bearer $TOKEN"
    echo
}

function update_profile() {
    TOKEN=$(get_token)
    curl -s -X PUT "$BASE_URL/users/me" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"first_name\":\"$1\",
            \"last_name\":\"$2\",
            \"phone\":\"$3\"
        }"
    echo
}

case "$1" in
    register)
        register_user "$2" "$3" "$4" "$5" "$6" "$7"
        ;;
    login)
        login_user "$2" "$3"
        ;;
    me)
        get_me
        ;;
    update)
        update_profile "$2" "$3" "$4"
        ;;
    *)
        echo "Available commands:"
        echo "./user_cli.sh register email username password first last phone"
        echo "./user_cli.sh login username password"
        echo "./user_cli.sh me"
        echo "./user_cli.sh update first last phone"
        ;;
esac
