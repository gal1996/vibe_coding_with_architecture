#!/bin/bash

echo "Simple Tax and Shipping Test"
echo "==========================="

# Login as user
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123"}')

USER_TOKEN=$(echo "$USER_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

echo "Login response: $USER_RESPONSE"
echo "User token: $USER_TOKEN"

# Test order: 2 products at 2000 yen = 4000 yen (should add 500 yen shipping)
echo ""
echo "Test: Order 4000 yen (2 x 2000 yen product)"
echo "Expected: Tax=400, Shipping=500, Total=4900"
echo ""
ORDER=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d '{"items":[{"product_id":"PROD-1234567890","quantity":2}]}')

echo "Order Result:"
echo $ORDER | python3 -m json.tool