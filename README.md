# Mini Payment System API

Mini Payment System is a backend assignment project built with **Go, Gin, Gorm, PostgreSQL, Docker Compose, and Clean Architecture**.

The goal is to demonstrate more than basic CRUD by modeling a small fintech-like transfer workflow where **money is moved between two accounts safely and atomically**.

---

## Why this domain

I intentionally chose a **mini payment system** instead of plain CRUD because it highlights backend concerns that matter in real systems:

- transactional consistency
- business rule validation
- concurrency awareness
- predictable API contracts
- testable use case design

A transfer operation is a stronger signal of backend reasoning than a standard create/update/delete endpoint.

---

## Features

- CRUD for **Users**
- CRUD for **Accounts**
- Transfer money between accounts
- Prevent negative balances
- Use a **database transaction** for money transfer
- Use **row-level locking** during transfer to reduce race-condition risk
- Unit tests for transfer business logic
- Dockerized local development environment
- OpenAPI / Swagger-ready documentation
- Consistent API response envelope

---

## Tech stack

- **Go**
- **Gin**
- **Gorm**
- **PostgreSQL**
- **Docker / Docker Compose**
- **OpenAPI / Swagger**

---

## Project structure

```bash
cmd/server                      # application entrypoint / dependency wiring
internal/config                 # environment config loading
internal/database               # database bootstrap
internal/domain/entities        # domain entities
internal/domain/errs            # domain-level error definitions
internal/domain/repositories    # repository interfaces
internal/usecase                # business logic / application services
internal/infrastructure         # Gorm repository implementations
internal/delivery/http          # handlers, DTOs, HTTP routing
pkg/response                    # shared response envelope helpers
.docs / docs                    # generated swagger docs or manual OpenAPI file
```

---

## Architecture overview

This project follows a lightweight **Clean Architecture** approach.

### Layers

1. **Delivery layer**
   - Gin handlers
   - Request binding / validation
   - HTTP response formatting

2. **Use case layer**
   - Business rules
   - Transfer orchestration
   - Transaction boundary

3. **Domain layer**
   - Entities
   - Repository contracts
   - Domain errors

4. **Infrastructure layer**
   - Gorm repository implementations
   - PostgreSQL persistence
   - Transaction manager

### Request flow

```text
HTTP Request
→ Gin Handler
→ DTO Binding / Validation
→ Use Case
→ Repository Interface
→ Gorm Repository
→ PostgreSQL
→ HTTP Response Envelope
```

For the transfer endpoint:

```text
POST /api/v1/transactions
→ TransactionHandler.Create
→ TransactionUseCase.CreateTransfer
→ TxManager.WithinTransaction(...)
→ lock source & destination accounts
→ validate balance / currency
→ update both balances
→ insert transaction record
→ commit / rollback
```

---

## Domain model

### User
Represents an application user.

### Account
Represents a wallet or bank-like account owned by a user.

Key fields:
- `user_id`
- `balance`
- `currency`

### Transaction
Represents a money transfer between two accounts.

Key fields:
- `from_account_id`
- `to_account_id`
- `amount`
- `currency`
- `status`
- `reference`

---

## Business rules

The transfer use case enforces these rules:

- `from_account_id` and `to_account_id` must exist
- source and destination accounts must be different
- transfer amount must be greater than zero
- source balance must be sufficient
- both accounts must use the same currency
- balance updates and transaction insert must succeed together or fail together

---

## Concurrency and consistency

A transfer is wrapped in a **database transaction** and uses **row-level locking** (`SELECT ... FOR UPDATE`) when reading the involved accounts.

This helps reduce race conditions such as two concurrent requests reading the same balance and updating it incorrectly.

To further reduce deadlock risk, the implementation locks accounts in a **deterministic order** before applying the balance updates.

---

## API base URL

```text
http://localhost:8080/api/v1
```

---

## Response shape

All successful responses use a consistent envelope:

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

List endpoints include `meta.count`:

```json
{
  "success": true,
  "data": [],
  "meta": {
    "count": 0
  }
}
```

Error responses use a predictable structure:

```json
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_BALANCE",
    "message": "insufficient balance"
  }
}
```

---

## API endpoints

### Health
- `GET /health`

### Users
- `POST /users`
- `GET /users`
- `GET /users/{id}`
- `PUT /users/{id}`
- `DELETE /users/{id}`

### Accounts
- `POST /accounts`
- `GET /accounts`
- `GET /accounts/{id}`
- `PUT /accounts/{id}`
- `DELETE /accounts/{id}`

### Transactions
- `POST /transactions`
- `GET /transactions`
- `GET /transactions/{id}`

---

## Example requests

### Create user

```json
{
  "name": "Alice",
  "email": "alice@example.com"
}
```

### Create account

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "initial_balance": 100000,
  "currency": "THB"
}
```

### Transfer money

```json
{
  "from_account_id": "550e8400-e29b-41d4-a716-446655440000",
  "to_account_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "amount": 5000,
  "reference": "invoice-1001"
}
```

---

## OpenAPI / Swagger

This repository includes a manual OpenAPI file at:

```text
docs/openapi.yaml
```

You can:
- open it in Swagger Editor
- import it into Postman
- use it as the reviewable API contract for the assignment

If you want to regenerate Swagger docs from code annotations, install `swag` and run:

```bash
swag init -g cmd/server/main.go -o docs
```

---

## Running locally with Docker

```bash
docker compose up --build
```

Service endpoints:
- API: `http://localhost:8080`
- Health check: `http://localhost:8080/health`
- Swagger UI: `http://localhost:8080/swagger/index.html`

---

## Running locally without Docker

1. Start PostgreSQL
2. Copy environment values from `.env.example`
3. Run the application:

```bash
go run ./cmd/server
```

---

## Environment variables

```env
APP_ENV=development
APP_PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=arise_assignment
DB_SSLMODE=disable
```

---

## Running tests

```bash
go test ./...
```

Current unit tests focus on the most critical business flow:
- successful transfer
- insufficient balance
- invalid transfer request
- currency mismatch
- transaction creation failure

---

## Design decisions and trade-offs

### Why Clean Architecture
I wanted the business logic to remain independent from Gin and Gorm so that:
- transfer logic is easier to test
- HTTP concerns do not leak into domain rules
- infrastructure can evolve without changing use case logic

### Why use `int64` for money in this assignment
For simplicity, this project models money as an integer amount in the smallest unit.
In a production-grade financial system, I would prefer a more explicit money representation depending on the currency and precision requirements.

### Why use `AutoMigrate`
For an assignment, `AutoMigrate` keeps setup simple and reviewer-friendly.
For production, I would switch to a versioned migration tool such as `golang-migrate`.

### Why manual OpenAPI file is included
A checked-in OpenAPI file makes the API contract easy to review immediately without requiring the reviewer to generate docs locally first.

---

## Possible future improvements

Given more time, I would add:

- versioned database migrations
- pagination and filtering for list endpoints
- idempotency key support for transfers
- structured logging and request ID middleware
- integration tests against a real PostgreSQL instance
- account status handling such as `active` / `inactive`
- audit trail / ledger-style transaction design
- authentication and authorization

---

## Reviewer notes

The most important part of this project is the **transfer money** use case. That flow demonstrates:
- business rule enforcement
- transaction handling
- concurrency awareness
- layered design
- unit-testable application logic

That is the main reason this project goes beyond a typical CRUD assignment.
