# GopherPay
The goal is to build a distributed system that handles user accounts, money transfers, and transaction history.

## System Architecture
Instead of one big application, we will split the logic into three core services:
- Identity Service: Manages user registration, profiles, and authentication (JWT).
- Account Service: Manages bank accounts, wallet balances, and currency types.
- Transaction Service: Handles the logic of moving money from Wallet A to Wallet B (The "Engine").

## Service Breakdown

### Service A: Identity & Auth
- Focus: Security and Middleware.
- Key Tasks: * Implement POST /register and POST /login.
- Use Bcrypt for password hashing.
- Generate JWT (JSON Web Tokens) for authenticated requests.

### Service B: Account Management
- Focus: CRUD operations and Database design.
- Key Tasks:
    - GET /accounts/me: Fetch current balance.
    - POST /accounts: Create a new currency wallet (e.g., USD, VND).

### Service C: Transaction Engine (The Core)
- Focus: Concurrency and Atomicity.
- Key Tasks: POST /transfer: Transfer money between users.
- Crucial Concept: Using Database Transactions to ensure that if the sender's deduction fails, the receiver's credit never happens.