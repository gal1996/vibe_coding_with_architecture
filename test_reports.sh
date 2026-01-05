#!/bin/bash

echo "===== 管理者レポートAPIテスト ====="

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

# 1. 管理者としてログイン
echo -e "\n${YELLOW}1. 管理者ログイン${NC}"
ADMIN_LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')
ADMIN_TOKEN=$(echo $ADMIN_LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "管理者トークン取得: ${ADMIN_TOKEN:0:20}..."

# 2. 通常ユーザーとしてログイン
echo -e "\n${YELLOW}2. 通常ユーザーログイン${NC}"
USER_LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"user123"}')
USER_TOKEN=$(echo $USER_LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [[ -z "$USER_TOKEN" ]]; then
  echo "通常ユーザーが存在しない場合は作成"
  curl -s -X POST http://localhost:8080/api/v1/register \
    -H "Content-Type: application/json" \
    -d '{"username":"user","password":"user123","is_admin":false}' > /dev/null

  USER_LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/login \
    -H "Content-Type: application/json" \
    -d '{"username":"user","password":"user123"}')
  USER_TOKEN=$(echo $USER_LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
fi
echo "ユーザートークン取得: ${USER_TOKEN:0:20}..."

# 3. 商品リストを取得してテスト用注文を作成
echo -e "\n${YELLOW}3. テスト用注文を作成${NC}"
PRODUCTS_JSON=$(curl -s http://localhost:8080/api/v1/products)
LAPTOP_ID=$(echo "$PRODUCTS_JSON" | python3 -c "import json, sys; data = json.load(sys.stdin); [print(p['id']) for p in data['products'] if p['name'] == 'Laptop']" 2>/dev/null | head -1)
MOUSE_ID=$(echo "$PRODUCTS_JSON" | python3 -c "import json, sys; data = json.load(sys.stdin); [print(p['id']) for p in data['products'] if p['name'] == 'Mouse']" 2>/dev/null | head -1)
KEYBOARD_ID=$(echo "$PRODUCTS_JSON" | python3 -c "import json, sys; data = json.load(sys.stdin); [print(p['id']) for p in data['products'] if p['name'] == 'Keyboard']" 2>/dev/null | head -1)

# いくつか注文を作成（管理者として）
echo "注文1: Laptop x2 with SAVE10"
ORDER1=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{
    \"items\": [{\"product_id\": \"$LAPTOP_ID\", \"quantity\": 2}],
    \"coupon_code\": \"SAVE10\"
  }")

echo "注文2: Mouse x10 without coupon"
ORDER2=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{
    \"items\": [{\"product_id\": \"$MOUSE_ID\", \"quantity\": 10}]
  }")

echo "注文3: Keyboard x3 with FLAT1000"
ORDER3=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "{
    \"items\": [{\"product_id\": \"$KEYBOARD_ID\", \"quantity\": 3}],
    \"coupon_code\": \"FLAT1000\"
  }")

# ユーザーとしても注文
echo "注文4: Laptop x1 with WELCOME20 (ユーザー)"
ORDER4=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{
    \"items\": [{\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1}],
    \"coupon_code\": \"WELCOME20\"
  }")

# 4. 通常ユーザーがレポートAPIにアクセス（拒否されるはず）
echo -e "\n${YELLOW}4. 通常ユーザーがレポートAPIにアクセス（拒否されるはず）${NC}"
USER_REPORT=$(curl -s -X GET http://localhost:8080/api/v1/admin/reports/sales \
  -H "Authorization: Bearer $USER_TOKEN")
echo "Response: $USER_REPORT"
if [[ $USER_REPORT == *"permission denied"* ]] || [[ $USER_REPORT == *"admin access required"* ]]; then
  echo -e "${GREEN}✓ 正しく拒否されました${NC}"
else
  echo -e "${RED}✗ 通常ユーザーがアクセスできてしまいました${NC}"
fi

# 5. 管理者がレポートAPIにアクセス
echo -e "\n${YELLOW}5. 管理者がレポートAPIにアクセス${NC}"
ADMIN_REPORT=$(curl -s -X GET http://localhost:8080/api/v1/admin/reports/sales \
  -H "Authorization: Bearer $ADMIN_TOKEN")
echo "Response: $ADMIN_REPORT"

# レポート内容を解析
if [[ $ADMIN_REPORT == *"sales_summary"* ]]; then
  echo -e "${GREEN}✓ 販売レポートが正常に取得できました${NC}"

  # Pretty print the report
  echo -e "\n${YELLOW}=== 販売レポート詳細 ===${NC}"
  echo "$ADMIN_REPORT" | python3 -m json.tool 2>/dev/null || echo "$ADMIN_REPORT"
else
  echo -e "${RED}✗ 販売レポートの取得に失敗しました${NC}"
fi

# 6. 認証なしでレポートAPIにアクセス（拒否されるはず）
echo -e "\n${YELLOW}6. 認証なしでレポートAPIにアクセス（拒否されるはず）${NC}"
NO_AUTH_REPORT=$(curl -s -X GET http://localhost:8080/api/v1/admin/reports/sales)
echo "Response: $NO_AUTH_REPORT"
if [[ $NO_AUTH_REPORT == *"Invalid token"* ]] || [[ $NO_AUTH_REPORT == *"authorization"* ]]; then
  echo -e "${GREEN}✓ 正しく拒否されました${NC}"
else
  echo -e "${RED}✗ 認証なしでアクセスできてしまいました${NC}"
fi

# Clean up
echo -e "\n${YELLOW}クリーンアップ...${NC}"
kill $SERVER_PID 2>/dev/null

echo -e "\n${GREEN}===== レポートAPIテスト完了 =====${NC}"