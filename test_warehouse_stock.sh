#!/bin/bash

# Create test script for warehouse stock functionality

# Login as admin
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  python3 -c "import sys, json; print(json.load(sys.stdin).get('token', ''))")

echo "Token obtained: ${TOKEN:0:20}..."

# Get product list to see warehouse stock info
echo -e "\n=== Product List with Warehouse Stock ==="
curl -s http://localhost:8080/api/v1/products/PROD-1234567890 | python3 -m json.tool

# Try to create an order to test stock allocation
echo -e "\n=== Create Order Test ==="
curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": "PROD-1234567890",
        "quantity": 1
      }
    ]
  }' | python3 -m json.tool || echo "Order creation failed (expected - no stock)"

echo -e "\n=== Test Complete ==="