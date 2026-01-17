# Notification Service

## Run Service
```Bash
docker-compose up -d
```

## Endpoints

### The Dashboard (User Interface)
- URL: `http://localhost:4200`
- Purpose: This is where you (the developer/admin) log in
- Key Actions: 
    - Create Workflows (e.g., "email-verification").
    - Design Templates using a Drag-and-Drop editor or HTML.
    - Configure Integrations (Connect your SendGrid, Twilio, or AWS SES accounts).
    - Monitor Activity Feed (See which emails were sent and if they were opened).
    
### The API Endpoints (For your Identity Service)
- Your Identity Service (Go) will communicate with Novu via the api service at `http://localhost:3000`.
- The most important endpoint you will use is the Trigger endpoint:
| Method | Endpoint | Description | 
| --- | --- | --- |
| POST | `/v1/events/trigger` | Used by your Identity Service to send a notification | 
| GET | `/v1/subscribers` | To manage your user profiles inside Novu | 
| GET | `/v1/workflows` | To fetch the list of available notification workflows | 
| GET | `/v1/notifications/stats` | To get analytics for your dashboard | 

### Health & Internal Endpoints
These are useful for checking if your services are running correctly:
- API Health Check: `http://localhost:3000/v1/health-check`
- Swagger UI (Documentation): Usually available at `http://localhost:3000/api` (depending on the version, it provides a full list of all available REST endpoints).