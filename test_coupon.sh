#!/bin/bash

echo "=== Testing Coupon Functionality ==="

# Kill any existing server
pkill -f "./server" 2>/dev/null
sleep 1

# Start server in background
./server &
SERVER_PID=$!
sleep 2

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get product list and extract IDs using Python (more reliable JSON parsing)
echo -e "\n${YELLOW}Fetching product IDs...${NC}"
PRODUCTS_JSON=$(curl -s http://localhost:8080/api/v1/products)

# Extract product IDs dynamically using Python for JSON parsing
LAPTOP_ID=$(echo "$PRODUCTS_JSON" | python3 -c "import json, sys; data = json.load(sys.stdin); [print(p['id']) for p in data['products'] if p['name'] == 'Laptop']" 2>/dev/null | head -1)
MOUSE_ID=$(echo "$PRODUCTS_JSON" | python3 -c "import json, sys; data = json.load(sys.stdin); [print(p['id']) for p in data['products'] if p['name'] == 'Mouse']" 2>/dev/null | head -1)
KEYBOARD_ID=$(echo "$PRODUCTS_JSON" | python3 -c "import json, sys; data = json.load(sys.stdin); [print(p['id']) for p in data['products'] if p['name'] == 'Keyboard']" 2>/dev/null | head -1)

echo "Product IDs found:"
echo "  Laptop: $LAPTOP_ID"
echo "  Mouse: $MOUSE_ID"
echo "  Keyboard: $KEYBOARD_ID"

# Login as admin
echo -e "\n${YELLOW}1. Logging in as admin...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Token obtained: ${TOKEN:0:20}..."

# Test 1: Order with percentage coupon (SAVE10 - 10% off, min 1000 yen)
echo -e "\n${YELLOW}Test 1: Order with SAVE10 (10% off, min 1000 yen)${NC}"
echo "Creating order with Keyboard (75*2=150) + Mouse (25*2=50) = 200 yen (below minimum)"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"items\": [
      {\"product_id\": \"$KEYBOARD_ID\", \"quantity\": 2},
      {\"product_id\": \"$MOUSE_ID\", \"quantity\": 2}
    ],
    \"coupon_code\": \"SAVE10\"
  }")
echo "Response: $RESPONSE"
if [[ $RESPONSE == *"does not meet minimum requirement"* ]]; then
  echo -e "${GREEN}✓ Correctly rejected: Order below minimum${NC}"
else
  echo -e "${RED}✗ Failed: Should reject order below minimum${NC}"
fi

# Test 2: Valid order with SAVE10 (meets minimum)
echo -e "\n${YELLOW}Test 2: Valid order with SAVE10 (meets minimum)${NC}"
echo "Creating order with Laptop (1200*1) = 1200 yen (meets minimum of 1000 yen)"
echo "Expected: Subtotal: 1200, Tax: 120, Subtotal+Tax: 1320, Shipping: 500"
echo "Discount: 10% of 1320 = 132, Final: 1320 - 132 + 500 = 1688"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"items\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1}
    ],
    \"coupon_code\": \"SAVE10\"
  }")
echo "Response: $RESPONSE"
TOTAL=$(echo $RESPONSE | grep -o '"total_price":[0-9]*' | cut -d':' -f2)
DISCOUNT=$(echo $RESPONSE | grep -o '"discount_amount":[0-9]*' | cut -d':' -f2)
APPLIED_COUPON=$(echo $RESPONSE | grep -o '"applied_coupon":"[^"]*' | cut -d'"' -f4)
if [[ "$TOTAL" == "1688" ]] && [[ "$DISCOUNT" == "132" ]] && [[ "$APPLIED_COUPON" == "SAVE10" ]]; then
  echo -e "${GREEN}✓ Correct: Total=1688, Discount=132, Coupon=SAVE10${NC}"
else
  echo -e "${RED}✗ Failed: Expected Total=1688, got Total=$TOTAL, Discount=$DISCOUNT${NC}"
fi

# Test 3: Fixed amount coupon (FLAT1000 - 1000 yen off, min 3000 yen)
echo -e "\n${YELLOW}Test 3: Order with FLAT1000 (1000 yen off, min 3000 yen)${NC}"
echo "Creating order with Laptop (1200*3) = 3600 yen"
echo "Expected: Subtotal: 3600, Tax: 360, Subtotal+Tax: 3960, Shipping: 500"
echo "Discount: 1000, Final: 3960 - 1000 + 500 = 3460"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"items\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 3}
    ],
    \"coupon_code\": \"FLAT1000\"
  }")
