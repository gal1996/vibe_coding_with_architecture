#!/bin/bash

echo "===== Wishlist & Recommendations Test Script ====="

# Kill any existing server
pkill -f "./server" 2>/dev/null
sleep 1

# Start the server
echo "Starting server..."
./server &
SERVER_PID=$!
sleep 2

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Test endpoints

# 1. Admin Login
echo -e "\n${YELLOW}1. Admin Login${NC}"
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$TOKEN" ]; then
    echo -e "${GREEN}✓ Login successful${NC}"
else
    echo -e "${RED}✗ Login failed${NC}"
    kill $SERVER_PID
    exit 1
fi

# 2. Get Products (with is_favorite field)
echo -e "\n${YELLOW}2. Get Products with is_favorite field${NC}"
PRODUCTS=$(curl -s -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN")

echo "$PRODUCTS" | grep -q "is_favorite"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Products include is_favorite field${NC}"
    echo "Sample: $(echo "$PRODUCTS" | python3 -c "import json, sys; data = json.load(sys.stdin); print(json.dumps(data['products'][0], indent=2) if data['products'] else 'No products')" 2>/dev/null | head -20)"
else
    echo -e "${RED}✗ is_favorite field not found${NC}"
fi

# Get first product ID for testing
PRODUCT_ID=$(echo "$PRODUCTS" | python3 -c "import json, sys; data = json.load(sys.stdin); print(data['products'][0]['id'] if data['products'] else '')" 2>/dev/null)

# 3. Add to Wishlist
echo -e "\n${YELLOW}3. Add Product to Wishlist${NC}"
ADD_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/wishlist/$PRODUCT_ID \
  -H "Authorization: Bearer $TOKEN")

echo "$ADD_RESPONSE" | grep -q "added to wishlist"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Product added to wishlist${NC}"
    echo "Response: $ADD_RESPONSE"
else
    echo -e "${RED}✗ Failed to add to wishlist${NC}"
    echo "Response: $ADD_RESPONSE"
fi

# 4. Get Wishlist
echo -e "\n${YELLOW}4. Get My Wishlist${NC}"
WISHLIST=$(curl -s -X GET http://localhost:8080/api/v1/wishlist \
  -H "Authorization: Bearer $TOKEN")

echo "$WISHLIST" | grep -q "wishlist"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Wishlist retrieved${NC}"
    echo "Wishlist: $WISHLIST"
else
    echo -e "${RED}✗ Failed to get wishlist${NC}"
fi

# 5. Check is_favorite is now true for added product
echo -e "\n${YELLOW}5. Check is_favorite is true for wishlist item${NC}"
PRODUCT_DETAIL=$(curl -s -X GET http://localhost:8080/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $TOKEN")

IS_FAVORITE=$(echo "$PRODUCT_DETAIL" | python3 -c "import json, sys; data = json.load(sys.stdin); print(data.get('is_favorite', False))" 2>/dev/null)

if [ "$IS_FAVORITE" = "True" ]; then
    echo -e "${GREEN}✓ Product is_favorite=true${NC}"
else
    echo -e "${RED}✗ Product is_favorite=false or missing${NC}"
fi

# 6. Get Recommendations
echo -e "\n${YELLOW}6. Get Recommendations${NC}"
RECOMMENDATIONS=$(curl -s -X GET http://localhost:8080/api/v1/users/me/recommendations \
  -H "Authorization: Bearer $TOKEN")

echo "$RECOMMENDATIONS" | grep -q "recommendations"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Recommendations retrieved${NC}"
    echo "Recommendations: $RECOMMENDATIONS"
else
    echo -e "${RED}✗ Failed to get recommendations${NC}"
    echo "Response: $RECOMMENDATIONS"
fi

# 7. Remove from Wishlist
echo -e "\n${YELLOW}7. Remove from Wishlist${NC}"
REMOVE_RESPONSE=$(curl -s -X DELETE http://localhost:8080/api/v1/wishlist/$PRODUCT_ID \
  -H "Authorization: Bearer $TOKEN")

echo "$REMOVE_RESPONSE" | grep -q "removed from wishlist"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Product removed from wishlist${NC}"
else
    echo -e "${RED}✗ Failed to remove from wishlist${NC}"
fi

# 8. Check is_favorite is now false
echo -e "\n${YELLOW}8. Check is_favorite is false after removal${NC}"
PRODUCT_DETAIL_AFTER=$(curl -s -X GET http://localhost:8080/api/v1/products/$PRODUCT_ID \
  -H "Authorization: Bearer $TOKEN")

IS_FAVORITE_AFTER=$(echo "$PRODUCT_DETAIL_AFTER" | python3 -c "import json, sys; data = json.load(sys.stdin); print(data.get('is_favorite', False))" 2>/dev/null)

if [ "$IS_FAVORITE_AFTER" = "False" ]; then
    echo -e "${GREEN}✓ Product is_favorite=false${NC}"
else
    echo -e "${RED}✗ Product is_favorite still true${NC}"
fi

# 9. Test unauthenticated access (should always show is_favorite=false)
echo -e "\n${YELLOW}9. Test Unauthenticated Access${NC}"
UNAUTH_PRODUCTS=$(curl -s -X GET http://localhost:8080/api/v1/products)

UNAUTH_FAVORITE=$(echo "$UNAUTH_PRODUCTS" | python3 -c "
import json, sys
data = json.load(sys.stdin)
if data.get('products'):
    print(data['products'][0].get('is_favorite', 'missing'))
else:
    print('no products')
" 2>/dev/null)

if [ "$UNAUTH_FAVORITE" = "False" ] || [ "$UNAUTH_FAVORITE" = "missing" ]; then
    echo -e "${GREEN}✓ Unauthenticated shows is_favorite=false${NC}"
else
    echo -e "${RED}✗ Unauthenticated shows is_favorite=true (incorrect)${NC}"
fi

# Clean up
kill $SERVER_PID 2>/dev/null

echo -e "\n${GREEN}===== Wishlist Test Complete =====${NC}"