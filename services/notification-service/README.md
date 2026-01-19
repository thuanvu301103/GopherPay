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
- Represents the top-level workspace that contains all your projects, environments, workflows, integration and settings
- The user who creates the organization will become its admin (automatically create a member relation with role `"admin"`)

### Environment 
- Represents an isolated workspace within an organization used to separate different stages of your notification system. Each environment has its own API keys, workflows, subscribers, templates, and settings. This allows you to safely test changes without affecting production.
- Common environments include:
    - Development – for local testing and experimentation
    - Staging – for pre‑production validation
    - Production – for live notification delivery
- Two default Environment entities are created when an Organization is created:
    - `Development` for local testing and experimentation
    - `Production` for live notification delivery

### Integration
- The "delivery arm" of the system. While the Workflow contains the logic (what to send and when), the Integration connects Novu to the actual service provider that delivers the message to the user.

- Novu does not send emails, SMS, or Push notifications itself; it orchestrates external providers to do the work.
    1. Key Responsibilities
        - Connectivity: It stores the credentials (API Keys, Secret Tokens, SMTP settings) required to talk to providers like SendGrid, Twilio, or AWS SES.
        - Abstraction: It allows you to switch providers without changing your Go code. If you switch from SendGrid to Postmark, you only update the Integration entity in Novu.
        - Failover: You can have multiple integrations for the same channel to ensure high availability.

    2. Supported Channels
        - Email: SendGrid, Mailgun, Postmark, AWS SES, SMTP.
        - SMS: Twilio, Plivo, MessageBird.
        - Push: Firebase (FCM), APNS (Apple).
        - Chat: Slack, Discord, MS Teams.
        - In-App: Novu's internal WebSocket-based notification center.

### Workflow/Template
- A Workflow is the "Master Blueprint" of a notification. It defines what to send, how to send it, and when it should be delivered.
- A Workflow in Novu is made up of steps, and each step defines a specific action in the notification process. Steps run in the order you arrange them, and they control how and when the subscriber receives a notification.

### Message/Activity

## Main Workflows
- *Known Issue*: The latest Novu Docker image is currently affected by a critical bug. Detailed technical information and tracking can be found in [Github Issue](https://github.com/novuhq/novu/issues/9569). Due to this issue, I suggest using API calling instead of UI.
- *Custom API Document*: Describe in `docs` folder. Import to Postman to use.
- *Caution*: These workflows are used for *email verification* feature only.
- *Base URL*: `http://localhost:3004`.

### Initial Root Admin & Organization
1. Register a new User (admin/developer): `POST /v1/auth/register`
2. Login: `POST /v1/auth/login`
3. Copy the returned token
4. Create an Organization: `POST /v1/organizations`
5. A member relation between the User and the Organization is created automatically with role `"admin"`

### Configure the Email Integration (GMail SMTP)
1. Connection Settings:
    - Host: `smtp.gmail.com`
    - Port: `465` (SSL) or `587` (TLS)
    - User/From: your full Gmail address (e.g., example@gmail.com)
2. Password (Google App Password): Since 2022, Google requires a 16-character unique code for third-party apps like Novu.
    1. Enable 2-Step Verification: Go to your [Google Security](https://myaccount.google.com/security) Settings and ensure 2FA is active.
    2. Generate App Password:
        - Search for "App Passwords" in the top search bar of your Google Account.
        - Alternatively, go to `Security > 2-Step Verification` scroll to the bottom to find App Passwords.
        - Create: Enter an app name (e.g., "Novu-Notification-Service") and click Create.
        - Copy: A 16-character code will appear . Save this code—it is the value you use for the "password" field in your Novu integration.
3. Create integration with body: `POST /v1/integrations`
```JSON
{
  "providerId": "nodemailer",
  "channel": "email",
  "_environmentId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
  "credentials": {
    "host": "smtp.gmail.com",
    "port": "465",
    "secure": true,
    "user": "admin@yourdomain.com",
    "password": "xxxx-xxxx-xxxx-xxxx",
    "from": "admin@yourdomain.com"
  },
  "active": true,
  "check": false
}
```
*Caution*: `check` (Optional): If set to true, Novu will attempt to verify the connection during creation. It is recommended to keep it false if you just want to save the settings first. 

### Create Notification Workflow (Verification Email Workflow)
1. Define Workflow steps (only 1 step):
    - `type`: Defines the communication channel. In this case, `"email"`. Novu will automatically route this to your active Email Integration (Gmail).
    - `template`: Contains the actual message details:
        - `subject`: The email's subject line. Supports variables like `{{otpCode}}`.
        - `content`: The body of the email (HTML supported). Variables used here must match the data sent from your Go backend.
        - `layoutIdentifier`: Points to a specific email layout (header/footer). `"default"` uses Novu's standard system layout.

### Trigger the Notification (Trigger the Workflow)
