# GopherPay
The goal is to build a distributed system that handles user accounts, money transfers, and transaction history.

## System Architecture
Instead of one big application, we will split the logic into three core services:
- Identity Service: Manages user registration, profiles, and authentication (JWT).
- Account Service: Manages bank accounts, wallet balances, and currency types.
- Transaction Service: Handles the logic of moving money from Wallet A to Wallet B (The "Engine").

Utility/Support Service:
- Nottification Service: A dedicated, asynchronous support service responsible for managing all outgoing communications between the system and its users

### Core Service Breakdown

#### Service A: Identity & Auth
- Focus: Security and Middleware.
- Key Tasks: 
    - `POST /auth/register` 
    - `POST /auth/login`.
- Use Bcrypt for password hashing.
- Generate JWT (JSON Web Tokens) for authenticated requests.

#### Service B: Account Management
- Focus: CRUD operations and Database design.
- Key Tasks:
    - `GET /accounts/me`: Fetch current balance.
    - `POST /accounts`: Create a new currency wallet (e.g., USD, VND).

#### Service C: Transaction Engine (The Core)
- Focus: Concurrency and Atomicity.
- Key Tasks: `POST /transfer`: Transfer money between users.
- Crucial Concept: Using Database Transactions to ensure that if the sender's deduction fails, the receiver's credit never happens.

## Security
The most secure way is to use the openssl tool, which is pre-installed on Linux, macOS, and Git Bash for Windows.

- For `JWT_SECRET`: Run this to get a long random string:
```Bash
openssl rand -base64 32
```

- For `STORE_ENCRYPTION_KEY`: Run this to get exactly 32 hex characters:
```Bash
openssl rand -hex 16
```