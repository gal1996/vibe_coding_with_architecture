#!/bin/bash

echo "=== Payment Gateway Test ==="
echo

# Create some products first
echo "Setting up test data..."
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Create products if they don't exist
curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Product","price":1000,"stock":100,"category":"Test"}' > /dev/null 2>&1

PRODUCT_ID=$(curl -s http://localhost:8080/api/v1/products | python3 -c "import sys, json; data = json.load(sys.stdin); print(data['products'][0]['id'] if data.get('products') else 'PROD-1234567890')" 2>/dev/null || echo "PROD-1234567890")

echo "Using product ID: $PRODUCT_ID"
echo

# Test payment success rate
echo "Testing 10 orders for payment success rate (90% expected):"
echo

SUCCESS=0
FAILED=0

for i in {1..10}; do
  ORDER=$(curl -s -X POST http://localhost:8080/api/v1/orders \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"items\": [{\"product_id\": \"$PRODUCT_ID\", \"quantity\": 1}]}")

  if echo "$ORDER" | grep -q '"status":"completed"'; then
    SUCCESS=$((SUCCESS+1))
    echo "Attempt $i: ✓ Payment SUCCESS"
  elif echo "$ORDER" | grep -q 'payment_failed\|payment declined'; then
    FAILED=$((FAILED+1))
    echo "Attempt $i: ✗ Payment FAILED"
  else
    echo "Attempt $i: ? ERROR - Check response"
    echo "$ORDER" | head -c 100
    echo
  fi
done

echo
echo "=== Test Results ==="
echo "Successful payments: $SUCCESS/10"
echo "Failed payments: $FAILED/10"
echo "Success rate: $((SUCCESS * 10))%"

# Expected rate is 90%, so we expect around 9 successes
if [ $SUCCESS -ge 7 ] && [ $SUCCESS -le 10 ]; then
  echo "✓ Payment success rate is within expected range (70-100%)"
else
  echo "✗ Payment success rate is outside expected range"
fi

echo
echo "=== Payment Gateway Test Complete ===">