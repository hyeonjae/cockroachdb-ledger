# Mini-Ledger

A stock trading API server built with Go and CockroachDB using Clean Architecture principles.

## Features

- Account balance management
- Stock holdings tracking
- Buy/Sell order creation and cancellation
- Transaction-based operations with proper error handling
- RESTful API with JSON responses
- CockroachDB for scalable, distributed database

## Architecture

The project follows Clean Architecture patterns with clear separation of concerns:

```
mini-ledger/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── api/                     # HTTP handlers and routes
│   ├── service/                 # Business logic
│   ├── repository/              # Data access layer
│   ├── domain/                  # Domain models and errors
│   ├── config/                  # Configuration management
│   └── db/                      # Database connection and migrations
├── migrations/                  # SQL migration files
├── tests/                       # Hurl API tests
├── Dockerfile                   # Container image
├── docker-compose.yml          # Development environment with CockroachDB
└── go.mod                       # Go module
```

## Technology Stack

- **Go 1.23** - Programming language
- **Chi** - HTTP router
- **CockroachDB** - Distributed SQL database
- **SQLX** - SQL toolkit with PostgreSQL driver
- **Fx** - Dependency injection framework
- **Docker** - Containerization

## API Endpoints

### Get Account Balance
```
GET /api/v1/accounts/{accountID}/balance
```
Response:
```json
{"account_number": "AC001", "balance": 1000000}
```

### Get Account Holdings
```
GET /api/v1/accounts/{accountID}/holdings
```
Response:
```json
[{"stock_code": "STOCK01", "quantity": 100}]
```

### Create Order
```
POST /api/v1/orders
Content-Type: application/json
```
Request:
```json
{
  "account_id": 1,
  "stock_code": "STOCK01",
  "type": "LIMIT",
  "direction": "BUY",
  "quantity": 10,
  "price": 50000
}
```
Response:
```json
{
  "id": 1,
  "account_id": 1,
  "stock_code": "STOCK01",
  "type": "LIMIT",
  "direction": "BUY",
  "quantity": 10,
  "price": 50000,
  "filled_quantity": 0,
  "status": "PENDING",
  "created_at": "2024-01-01T10:00:00Z"
}
```

### Cancel Order
```
DELETE /api/v1/orders/{orderID}
```
Returns the canceled order information.

## Business Logic

### Buy Orders
1. Verify account exists
2. Check sufficient balance (balance >= price × quantity)
3. Deduct amount from account balance
4. Create order with PENDING status
5. All operations in a transaction

### Sell Orders
1. Verify account exists
2. Check sufficient holdings (holdings >= quantity)
3. Deduct from holdings
4. Create order with PENDING status
5. All operations in a transaction

### Order Cancellation
1. Verify order exists
2. Check order is cancelable (PENDING or PARTIAL status)
3. Restore funds (buy) or holdings (sell) for unfilled quantity
4. Update order status to CANCELED
5. All operations in a transaction

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `400 Bad Request` - Invalid input or business rule violations
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server errors

Error response format:
```json
{"error": "error message"}
```

## Environment Variables

- `DATABASE_URL` - CockroachDB connection string (default: "postgresql://root@localhost:26257/mini_ledger?sslmode=disable")
- `HTTP_PORT` - HTTP server port (default: "8080")

## Quick Start

### Using Docker Compose (Recommended)
```bash
# Clone the repository
git clone <repository-url>
cd mini-ledger

# Start CockroachDB and the application
docker-compose up --build

# The API server will be available at http://localhost:8081
# CockroachDB UI will be available at http://localhost:8080
```

### Manual Setup
```bash
# Start CockroachDB locally
cockroach start-single-node --insecure --store=node1 --listen-addr=localhost:26257 --http-addr=localhost:8080

# Create database and run migrations
cockroach sql --insecure --host=localhost:26257 -e "CREATE DATABASE mini_ledger;"

# Install dependencies
go mod download

# Run the server
DATABASE_URL="postgresql://root@localhost:26257/mini_ledger?sslmode=disable" go run cmd/server/main.go
```

## Testing

The project includes comprehensive API tests using Hurl:

```bash
# Make sure the services are running
docker-compose up -d

# Wait for services to be ready
sleep 30

# Run API tests
hurl --test tests/api_tests.hurl

# Or run the interactive test script
./test_api.sh
```

Test scenarios covered:
- Account balance retrieval
- Holdings management
- Order creation (buy/sell)
- Order cancellation
- Error cases (insufficient funds, holdings, etc.)

## CockroachDB Schema

The application uses CockroachDB with the following tables:

- **accounts** - User accounts with balances (DECIMAL for precision)
- **holdings** - Stock holdings per account with unique constraints
- **orders** - Trading orders with status tracking

CockroachDB-specific features used:
- SERIAL PRIMARY KEY for auto-incrementing IDs
- STRING data type for text fields
- DECIMAL(15,2) for monetary values
- TIMESTAMPTZ for timestamps with timezone
- Foreign key constraints with proper referential integrity
- UPSERT operations with ON CONFLICT clauses

Initial test data includes:
- Account AC001 with 1,000,000 balance
- 100 shares of STOCK01

## Development

### Project Structure
- Clean Architecture with dependency injection
- Repository pattern for data access
- Service layer for business logic
- HTTP handlers for API endpoints
- CockroachDB-compatible SQL queries

### Dependencies
- Automatic database migration on startup
- Transaction-based operations with CockroachDB's serializable isolation
- Comprehensive error handling
- JSON API responses

## Building

### Docker Build
```bash
docker build -t mini-ledger .
```

### Go Build
```bash
go build -o mini-ledger ./cmd/server
```

## CockroachDB Features

This implementation leverages CockroachDB's distributed SQL capabilities:

1. **ACID Transactions** - All trading operations are transactional
2. **Serializable Isolation** - Prevents race conditions in concurrent trades
3. **Automatic Retries** - CockroachDB handles serialization conflicts
4. **Horizontal Scaling** - Can be scaled across multiple nodes
5. **SQL Standard** - Uses standard SQL with PostgreSQL compatibility

## License

This project is for educational purposes.