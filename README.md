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
  - Request body:
    ```json
    {
      "transactionID": "string",      # required, must be unique for each request. Ideally UUID
      "accountNumber": "string",      # required
      "amount": number,               # required, must be positive
      "currency": "string",           # required, only "MYR" is supported
      "description": "string"         # required
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
      "amount": number,           # required, must be negative
      "currency": "string",       # required, only "MYR" is supported
      "description": "string"     # required
    }
    ```

### Balance

- **GET** `/api/v1/balances?accountNumber=string`
  - Get the current balance for an account
  - Query parameters:
    - `accountNumber`: The account number to check balance for. Required.

### Transactions

- **GET** `/api/v1/transactions`

  - Get a list of transactions
  - Query parameters:
    - `accountNumber`: Filter by account number. Required.
    - `limit`: Maximum number of transactions to return (default: 10)
    - `page`: Page number for pagination (default: 1)

- **GET** `/api/v1/transactions/{transactionID}`

  - Get details of a specific transaction
  - Path parameters:
    - `transactionID`: The ID of the transaction to retrieve

## Testing

### Running Tests

1. Unit tests:

   ```bash
   make test
   ```

2. API tests:

   ```bash
   make api_test
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

   ```bash
   make docker-api-test
   ```

3. Coverage report:
   ```bash
   make docker-coverage
   ```

