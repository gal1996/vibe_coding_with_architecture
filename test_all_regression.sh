#!/bin/bash

echo "===== 総合リグレッションテスト ====="

# Kill any existing server
pkill -f "./server" 2>/dev/null
sleep 1

# Start server
./server &
SERVER_PID=$!
sleep 2

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS=0
FAIL=0

# Test function
test_result() {
    if [ "$1" = "true" ]; then
        echo -e "${GREEN}✓ $2${NC}"
        PASS=$((PASS + 1))
    else
        echo -e "${RED}✗ $2${NC}"
        FAIL=$((FAIL + 1))
    fi
}

# 1. ヘルスチェック
echo -e "\n${YELLOW}1. ヘルスチェック${NC}"
HEALTH=$(curl -s http://localhost:8080/health)
[[ $HEALTH == *"ok"* ]] && test_result "true" "ヘルスチェックOK" || test_result "false" "ヘルスチェック失敗"

# 2. 認証機能
echo -e "\n${YELLOW}2. 認証機能${NC}"
LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}')
TOKEN=$(echo $LOGIN | grep -o '"token":"[^"]*' | cut -d'"' -f4)
[[ -n "$TOKEN" ]] && test_result "true" "管理者ログイン成功" || test_result "false" "管理者ログイン失敗"

# 3. 製品管理（TASK_001）
echo -e "\n${YELLOW}3. 製品管理（TASK_001）${NC}"
PRODUCTS=$(curl -s http://localhost:8080/api/v1/products)
[[ $PRODUCTS == *"products"* ]] && test_result "true" "製品一覧取得" || test_result "false" "製品一覧取得失敗"

# 4. 税金・送料計算（TASK_002）
echo -e "\n${YELLOW}4. 税金・送料計算（TASK_002）${NC}"
LAPTOP_ID=$(echo "$PRODUCTS" | python3 -c "import json, sys; data = json.load(sys.stdin); [print(p['id']) for p in data['products'] if p['name'] == 'Laptop']" 2>/dev/null | head -1)
ORDER=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"items\": [{\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1}]}")
TOTAL=$(echo $ORDER | grep -o '"total_price":[0-9]*' | cut -d':' -f2)
SHIPPING=$(echo $ORDER | grep -o '"shipping_fee":[0-9]*' | cut -d':' -f2)
# 1200 + 税120 + 送料500 = 1820
[[ "$TOTAL" == "1820" ]] && test_result "true" "税金計算（10%）" || test_result "false" "税金計算エラー"
[[ "$SHIPPING" == "500" ]] && test_result "true" "送料計算（5000円未満）" || test_result "false" "送料計算エラー"

# 5. 決済機能（TASK_003）
echo -e "\n${YELLOW}5. 決済機能（TASK_003）${NC}"
STATUS=$(echo $ORDER | grep -o '"status":"[^"]*' | cut -d'"' -f4)
[[ "$STATUS" == "completed" ]] && test_result "true" "決済処理成功" || test_result "false" "決済処理失敗"

# 6. マルチ倉庫在庫管理（TASK_004）
echo -e "\n${YELLOW}6. マルチ倉庫在庫管理（TASK_004）${NC}"
PRODUCT_DETAIL=$(echo "$PRODUCTS" | python3 -c "import json, sys; data = json.load(sys.stdin); p = [p for p in data['products'] if p['name'] == 'Laptop'][0]; print(json.dumps(p))" 2>/dev/null)
[[ $PRODUCT_DETAIL == *"stocks"* ]] && test_result "true" "倉庫別在庫表示" || test_result "false" "倉庫別在庫表示エラー"
[[ $PRODUCT_DETAIL == *"total_stock"* ]] && test_result "true" "在庫合計表示" || test_result "false" "在庫合計表示エラー"

# 7. クーポン機能（TASK_005）
echo -e "\n${YELLOW}7. クーポン機能（TASK_005）${NC}"
MOUSE_ID=$(echo "$PRODUCTS" | python3 -c "import json, sys; data = json.load(sys.stdin); [print(p['id']) for p in data['products'] if p['name'] == 'Mouse']" 2>/dev/null | head -1)
COUPON_ORDER=$(curl -s -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"items\": [{\"product_id\": \"$LAPTOP_ID\", \"quantity\": 1}], \"coupon_code\": \"SAVE10\"}")
DISCOUNT=$(echo $COUPON_ORDER | grep -o '"discount_amount":[0-9]*' | cut -d':' -f2)
[[ "$DISCOUNT" == "132" ]] && test_result "true" "クーポン割引適用（10%）" || test_result "false" "クーポン割引エラー"

# 8. 管理者レポート（TASK_006）
echo -e "\n${YELLOW}8. 管理者レポート（TASK_006）${NC}"
REPORT=$(curl -s -X GET http://localhost:8080/api/v1/admin/reports/sales -H "Authorization: Bearer $TOKEN")
[[ $REPORT == *"sales_summary"* ]] && test_result "true" "販売サマリー取得" || test_result "false" "販売サマリー取得エラー"
[[ $REPORT == *"top_products"* ]] && test_result "true" "人気商品ランキング" || test_result "false" "人気商品ランキングエラー"
[[ $REPORT == *"warehouse_stock"* ]] && test_result "true" "倉庫別在庫サマリー" || test_result "false" "倉庫別在庫サマリーエラー"
[[ $REPORT == *"coupon_analytics"* ]] && test_result "true" "クーポン分析" || test_result "false" "クーポン分析エラー"

# 9. 認可チェック
echo -e "\n${YELLOW}9. 認可チェック${NC}"
UNAUTH=$(curl -s -X GET http://localhost:8080/api/v1/admin/reports/sales)
[[ $UNAUTH == *"Authorization"* ]] && test_result "true" "未認証アクセス拒否" || test_result "false" "未認証アクセス許可されてしまった"

# Clean up
kill $SERVER_PID 2>/dev/null

# Summary
echo -e "\n${YELLOW}===== テスト結果サマリー =====${NC}"
echo -e "成功: ${GREEN}$PASS${NC}"
echo -e "失敗: ${RED}$FAIL${NC}"

if [ $FAIL -eq 0 ]; then
    echo -e "\n${GREEN}すべてのテストが成功しました！✨${NC}"
    exit 0
else
    echo -e "\n${RED}$FAIL 件のテストが失敗しました${NC}"
    exit 1
fi