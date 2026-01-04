#!/bin/bash

echo "=== Complete Multi-Warehouse Stock Management Test ==="

# Clean restart
echo "Restarting server..."
pkill -f ./server 2>/dev/null
sleep 1
./server &
sleep 2

# Login as admin
echo -e "\n1. Login as admin..."
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  python3 -c "import sys, json; print(json.load(sys.stdin).get('token', ''))")

if [ -z "$TOKEN" ]; then
  echo "Failed to get token. Test aborted."
  exit 1
fi

echo "Token obtained successfully"

# Create products using admin token
echo -e "\n2. Creating products..."
curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","price":1200,"category":"Electronics"}' | python3 -m json.tool

curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Mouse","price":25,"category":"Electronics"}' | python3 -m json.tool

# Get product list
echo -e "\n3. Fetching product list..."
curl -s http://localhost:8080/api/v1/products | python3 -m json.tool

# Try to create an order (should fail due to no stock)
echo -e "\n4. Attempting to create order (should fail - no stock)..."
curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": "PROD-1234567890",
        "quantity": 2
      }
    ]
  }' 2>/dev/null | python3 -m json.tool

echo -e "\n=== Test Complete ==="