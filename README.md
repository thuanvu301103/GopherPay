# GopherPay
The goal is to build a distributed system that handles user accounts, money transfers, and transaction history.

## System Architecture
Instead of one big application, we will split the logic into three core services:
- Identity Service: Manages user registration, profiles, and authentication (JWT).
- Account Service: Manages bank accounts, wallet balances, and currency types.
- Transaction Service: Handles the logic of moving money from Wallet A to Wallet B (The "Engine").

### Service Breakdown

#### Service A: Identity & Auth
- Focus: Security and Middleware.
- Key Tasks: Implement POST /register and POST /login.
- Use Bcrypt for password hashing.
- Generate JWT (JSON Web Tokens) for authenticated requests.

#### Service B: Account Management
- Focus: CRUD operations and Database design.
- Key Tasks:
    - GET /accounts/me: Fetch current balance.
    - POST /accounts: Create a new currency wallet (e.g., USD, VND).

#### Service C: Transaction Engine (The Core)
- Focus: Concurrency and Atomicity.
- Key Tasks: POST /transfer: Transfer money between users.
- Crucial Concept: Using Database Transactions to ensure that if the sender's deduction fails, the receiver's credit never happens.

## API Contract Document

- Update document:
```Bash
swag init -g cmd/api/main.go
```

- Endpoint: `http://localhost:3000/docs/index.html`

## Database Migration

### Using Atlas
- Change the `DB_AUTO_MIGRATE` in `.env` to `true`
- Create migration file: 
```bash
atlas migrate diff <file_name> --env gorm
```
- Run migration:
```bash
atlas migrate apply --env gorm --url "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

### Using AutoMigration
- Change the `DB_AUTO_MIGRATE` in `.env` to `false`
- Run the server

