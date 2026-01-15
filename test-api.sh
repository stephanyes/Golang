#!/bin/bash

# Test GO API - Coin Balance Endpoint
# Server runs on localhost:8000

echo "=== Test 1: Missing Authorization Header ==="
curl -X GET "http://localhost:8000/account/coins?username=testuser" \
  -H "Content-Type: application/json" \
  -v

echo -e "\n\n=== Test 2: Missing Username Parameter ==="
curl -X GET "http://localhost:8000/account/coins" \
  -H "Authorization: test-token-123" \
  -H "Content-Type: application/json" \
  -v

echo -e "\n\n=== Test 3: Valid Request - User 'alex' (token: 123ABC, expected balance: 100) ==="
curl -X GET "http://localhost:8000/account/coins?username=alex" \
  -H "Authorization: 123ABC" \
  -H "Content-Type: application/json" \
  -v

echo -e "\n\n=== Test 4: Valid Request - User 'jason' (token: 456DEF, expected balance: 200) ==="
curl -X GET "http://localhost:8000/account/coins?username=jason" \
  -H "Authorization: 456DEF" \
  -H "Content-Type: application/json" \
  -v

echo -e "\n\n=== Test 5: Valid Request - User 'marie' (token: 789GHI, expected balance: 300) ==="
curl -X GET "http://localhost:8000/account/coins?username=marie" \
  -H "Authorization: 789GHI" \
  -H "Content-Type: application/json" \
  -v

echo -e "\n\n=== Test 6: Invalid Token (wrong token for username) ==="
curl -X GET "http://localhost:8000/account/coins?username=alex" \
  -H "Authorization: wrong-token" \
  -H "Content-Type: application/json" \
  -v

echo -e "\n\n=== Test 7: Pretty JSON Output (alex) ==="
curl -X GET "http://localhost:8000/account/coins?username=alex" \
  -H "Authorization: 123ABC" \
  -H "Content-Type: application/json" \
  -s | jq '.'

