#!/bin/bash

# EC Site Backend API Test Script with Tax and Shipping Logic
# This script tests the API endpoints and verifies tax and shipping calculations

BASE_URL="http://localhost:8080/api/v1"
ADMIN_TOKEN=""
USER_TOKEN=""

echo "===== EC Site Backend API Test - Tax & Shipping Logic ====="
echo ""

# Health check
echo "1. Health Check"
curl -s ${BASE_URL%/api/v1}/health | jq '.'
echo ""

# Register admin
echo "2. Register Admin"
curl -s -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testadmin","password":"admin123","is_admin":true}' | jq '.'
echo ""

# Register user
echo "3. Register User"
curl -s -X POST $BASE_URL/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123","is_admin":false}' | jq '.'
echo ""

# Login as admin
echo "4. Login as Admin"
ADMIN_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')
echo $ADMIN_RESPONSE | jq '.'
ADMIN_TOKEN=$(echo $ADMIN_RESPONSE | jq -r '.token')
echo "Admin Token: ${ADMIN_TOKEN:0:20}..."
echo ""

# Login as user
echo "5. Login as User"
USER_RESPONSE=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123"}')
echo $USER_RESPONSE | jq '.'
USER_TOKEN=$(echo $USER_RESPONSE | jq -r '.token')
echo "User Token: ${USER_TOKEN:0:20}..."
echo ""

# List products (public)
echo "6. List Products (Public)"
curl -s -X GET $BASE_URL/products | jq '.'
echo ""

# List products with category filter
echo "7. List Products - Electronics Category"
curl -s -X GET "$BASE_URL/products?category=Electronics" | jq '.'
echo ""

# Create product (admin only)
echo "8. Create Product (Admin Only)"
curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"New Laptop","price":1500,"stock":5,"category":"Electronics"}' | jq '.'
echo ""

# Try to create product as regular user (should fail)
echo "9. Create Product as User (Should Fail)"
curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d '{"name":"Unauthorized Product","price":100,"stock":10,"category":"Test"}' | jq '.'
echo ""

# Create sample products for tax/shipping testing
echo "10. Creating Test Products for Tax & Shipping Tests"

# Product 1: 2000 yen
PROD1_RESPONSE=$(curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"Product A - 2000yen","price":2000,"stock":100,"category":"Test"}')
PROD1_ID=$(echo $PROD1_RESPONSE | jq -r '.id')
echo "Created Product A (2000 yen): ID=$PROD1_ID"

# Product 2: 1500 yen
PROD2_RESPONSE=$(curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"Product B - 1500yen","price":1500,"stock":50,"category":"Test"}')
PROD2_ID=$(echo $PROD2_RESPONSE | jq -r '.id')
echo "Created Product B (1500 yen): ID=$PROD2_ID"

# Product 3: 3000 yen
PROD3_RESPONSE=$(curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"name":"Product C - 3000yen","price":3000,"stock":30,"category":"Test"}')
PROD3_ID=$(echo $PROD3_RESPONSE | jq -r '.id')
echo "Created Product C (3000 yen): ID=$PROD3_ID"
echo ""

# Test Case 1: Order under 5000 yen (should add 500 yen shipping)
echo "11. Test Case 1: Order Total 4000 yen (Tax: 400, Shipping: 500)"
echo "   Order: Product A x2 = 4000 yen"
echo "   Expected: Subtotal: 4000, Tax: 400, Shipping: 500, Total: 4900"
ORDER1_RESPONSE=$(curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"items\":[{\"product_id\":\"$PROD1_ID\",\"quantity\":2}]}")
echo $ORDER1_RESPONSE | jq '.'
echo ""

# Test Case 2: Order exactly 5000 yen (free shipping)
echo "12. Test Case 2: Order Total 5000 yen (Tax: 500, Shipping: 0)"
echo "   Order: Product A x1 + Product C x1 = 5000 yen"
echo "   Expected: Subtotal: 5000, Tax: 500, Shipping: 0, Total: 5500"
ORDER2_RESPONSE=$(curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"items\":[{\"product_id\":\"$PROD1_ID\",\"quantity\":1},{\"product_id\":\"$PROD3_ID\",\"quantity\":1}]}")
echo $ORDER2_RESPONSE | jq '.'
echo ""

# Test Case 3: Order over 5000 yen (free shipping)
echo "13. Test Case 3: Order Total 8500 yen (Tax: 850, Shipping: 0)"
echo "   Order: Product A x2 + Product B x1 + Product C x1 = 8500 yen"
echo "   Expected: Subtotal: 8500, Tax: 850, Shipping: 0, Total: 9350"
ORDER3_RESPONSE=$(curl -s -X POST $BASE_URL/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"items\":[{\"product_id\":\"$PROD1_ID\",\"quantity\":2},{\"product_id\":\"$PROD2_ID\",\"quantity\":1},{\"product_id\":\"$PROD3_ID\",\"quantity\":1}]}")
echo $ORDER3_RESPONSE | jq '.'
echo ""

# List user orders to verify all created orders
echo "14. List All User Orders"
curl -s -X GET $BASE_URL/orders \
  -H "Authorization: Bearer $USER_TOKEN" | jq '.'
echo ""

echo "===== Tax & Shipping Logic Test Complete ====="
echo ""
echo "Summary:"
echo "- Orders under 5000 yen: 500 yen shipping fee added"
echo "- Orders 5000 yen or more: Free shipping (0 yen)"
echo "- All orders: 10% consumption tax applied"