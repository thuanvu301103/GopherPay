# GopherPay
The goal is to build a distributed system that handles user accounts, money transfers, and transaction history.

## System Architecture
Instead of one big application, we will split the logic into three core services:
- Identity Service: Manages user registration, profiles, and authentication (JWT).
- Account Service: Manages bank accounts, wallet balances, and currency types.
- Transaction Service: Handles the logic of moving money from Wallet A to Wallet B (The "Engine").

Utility/Support Service:
- Nottification Service (Novu): A dedicated, asynchronous support service responsible for managing all outgoing communications between the system and its users

Infrastructures:
- Message Broker (Apache Kafka): The asynchronous backbone that decouples the core services from the support services

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

## Connecting Separate Docker Compose Projects

### Create the Shared Network

1. Create a global network that exists outside of any specific `docker-compose.yml` file. Open terminal (PowerShell) and run:

```Bash
docker network create gopher_network
```

2. Verify it exists by running 

```Bash
docker network ls
```

### Update the Docker Compose Files

Tell each service to join this external network. Add the networks configuration to the `kafka_connect` (and all other services) service and define the network at the bottom.
```YAML
services:
  kafka_connect:
    # ... existing config ...
    networks:
      - gopher_net

networks:
  gopher_net:
    external: true
    name: gopher_network
```

### Restart Services

Navigate to each directory and restart the containers to apply the network changes:
```Bash
docker compose up -d
```

### Verification

1. Check if both containers have successfully joined the "bridge":

```Bash
docker network inspect gopher_network
```

2. Look for the Containers section. You should see all services listed there with their internal IPs.

## Step-by-step Start Services

1. Start Infra: Message Broker (Apache Kafka)
2. Setup Service: Notification Service
   1. Start service
   2. Run seed data script `seed.sh`
   3. Retrieve `API-key` from seed script console or from API
3. Run seed data script for Infra: Message Broke (Apache Kafaka) `seed.sh`
4. Start Service: Auth Service
