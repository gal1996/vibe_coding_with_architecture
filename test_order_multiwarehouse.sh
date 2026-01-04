#!/bin/bash

echo "=== Multi-Warehouse Order Test ==="
echo

# Login as admin
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  grep -o '"token":"[^"]*"' | cut -d'"' -f4)

echo "Token obtained: ${TOKEN:0:20}..."

# Create order that will allocate from multiple warehouses
echo -e "\nCreating order for Laptop x2 + Mouse x3..."
ORDER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": "PROD-1767504522-3569",
        "quantity": 2
      },
      {
        "product_id": "PROD-1767504522-175",
        "quantity": 3
      }
    ]
  }')

echo "$ORDER_RESPONSE" | python3 -m json.tool

# Check if payment successful (90% success rate)
if echo "$ORDER_RESPONSE" | grep -q "PaymentFailed"; then
  echo -e "\nPayment failed - order not completed"
else
  echo -e "\nPayment successful - stock should be reduced from warehouses"

  # Show updated stock levels
  echo -e "\n=== Updated Stock Levels ==="
  curl -s "http://localhost:8080/api/v1/products/PROD-1767504522-3569" | python3 -c "
import sys, json
data = json.load(sys.stdin)
print(f\"Laptop: Total Stock = {data.get('total_stock', 0)}\")
for stock in data.get('stocks', []):
    print(f\"  {stock['warehouse_name']}: {stock['quantity']}\")
"

  curl -s "http://localhost:8080/api/v1/products/PROD-1767504522-175" | python3 -c "
import sys, json
data = json.load(sys.stdin)
print(f\"Mouse: Total Stock = {data.get('total_stock', 0)}\")
for stock in data.get('stocks', []):
    print(f\"  {stock['warehouse_name']}: {stock['quantity']}\")
"
fi

echo -e "\n=== Test Complete ==="