echo "Response: $RESPONSE"
TOTAL=$(echo $RESPONSE | grep -o '"total_price":[0-9]*' | cut -d':' -f2)
DISCOUNT=$(echo $RESPONSE | grep -o '"discount_amount":[0-9]*' | cut -d':' -f2)
if [[ "$TOTAL" == "3460" ]] && [[ "$DISCOUNT" == "1000" ]]; then
  echo -e "${GREEN}✓ Correct: Total=3460, Discount=1000${NC}"
else
  echo -e "${RED}✗ Failed: Expected Total=3460, got Total=$TOTAL, Discount=$DISCOUNT${NC}"
fi

# Test 4: Order above 5000 yen (free shipping) with WELCOME20 (20% off)
echo -e "\n${YELLOW}Test 4: Large order with WELCOME20 (20% off, free shipping)${NC}"
echo "Creating order with Laptop (1200*5) = 6000 yen"
echo "Expected: Subtotal: 6000, Tax: 600, Subtotal+Tax: 6600, Shipping: 0 (free)"
echo "Discount: 20% of 6600 = 1320, Final: 6600 - 1320 + 0 = 5280"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"items\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 5}
    ],
    \"coupon_code\": \"WELCOME20\"
  }")
echo "Response: $RESPONSE"
TOTAL=$(echo $RESPONSE | grep -o '"total_price":[0-9]*' | cut -d':' -f2)
DISCOUNT=$(echo $RESPONSE | grep -o '"discount_amount":[0-9]*' | cut -d':' -f2)
SHIPPING=$(echo $RESPONSE | grep -o '"shipping_fee":[0-9]*' | cut -d':' -f2)
if [[ "$TOTAL" == "5280" ]] && [[ "$DISCOUNT" == "1320" ]] && [[ "$SHIPPING" == "0" ]]; then
  echo -e "${GREEN}✓ Correct: Total=5280, Discount=1320, Shipping=0${NC}"
else
  echo -e "${RED}✗ Failed: Expected Total=5280, got Total=$TOTAL, Discount=$DISCOUNT, Shipping=$SHIPPING${NC}"
fi

# Test 5: Invalid coupon code
echo -e "\n${YELLOW}Test 5: Invalid coupon code${NC}"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"items\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1}
    ],
    \"coupon_code\": \"INVALID999\"
  }")
echo "Response: $RESPONSE"
if [[ $RESPONSE == *"coupon not found"* ]] || [[ $RESPONSE == *"invalid coupon"* ]]; then
  echo -e "${GREEN}✓ Correctly rejected invalid coupon${NC}"
else
  echo -e "${RED}✗ Failed: Should reject invalid coupon${NC}"
fi

# Test 6: Order without coupon (verify it still works)
echo -e "\n${YELLOW}Test 6: Order without coupon${NC}"
echo "Creating order with Mouse (25*10) = 250 yen"
echo "Expected: Subtotal: 250, Tax: 25, Subtotal+Tax: 275, Shipping: 500, Total: 775"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"items\": [
      {\"product_id\": \"$MOUSE_ID\", \"quantity\": 10}
    ]
  }")
echo "Response: $RESPONSE"
TOTAL=$(echo $RESPONSE | grep -o '"total_price":[0-9]*' | cut -d':' -f2)
if [[ "$TOTAL" == "775" ]]; then
  echo -e "${GREEN}✓ Correct: Order without coupon works, Total=775${NC}"
else
  echo -e "${RED}✗ Failed: Expected Total=775, got Total=$TOTAL${NC}"
fi

# Test 7: Check coupon usage limit
echo -e "\n${YELLOW}Test 7: Coupon usage tracking${NC}"
echo "Using FLASH500 coupon (should work first time)"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"items\": [
      {\"product_id\": \"$LAPTOP_ID\", \"quantity\": 2}
    ],
    \"coupon_code\": \"FLASH500\"
  }")
echo "Response: $RESPONSE"
ORDER_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
if [[ -n "$ORDER_ID" ]]; then
  echo -e "${GREEN}✓ FLASH500 coupon applied successfully${NC}"

  # Show order details
  TOTAL=$(echo $RESPONSE | grep -o '"total_price":[0-9]*' | cut -d':' -f2)
  DISCOUNT=$(echo $RESPONSE | grep -o '"discount_amount":[0-9]*' | cut -d':' -f2)
  echo "  Order Total: $TOTAL yen (Discount: $DISCOUNT yen)"
else
  echo -e "${RED}✗ Failed: FLASH500 should work${NC}"
  echo "  Note: Expected (2*1200=2400, tax:240, shipping:500, -500 discount = 2640 yen total)"
fi

# Clean up
echo -e "\n${YELLOW}Cleaning up...${NC}"
kill $SERVER_PID 2>/dev/null

echo -e "\n${GREEN}=== Coupon Functionality Tests Completed ===${NC}"