#!/bin/bash

echo "=== Mini-Ledger CockroachDB API Test ==="
echo "Make sure to start the services first:"
echo "docker-compose up -d"
echo
echo "API endpoint: http://localhost:8081"
echo

# Test 1: Get account balance
echo "1. Getting account balance:"
curl -s http://localhost:8081/api/v1/accounts/1/balance | jq .
echo

# Test 2: Get account holdings
echo "2. Getting account holdings:"
curl -s http://localhost:8081/api/v1/accounts/1/holdings | jq .
echo

# Test 3: Create a buy order
echo "3. Creating buy order (10 shares of STOCK02 at 50,000 each):"
BUY_ORDER=$(curl -s -X POST http://localhost:8081/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "stock_code": "STOCK02", "type": "LIMIT", "direction": "BUY", "quantity": 10, "price": 50000}')
echo $BUY_ORDER | jq .
BUY_ORDER_ID=$(echo $BUY_ORDER | jq -r .id)
echo

# Test 4: Check balance after buy order
echo "4. Account balance after buy order (should be 500,000):"
curl -s http://localhost:8081/api/v1/accounts/1/balance | jq .
echo

# Test 5: Create a sell order
echo "5. Creating sell order (50 shares of STOCK01 at 60,000 each):"
SELL_ORDER=$(curl -s -X POST http://localhost:8081/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "stock_code": "STOCK01", "type": "LIMIT", "direction": "SELL", "quantity": 50, "price": 60000}')
echo $SELL_ORDER | jq .
SELL_ORDER_ID=$(echo $SELL_ORDER | jq -r .id)
echo

# Test 6: Check holdings after sell order
echo "6. Account holdings after sell order (should be 50 STOCK01):"
curl -s http://localhost:8081/api/v1/accounts/1/holdings | jq .
echo

# Test 7: Cancel buy order
echo "7. Canceling buy order (should restore 500,000 to balance):"
curl -s -X DELETE http://localhost:8081/api/v1/orders/$BUY_ORDER_ID | jq .
echo

# Test 8: Check balance after cancel
echo "8. Account balance after canceling buy order (should be 1,000,000):"
curl -s http://localhost:8081/api/v1/accounts/1/balance | jq .
echo

# Test 9: Cancel sell order
echo "9. Canceling sell order (should restore 50 shares to holdings):"
curl -s -X DELETE http://localhost:8081/api/v1/orders/$SELL_ORDER_ID | jq .
echo

# Test 10: Check final holdings
echo "10. Final account holdings (should be 100 STOCK01):"
curl -s http://localhost:8081/api/v1/accounts/1/holdings | jq .
echo

# Test 11: Error case - insufficient funds
echo "11. Testing insufficient funds error:"
curl -s -X POST http://localhost:8081/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "stock_code": "STOCK03", "type": "LIMIT", "direction": "BUY", "quantity": 100, "price": 20000}' | jq .
echo

# Test 12: Error case - insufficient holdings
echo "12. Testing insufficient holdings error:"
curl -s -X POST http://localhost:8081/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"account_id": 1, "stock_code": "STOCK01", "type": "LIMIT", "direction": "SELL", "quantity": 200, "price": 50000}' | jq .
echo

echo "=== Test completed! ==="