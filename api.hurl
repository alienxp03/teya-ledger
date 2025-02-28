GET http://localhost:8080/api/v1/transactions
HTTP 200
Authorization: Bearer USER_TOKEN_1
[Asserts]
jsonpath "$.data.transactions" count == 0
