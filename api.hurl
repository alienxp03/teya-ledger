# Get transactions
GET http://localhost:8080/api/v1/transactions?accountNumber=ACCOUNT_NUMBER_1
Authorization: USER_TOKEN_1
HTTP 200
[Asserts]
jsonpath "$.transactions" count >= 1

# Get transactions
GET http://localhost:8080/api/v1/transactions?accountNumber=ACCOUNT_NUMBER_2
Authorization: USER_TOKEN_1
HTTP 200
[Asserts]
jsonpath "$.transactions" count == 0

# Get transactions with pagination
GET http://localhost:8080/api/v1/transactions?limit=1&page=1&accountNumber=111222
Authorization: USER_TOKEN_1
HTTP 200
[Asserts]
jsonpath "$.transactions" isCollection

# POST deposits
POST http://localhost:8080/api/v1/deposits
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{newUuid}}",
    "accountNumber": "ACCOUNT_NUMBER_1",
    "amount": 1000,
    "currency": "MYR",
    "description": "description"
}
HTTP 200
[Asserts]
jsonpath "$.transaction.transactionID" exists
jsonpath "$.transaction.status" == "pending"
jsonpath "$.transaction.amount" > 0
jsonpath "$.transaction.currency" == "MYR"
jsonpath "$.transaction.description" == "description"
jsonpath "$.transaction.createdAt" exists
jsonpath "$.transaction.updatedAt" exists
[Captures]
depositTransactionID: jsonpath "$.transaction.transactionID"

# duplicate transactionID
POST http://localhost:8080/api/v1/deposits
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{depositTransactionID}}",
    "accountNumber": "ACCOUNT_NUMBER_1",
    "amount": 100,
    "currency": "MYR",
    "description": "description"
}
HTTP 400
[Asserts]
jsonpath "$.message" contains "already exists"

# POST deposits with invalid account number
POST http://localhost:8080/api/v1/deposits
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{newUuid}}",
    "accountNumber": "ACCOUNT_NUMBER_2",
    "amount": 100,
    "currency": "MYR",
    "description": "description"
}
HTTP 404
[Asserts]
jsonpath "$.code" contains "NOT_FOUND"
jsonpath "$.message" contains "not found"

# POST deposits with invalid amount
POST http://localhost:8080/api/v1/deposits
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{newUuid}}",
    "accountNumber": "ACCOUNT_NUMBER_1",
    "amount": -100,
    "currency": "MYR",
    "description": "description"
}
HTTP 400
[Asserts]
jsonpath "$.message" contains "Invalid"

# POST deposits with invalid currency
POST http://localhost:8080/api/v1/deposits
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{newUuid}}",
    "accountNumber": "ACCOUNT_NUMBER_1",
    "amount": 100,
    "currency": "USD",
    "description": "description"
}
HTTP 400
[Asserts]
jsonpath "$.message" contains "Invalid"

# POST withdrawals
POST http://localhost:8080/api/v1/withdrawals
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{newUuid}}",
    "accountNumber": "ACCOUNT_NUMBER_1",
    "amount": -100,
    "currency": "MYR",
    "description": "withdrawal description"
}
HTTP 200
[Asserts]
jsonpath "$.transaction.transactionID" exists
jsonpath "$.transaction.status" == "pending"
jsonpath "$.transaction.amount" == -100
jsonpath "$.transaction.currency" == "MYR"
jsonpath "$.transaction.description" == "withdrawal description"
jsonpath "$.transaction.createdAt" exists
jsonpath "$.transaction.updatedAt" exists
[Captures]
withdrawalTransactionID: jsonpath "$.transaction.transactionID"

# duplicate withdrawal transactionID
POST http://localhost:8080/api/v1/withdrawals
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{withdrawalTransactionID}}",
    "accountNumber": "ACCOUNT_NUMBER_1",
    "amount": -100,
    "currency": "MYR",
    "description": "withdrawal description"
}
HTTP 400
[Asserts]
jsonpath "$.message" contains "already exists"

# POST withdrawals with invalid account number
POST http://localhost:8080/api/v1/withdrawals
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{newUuid}}",
    "accountNumber": "ACCOUNT_NUMBER_2",
    "amount": -100,
    "currency": "MYR",
    "description": "withdrawal description"
}
HTTP 404
[Asserts]
jsonpath "$.code" contains "NOT_FOUND"
jsonpath "$.message" contains "not found"

# POST withdrawals with invalid amount
POST http://localhost:8080/api/v1/withdrawals
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{newUuid}}",
    "accountNumber": "ACCOUNT_NUMBER_1",
    "amount": 100,
    "currency": "MYR",
    "description": "withdrawal description"
}
HTTP 400
[Asserts]
jsonpath "$.message" contains "Invalid"

# POST withdrawals with invalid currency
POST http://localhost:8080/api/v1/withdrawals
Authorization: USER_TOKEN_1
Content-Type: application/json
{
    "transactionID": "{{newUuid}}",
    "accountNumber": "ACCOUNT_NUMBER_1",
    "amount": -10,
    "currency": "USD",
    "description": "withdrawal description"
}
HTTP 400
[Asserts]
jsonpath "$.message" contains "Invalid"

# Get balance
GET http://localhost:8080/api/v1/balances?accountNumber=ACCOUNT_NUMBER_1
Authorization: USER_TOKEN_1
Content-Type: application/json
HTTP/1.1 200
[Asserts]
jsonpath "$.balance.amount" > 0
jsonpath "$.balance.currency" == "MYR"

# Get balance - invalid account
GET http://localhost:8080/api/v1/balances?accountNumber=INVALID_ACCOUNT
Authorization: USER_TOKEN_1
Content-Type: application/json
HTTP/1.1 404
[Asserts]
jsonpath "$.code" == "NOT_FOUND"
jsonpath "$.message" == "not found"

# Get balance - unauthorized
GET http://localhost:8080/api/v1/balances
Content-Type: application/json
HTTP/1.1 401
