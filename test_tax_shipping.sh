#!/bin/bash

# Simple test for tax and shipping calculation
echo "===== Tax and Shipping Test ====="

BASE_URL="http://localhost:8080/api/v1"

# Use predefined users (admin and user)
echo "Logging in with existing accounts..."

# Login as admin
ADMIN_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Admin token obtained"

# Login as user
USER_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123"}')
USER_TOKEN=$(echo $USER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "User token obtained"
echo ""

# Create products
echo "Creating test products..."
PROD1=$(curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"Product 2000yen","price":2000,"stock":100,"category":"Test"}')
PROD1_ID=$(echo $PROD1 | grep -o '"id":"[^"]*' | cut -d'"' -f4)

PROD2=$(curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"Product 3000yen","price":3000,"stock":50,"category":"Test"}')
PROD2_ID=$(echo $PROD2 | grep -o '"id":"[^"]*' | cut -d'"' -f4)

echo "Product 1 (2000 yen): $PROD1_ID"
echo "Product 2 (3000 yen): $PROD2_ID"
echo ""

# Test 1: Order under 5000 yen
echo "===== Test 1: Order Total 4000 yen (under 5000) ====="
echo "Order: Product1 x2 = 4000 yen"
echo "Expected: Tax=400, Shipping=500, Total=4900"
ORDER1=$(curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"items\":[{\"product_id\":\"$PROD1_ID\",\"quantity\":2}]}")

TOTAL1=$(echo $ORDER1 | grep -o '"total_price":[0-9]*' | cut -d':' -f2)
SHIPPING1=$(echo $ORDER1 | grep -o '"shipping_fee":[0-9]*' | cut -d':' -f2)
echo "Result: Total=$TOTAL1, Shipping=$SHIPPING1"
if [ "$TOTAL1" = "4900" ] && [ "$SHIPPING1" = "500" ]; then
  echo "✓ Test 1 PASSED"
else
  echo "✗ Test 1 FAILED"
fi
echo ""

# Test 2: Order exactly 5000 yen
echo "===== Test 2: Order Total 5000 yen (exactly 5000) ====="
echo "Order: Product1 x1 + Product2 x1 = 5000 yen"
echo "Expected: Tax=500, Shipping=0, Total=5500"
ORDER2=$(curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"items\":[{\"product_id\":\"$PROD1_ID\",\"quantity\":1},{\"product_id\":\"$PROD2_ID\",\"quantity\":1}]}")

TOTAL2=$(echo $ORDER2 | grep -o '"total_price":[0-9]*' | cut -d':' -f2)
SHIPPING2=$(echo $ORDER2 | grep -o '"shipping_fee":[0-9]*' | cut -d':' -f2)
echo "Result: Total=$TOTAL2, Shipping=$SHIPPING2"
if [ "$TOTAL2" = "5500" ] && [ "$SHIPPING2" = "0" ]; then
  echo "✓ Test 2 PASSED"
else
  echo "✗ Test 2 FAILED"
fi
echo ""

# Test 3: Order over 5000 yen
echo "===== Test 3: Order Total 7000 yen (over 5000) ====="
echo "Order: Product1 x2 + Product2 x1 = 7000 yen"
echo "Expected: Tax=700, Shipping=0, Total=7700"
ORDER3=$(curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"items\":[{\"product_id\":\"$PROD1_ID\",\"quantity\":2},{\"product_id\":\"$PROD2_ID\",\"quantity\":1}]}")

TOTAL3=$(echo $ORDER3 | grep -o '"total_price":[0-9]*' | cut -d':' -f2)
SHIPPING3=$(echo $ORDER3 | grep -o '"shipping_fee":[0-9]*' | cut -d':' -f2)
echo "Result: Total=$TOTAL3, Shipping=$SHIPPING3"
if [ "$TOTAL3" = "7700" ] && [ "$SHIPPING3" = "0" ]; then
  echo "✓ Test 3 PASSED"
else
  echo "✗ Test 3 FAILED"
fi
echo ""

echo "===== Test Complete ====="