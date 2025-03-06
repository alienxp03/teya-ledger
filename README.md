# Teya Ledger API

A simple ledger API for managing transactions.

## Getting Started

### Assumptions

- Operations are not atomic.
- All operations are idempotent.
- No database is used for persistence. Each data will be reset on each server run.
- No logging.
- Each operation require `Authorization` header. There are two default users:
  - `USER_TOKEN_1` with `ACCOUNT_NUMBER_1`
  - `USER_TOKEN_2` with `ACCOUNT_NUMBER_2`

### Folder structure

- `/api`
  - API implementation
  - Keep it lightweight. Ideally we only want to deal with the API responses. No business logic should be included here.
- `/cmd`
  - Run command
- `/db`
  - Database connector
- `/handler`
  - Logic handler. This is where the business logic is implemented.
- `/server`
  - Handle howe we run the server
- `/storage`
  - Storage implementation. This is where the data is persisted.
  - We use an in-memory storage for now which will be reset on each server run.
- `/types`
  - Shared types between packages

### Prerequisites

- Go 1.24 or later
- Docker and Docker Compose (optional)

### Running Locally

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the application:
   ```bash
   make run
   ```

### Running with Docker

1. Build and start the containers:

   ```bash
   make docker
   ```

2. Stop the containers:
   ```bash
   make docker-down
   ```

## API Endpoints

- All endpoints require a `Authorization` header for authentication.
- Example: `Authorization: <token>`

### Deposits

- **POST** `/api/v1/deposits`
  - Create a new deposit transaction
  - Each request will create a new transaction with a status of `pending`.
  - The transaction will be updated to `completed` after a delay of 200ms.
  - Request body:
    ```json
    {
      "transactionID": "string",      # required, must be unique for each request. Ideally UUID
      "accountNumber": "string",      # required
      "amount": number,               # required, must be positive. In cents value
      "currency": "string",           # required, only "MYR" is supported
      "description": "string"         # required
    }
    ```
  - Response:
    ```json
    {
      "transaction": {
        "transactionID": "string",
        "status": "string",
        "amount": number,
        "currency": "string",
        "description": "string",
        "createdAt": "string",
        "updatedAt": "string"
      }
    }
    ```

### Withdrawals

- **POST** `/api/v1/withdrawals`
  - Create a new withdrawal transaction
  - Request body:
    ```json
    {
      "transactionID": "string",  # required, must be unique for each request. Ideally UUID
      "accountNumber": "string",  # required
      "amount": number,           # required, must be negative. In cents value
      "currency": "string",       # required, only "MYR" is supported
      "description": "string"     # required
    }
    ```
  - Response:
    ```json
    {
      "transaction": {
        "transactionID": "string",
        "status": "string",
        "amount": number,
        "currency": "string",
        "description": "string",
        "createdAt": "string",
        "updatedAt": "string"
      }
    }
    ```

### Balance

- **GET** `/api/v1/balances?accountNumber=string`
  - Get the current balance for an account
  - Query parameters:
    - `accountNumber`: The account number to check balance for. Required.
  - Response:
    ```json
    {
      "balance": {
        "amount": number,
        "currency": "string"
      }
    }
    ```

### Transactions

- **GET** `/api/v1/transactions`

  - Get a list of transactions
  - Query parameters:
    - `accountNumber`: Filter by account number. Required.
    - `limit`: Maximum number of transactions to return (default: 10)
    - `page`: Page number for pagination (default: 1)
  - Response:
    ```json
    {
      "transactions": [
        {
          "transactionID": "string",
          "status": "string",
          "amount": number,
          "currency": "string",
          "description": "string",
          "createdAt": "string",
          "updatedAt": "string"
        }
      ]
    }
    ```

- **GET** `/api/v1/transactions/{transactionID}`

  - Get details of a specific transaction
  - Path parameters:
    - `transactionID`: The ID of the transaction to retrieve
  - Response:
    ```json
    {
      "transaction": {
        "transactionID": "string",
        "status": "string",
        "amount": number,
        "currency": "string",
        "description": "string",
        "createdAt": "string",
        "updatedAt": "string"
      }
    }
    ```

## Manual Tests

- You can manually test the API using `curl` with the following steps (assuming you have the server running):

  - Start the server
  - Run the following commands:

    - Get user balance

    ```bash
    curl -X GET "http://localhost:8080/api/v1/balances?accountNumber=ACCOUNT_NUMBER_1" \
      -H "Authorization: USER_TOKEN_1"
    ```

    - Get transactions

    ```bash
    curl -X GET "http://localhost:8080/api/v1/transactions?accountNumber=ACCOUNT_NUMBER_1" \
      -H "Authorization: USER_TOKEN_1"
    ```

    - Create deposit

    ```bash
    curl -X POST http://localhost:8080/api/v1/deposits \
        -H "Authorization: USER_TOKEN_1" \
        -H "Content-Type: application/json" \
        -d '{
        "transactionID": "REPLACE_UNIQUE_ID_1",
        "accountNumber": "ACCOUNT_NUMBER_1",
        "amount": 1000,
        "currency": "MYR",
        "description": "deposit description"
        }'
    ```

    - Create withdrawal

      ```bash
      curl -X POST http://localhost:8080/api/v1/withdrawals \
        -H "Authorization: USER_TOKEN_1" \
        -H "Content-Type: application/json" \
        -d '{
        "transactionID": "REPLACE_UNIQUE_ID_2",
        "accountNumber": "ACCOUNT_NUMBER_1",
        "amount": -1,
        "currency": "MYR",
        "description": "withdrawal description"
           }'
      ```

    - Get transaction details

    ```bash
    curl -X GET http://localhost:8080/api/v1/transactions/123456 \
      -H "Authorization: USER_TOKEN_1"
    ```

    - Get user balance

    ```bash
    curl -X GET "http://localhost:8080/api/v1/balances?accountNumber=ACCOUNT_NUMBER_1" \
      -H "Authorization: USER_TOKEN_1"
    ```

## Automated Tests

- You can refer to api.hurl file for the integration tests.

### Running Tests

1. Unit tests:

   ```bash
   make test
   ```

2. API tests

   - Integration tests using hurl
   - Requires `hurl` to be installed
   - Requires server to be running

   ```bash
   make run       # Start server
   make api_test  # Run tests
   ```

3. Coverage report:
   ```bash
   make coverage
   ```

### Running Tests in Docker

1. Unit tests:

   ```bash
   make docker-test
   ```

2. API tests:

   - Integration tests using hurl
   - Requires server to be running

   ```bash
   make docker
   make docker-api-test
   ```

3. Coverage report:
   ```bash
   make docker-coverage
   ```
