#!/bin/bash

echo "=== Regression Test for Tax/Shipping with Multi-Warehouse ==="

# Get admin token
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

TOKEN=$(echo "$RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
echo "Admin token obtained"

# Test 1: Order under 5000 yen (should have 500 yen shipping)
echo -e "\nTest 1: Order with Laptop x2 (2400 yen + tax + shipping)..."
ORDER1=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": "PROD-1767505093-3883",
        "quantity": 2
      }
    ]
  }')

echo "$ORDER1" | python3 -c "
import sys, json
data = json.load(sys.stdin)
if 'error' not in data:
    print(f\"Total: {data.get('total_price', 0)}, Shipping: {data.get('shipping_fee', 0)}, Status: {data.get('status', 'unknown')}\")
else:
    print(f\"Error: {data['error']}\")
"

# Test 2: Order over 5000 yen (should have free shipping)
echo -e "\nTest 2: Order with Laptop x5 (6000 yen + tax, no shipping)..."
ORDER2=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": "PROD-1767505093-3883",
        "quantity": 5
      }
    ]
  }')

echo "$ORDER2" | python3 -c "
import sys, json
data = json.load(sys.stdin)
if 'error' not in data:
    print(f\"Total: {data.get('total_price', 0)}, Shipping: {data.get('shipping_fee', 0)}, Status: {data.get('status', 'unknown')}\")
else:
    print(f\"Error: {data['error']}\")
"

echo
echo "=== Regression Test Complete ==="
