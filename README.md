# Banking System

Microservices banking system with REST API, gRPC communication, and idempotent transactions.


Three services:
- **Web Server** - REST API (port 8080)
- **DB Service** - Database operations (gRPC)
- **Idempotency Service** - Duplicate transaction prevention (gRPC)

## Quick Start
```bash
docker-compose up --build
```

Wait for services to start, then test at `http://localhost:8080`

## API Usage

**1. Login**
```bash
POST /login
{
"username": "test", 
"password": "test1234"
}

Response: {"error": 0, "token": "abcdef"}
```

**2. Deposit**
```bash
POST /deposit
{
  "id": "unique-tx-id",
  "token": "abcdef",
  "amount": 100,
  "timestamp": 1663237199324
}

Response: {"error": 0, "balance": 650}
```

**3. Withdraw**
```bash
POST /withdraw
{
  "id": "unique-tx-id",
  "token": "abcdef",
  "amount": 30,
  "timestamp": 1663237199324
}

Response: {"error": 0, "balance": 620}
```

**Note:** Sending the same transaction `id` twice returns the same result.

## Key Features

Token authentication (1h expiry)
Idempotent transactions (duplicate prevention via Redis)
Insufficient funds validation
gRPC for internal services, REST for external API

## Tech Stack

Go • gRPC • MySQL • Redis • Docker • Protocol Buffers


## Testing

Test user: `test` / `test1234`

Use Postman, curl, or any HTTP client. All endpoints require `POST` method with JSON body.
