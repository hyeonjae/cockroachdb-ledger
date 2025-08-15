# Test 1: Get account balance
GET http://localhost:8081/api/v1/accounts/1/balance
HTTP 200
[Asserts]
jsonpath "$.account_number" == "AC001"
jsonpath "$.balance" == 1000000

# Test 2: Get account holdings
GET http://localhost:8081/api/v1/accounts/1/holdings
HTTP 200
[Asserts]
jsonpath "$" count == 1
jsonpath "$[0].stock_code" == "STOCK01"
jsonpath "$[0].quantity" == 100

# Test 3: Create buy order
POST http://localhost:8081/api/v1/orders
Content-Type: application/json
{
    "account_id": 1,
    "stock_code": "STOCK02",
    "type": "LIMIT",
    "direction": "BUY",
    "quantity": 10,
    "price": 50000
}
HTTP 201
[Asserts]
jsonpath "$.account_id" == 1
jsonpath "$.stock_code" == "STOCK02"
jsonpath "$.type" == "LIMIT"
jsonpath "$.direction" == "BUY"
jsonpath "$.quantity" == 10
jsonpath "$.price" == 50000
jsonpath "$.filled_quantity" == 0
jsonpath "$.status" == "PENDING"
[Captures]
buy_order_id: jsonpath "$.id"

# Test 4: Verify balance after buy order
GET http://localhost:8081/api/v1/accounts/1/balance
HTTP 200
[Asserts]
jsonpath "$.balance" == 500000

# Test 5: Create sell order
POST http://localhost:8081/api/v1/orders
Content-Type: application/json
{
    "account_id": 1,
    "stock_code": "STOCK01",
    "type": "LIMIT",
    "direction": "SELL",
    "quantity": 50,
    "price": 60000
}
HTTP 201
[Asserts]
jsonpath "$.account_id" == 1
jsonpath "$.stock_code" == "STOCK01"
jsonpath "$.type" == "LIMIT"
jsonpath "$.direction" == "SELL"
jsonpath "$.quantity" == 50
jsonpath "$.price" == 60000
jsonpath "$.filled_quantity" == 0
jsonpath "$.status" == "PENDING"
[Captures]
sell_order_id: jsonpath "$.id"

# Test 6: Verify holdings after sell order
GET http://localhost:8081/api/v1/accounts/1/holdings
HTTP 200
[Asserts]
jsonpath "$" count == 1
jsonpath "$[0].stock_code" == "STOCK01"
jsonpath "$[0].quantity" == 50

# Test 7: Cancel buy order
DELETE http://localhost:8081/api/v1/orders/{{buy_order_id}}
HTTP 200
[Asserts]
jsonpath "$.status" == "CANCELED"

# Test 8: Verify balance after canceling buy order
GET http://localhost:8081/api/v1/accounts/1/balance
HTTP 200
[Asserts]
jsonpath "$.balance" == 1000000

# Test 9: Cancel sell order
DELETE http://localhost:8081/api/v1/orders/{{sell_order_id}}
HTTP 200
[Asserts]
jsonpath "$.status" == "CANCELED"

# Test 10: Verify holdings after canceling sell order
GET http://localhost:8081/api/v1/accounts/1/holdings
HTTP 200
[Asserts]
jsonpath "$" count == 1
jsonpath "$[0].stock_code" == "STOCK01"
jsonpath "$[0].quantity" == 100

# Test 11: Test insufficient funds
POST http://localhost:8081/api/v1/orders
Content-Type: application/json
{
    "account_id": 1,
    "stock_code": "STOCK03",
    "type": "LIMIT",
    "direction": "BUY",
    "quantity": 100,
    "price": 20000
}
HTTP 400
[Asserts]
jsonpath "$.error" == "insufficient funds"

# Test 12: Test insufficient holding quantity
POST http://localhost:8081/api/v1/orders
Content-Type: application/json
{
    "account_id": 1,
    "stock_code": "STOCK01",
    "type": "LIMIT",
    "direction": "SELL",
    "quantity": 200,
    "price": 50000
}
HTTP 400
[Asserts]
jsonpath "$.error" == "insufficient holding quantity"

# Test 13: Test account not found
GET http://localhost:8081/api/v1/accounts/999/balance
HTTP 404
[Asserts]
jsonpath "$.error" == "account not found"

# Test 14: Test order not found
DELETE http://localhost:8081/api/v1/orders/999
HTTP 404
[Asserts]
jsonpath "$.error" == "order not found"