# Notification Service

## Run Service using Docker
```Bash
docker-compose up -d
```

## Important Endpoints

### The Dashboard (User Interface)
- URL: `http://localhost:4200`
- Purpose: This is where you (the developer/admin) log in
- Key Actions: 
    - Create Workflows (e.g., "email-verification").
    - Design Templates using a Drag-and-Drop editor or HTML.
    - Configure Integrations (Connect your SendGrid, Twilio, or AWS SES accounts).
    - Monitor Activity Feed (See which emails were sent and if they were opened).
    
### The API Endpoints (For your Identity Service)
- Your Identity Service (Go) will communicate with Novu via the api service at `http://localhost:3004`.
- The most important endpoint you will use is the Trigger endpoint:
| Method | Endpoint | Description | 
| --- | --- | --- |
| POST | `/v1/events/trigger` | Used by your Identity Service to send a notification | 
| GET | `/v1/subscribers` | To manage your user profiles inside Novu | 
| GET | `/v1/workflows` | To fetch the list of available notification workflows | 
| GET | `/v1/notifications/stats` | To get analytics for your dashboard | 

### Health & Internal Endpoints
These are useful for checking if your services are running correctly:
- API Health Check: `http://localhost:3004/v1/health-check`
- Swagger UI - Documentation (DEPRECATED): Usually available at `http://localhost:3004/api` (depending on the version, it provides a full list of all available REST endpoints).
- OpenAPI JSON, YAML - (Import to `Postman`): `http://localhost:3004/api-json` or `http://localhost:3004/api-yaml`

## Core Entities

### User
An internal entity that manages or interacts with the Novu platform, such as developers or administrators

### Subscriber
The recipient of the notification. Every subscriber has a unique subscriberId (usually matching your database user ID). They store contact metadata like email, phone number, and device tokens for push notifications

### Organization
- Represents the top-level workspace that contains all your projects, environments, workflows, and settings
- The user who creates the organization will become its admin (automatically create a member relation with role `"admin"`)

## Main Workflows
- *Known Issue*: The latest Novu Docker image is currently affected by a critical bug. Detailed technical information and tracking can be found in [Github Issue](https://github.com/novuhq/novu/issues/9569). Due to this issue, I suggest using API calling instead of UI
- *Custom API Document*: Describe in `docs` folder 
- *Caution*: These workflows are used for *email verification* feature only
- *Base URL*: `http://localhost:3004`

### Initial
1. Register a new User (admin/developer): `POST /v1/auth/register`
2. Login: `POST /v1/auth/login`
3. Copy the returned token
4. Create an Organization: `POST /v1/organizations`
5. A member relation between the User and the Organization is created automatically with role `"admin"`

### 