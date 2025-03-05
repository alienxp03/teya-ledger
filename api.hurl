# Get transactions
GET http://localhost:8080/api/v1/transactions?accountNumber=ACCOUNT_NUMBER_1
Authorization: Bearer USER_TOKEN_1
HTTP 200
[Asserts]
jsonpath "$.transactions" count == 1

# Get transactions
GET http://localhost:8080/api/v1/transactions?accountNumber=ACCOUNT_NUMBER_2
Authorization: Bearer USER_TOKEN_1
HTTP 200
[Asserts]
jsonpath "$.transactions" count == 0

# Get transactions with pagination
GET http://localhost:8080/api/v1/transactions?limit=1&page=1&accountNumber=111222
Authorization: Bearer USER_TOKEN_1
HTTP 200
[Asserts]
jsonpath "$.transactions" isCollection

# Get transactions with invalid query parameter
GET http://localhost:8080/api/v1/transactions?limit=1&page=1&invalidParam=1&accountNumber=111222
Authorization: Bearer USER_TOKEN_1
HTTP 400
[Asserts]
jsonpath "$.error" contains "invalid path"
