#!/bin/bash

# Test script for payment gateway integration

echo "=== Payment Gateway Integration Test ==="
echo

# Login as admin
echo "1. Logging in as admin..."
ADMIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Admin logged in successfully"
echo

# Login as regular user
echo "2. Logging in as regular user..."
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123"}')
USER_TOKEN=$(echo $USER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "User logged in successfully"
echo

# Get products
echo "3. Getting product list..."
PRODUCTS=$(curl -s http://localhost:8080/api/v1/products)
echo "Products available:"
echo $PRODUCTS | python3 -m json.tool 2>/dev/null || echo $PRODUCTS
echo

# Try to create 10 orders to test 90% success rate
echo "4. Testing payment success rate (90% expected)..."
SUCCESS_COUNT=0
FAIL_COUNT=0

for i in {1..10}; do
  echo -n "Attempt $i: "

  ORDER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
    -H "Authorization: Bearer $USER_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
      "items": [
        {"product_id": "PRD-1", "quantity": 1},
        {"product_id": "PRD-2", "quantity": 2}
      ]
    }')

  if echo $ORDER_RESPONSE | grep -q "payment declined\|payment_failed"; then
    echo "Payment FAILED"
    FAIL_COUNT=$((FAIL_COUNT + 1))
    echo "Response: $ORDER_RESPONSE"
  elif echo $ORDER_RESPONSE | grep -q '"status":"completed"'; then
    echo "Payment SUCCESS"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
  else
    echo "Unknown response: $ORDER_RESPONSE"
  fi

  # Small delay between attempts
  sleep 0.5
done

echo
echo "=== Test Results ==="
echo "Successful payments: $SUCCESS_COUNT/10"
echo "Failed payments: $FAIL_COUNT/10"
echo "Success rate: $((SUCCESS_COUNT * 10))%"
echo

# Get user's orders
echo "5. Getting user's orders..."
USER_ORDERS=$(curl -s -X GET http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $USER_TOKEN")
echo "User's orders:"
echo $USER_ORDERS | python3 -m json.tool 2>/dev/null || echo $USER_ORDERS

echo
echo "=== Payment Gateway Integration Test Complete ==="