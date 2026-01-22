# Notification Service

## Run Service using Docker
```Bash
docker compose down
docker compose up -d --pull always
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

### Workflow/Notification Template
- A Workflow is the "Master Blueprint" of a notification. It defines what to send, how to send it, and when it should be delivered.
- A Workflow in Novu is made up of steps, and each step defines a specific action in the notification process. Steps run in the order you arrange them, and they control how and when the subscriber receives a notification.

### Notification
Represent the complete journey of a message triggered by an event. They encapsulate all the execution logic and delivery metadata in a single traceable unit

### Message
Message is a single notification that is sent to a subscriber. Each channel step in the workflow generates one or more message(s)

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
6. Make sure that at least one Environement belongs to user's Organization. Else, create a default Environtment `POST /v1/environments`

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
  "type": "email",
  "_environmentId": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx",
  "credentials": {
    "host": "smtp.gmail.com",
    "port": "465",
    "secure": true,
    "user": "admin@yourdomain.com",
    "password": "xxxx-xxxx-xxxx-xxxx",
    "from": "admin@yourdomain.com",
    "user": "admin@yourdomain.com"
  },
  "active": true,
  "check": false
}
```
4. API Service creates an NotoficationTemplate Entity:
```JSON
{
  "_id": {
    "$oid": "6971ab0493e98bfba5435ba9"
  },
  "name": "Email Verification",
  "description": "Email Verification",
  "active": true,
  "type": "workflow",
  "steps": [
    {
      "active": true,
      "replyCallback": {},
      "shouldStopOnFail": false,
      "issues": {},
      "stepId": "send-email",
      "name": "Send email",
      "type": "REGULAR",
      "filters": [],
      "_templateId": {
        "$oid": "6971ab0493e98bfba5435ba5"  // Id of related MessageTemplate
      },
      "_id": {
        "$oid": "6971ab0493e98bfba5435ba5"
      },
      "variants": []
    }
  ]
}
```
5. API Service creates MessageTemplate related to the steps in Workflow. Make sure the `subject` - `content` field is defined properly, `contentType` must be `"html"`
```JSON
{
  "_id": {
    "$oid": "6971ab0493e98bfba5435ba5"
  },
  "type": "email",
  "active": true,
  "name": "Send email",
  "subject": "Confirm Email",
  "contentType": "html",
  "content": "Hi {{subscriber.firstName}}, click here: {{verificationLink}}"
  "_layoutId": {
    "$oid": "6971ab0493e98bfba543555"
  }
}
``` 
*Caution*: In NotificationTemplate 
    - `check` (Optional): If set to true, Novu will attempt to verify the connection during creation. It is recommended to keep it false if you just want to save the settings first.
    - Make sure that `type` is `"email"` 

### Create Notification Workflow (Verification Email Workflow)
1. Define Workflow steps (only 1 step):
    - `type`: Defines the communication channel. In this case, `"email"`. Novu will automatically route this to your active Email Integration (Gmail).
    - `template`: Contains the actual message details:
        - `name`: The indentifier of the step
        - `subject`: The email's subject line. Supports variables like `{{variable}}` (e.g. `"Mã xác nhận của bạn là {{otpCode}}"`).
        - `content`: The body of the email (HTML supported). Variables used here must match the data sent from your Go backend.
        - `layoutIdentifier`: Points to a specific email layout (header/footer). `"default"` uses Novu's standard system layout.
2. Add a header's key `Novu-Environment-Id` using an Environment that belong to the User's Organization (This solve the Critical Bug - Issue of the UI)
3. Create workflow: `POST /v2/workflows`

*Caution*: When passing variables in the email body or subject, you must include the correct namespace. Variables from the trigger call should be prefixed with payload (e.g., `{{payload.variableName}}`), while subscriber data must use the subscriber prefix (e.g., `{{subscriber.firstName}}`). Using naked variables like `{{variableName}}` will result in a validation error.
```HTML
<div style=\"font-family: Arial; padding: 20px;\">
    <h2>Verify Email</h2>
    <p>Hi {{payload.firstName}},</p>
    <a href=\"{{payload.verificationLink}}\" 
        style=\"background: #4F46E5; 
        color: white; 
        padding: 12px 20px; 
        text-decoration: none; 
        border-radius: 5px;\">
        Verify Now
    </a>
</div>"
```

### Trigger the Notification (Trigger the Workflow)
This workflow is usually peform by a third-party (Auth Service)
1. Define Subcriber (Email receiver)'s infomation
```JSON
{
    "subscriberId": "123456789",    // Usually user's id stored in third-party's database
    "email": "receiver@gmail.com",
    "firstName": "Thuận"
}
```
2. Define Email's payload
```JSON
{
    "verificationLink": "https://your-app.com/verify?token=secret-token-123"
}
```
3. Add a header's key `Novu-Environment-Id` using an Environment that belong to the User's Organization
4. Trigger event: `POST /v1/events/trigger`
```JSON
{
    "name": "email-verification",    // Must match workflow_indentifier 
    "to": {subcriber},
    "payload": {payload}
}
```
5. Automatic Subscriber Upsert: Novu follows a "Just-in-Time" provisioning logic. If the subscriberId ("123456789") does not exist in the Novu database, Novu will:
    1. Create a new subscriber record automatically
    2. Store the email and firstName provided in the to object
    3. Proceed to the next step
    
6. Worker Pickup & Template Rendering: Once the job is added to the Redis Queue, the Novu Worker takes over the execution:
    1. Job Acquisition: The Worker service monitors the "standard" queue. It picks up the job created in Step 4 and changes its  tatus from queued to processing.
    2. Logic Execution: The Worker fetches the Workflow definition from the database. It identifies the step `type` (e.g., email) and verifies if any filters or delays need to be applied.
    3. Variable Injection: The Worker retrieves your HTML template and uses a rendering engine (Handlebars) to replace the placeholders with your actual payload data
    
7. Notification & Message Entity Creation: Before sending the actual email, the Worker performs a "Write" operation to the database:
    1. Notification Record: It creates a Notification entity (the one you found in MongoDB) which acts as the "Parent" record for this specific trigger.
    2. Message Record: It creates a Message entity for the specific channel (Email). This record stores the final rendered HTML and is used to track whether the email is sent, delivered, or failed.
    
8. Email Dispatch (Integration Layer): The Worker now acts as a client to your Email Provider:
    1. Provider Handshake: The Worker looks for an Active Integration (e.g., Gmail SMTP). It establishes a connection using the credentials you provided (Host, Port, App Password).
    2. SMTP Transmission: The Worker sends the rendered email to the SMTP server.Final Update: Once the SMTP server accepts the mail, the Worker updates the Message status to sent and completes the job.

### System Delivery Diagnostics
1. Get message detail of a transaction `GET /v1/messages?transactionId={transactionId}`