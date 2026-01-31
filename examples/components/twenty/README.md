# Twenty Component

This component deploys [Twenty](https://twenty.com/), an open-source CRM platform that gives you full control over your customer relationship data.

## Overview

Twenty provides:
- Modern, intuitive CRM interface
- Full data ownership and compliance capabilities
- Gmail and Google Calendar integration
- Microsoft 365 integration
- Custom objects and fields
- Workflow automation
- REST and GraphQL APIs

## Architecture

This component deploys the following services:

| Service | Image | Description |
|---------|-------|-------------|
| `server` | twentycrm/twenty | Main application server (handles API, UI, and migrations) |
| `worker` | twentycrm/twenty | Background worker for async jobs |

Plus these databases:
- **PostgreSQL 16** - Primary data storage
- **Redis 7** - Job queue and caching

## System Requirements

- **RAM**: Minimum 2GB recommended
- The server runs database migrations on startup
- The worker processes background jobs asynchronously

## Required Variables

Before deploying, you must provide these required values:

### Server URL

```bash
# The public URL where Twenty will be accessible
server_url: "https://crm.example.com"
```

### Security Keys

```bash
# Generate APP_SECRET
openssl rand -base64 32
```

## Example Environment Configuration

### Basic Setup

```yaml
# environment.yml
name: twenty-production
datacenter: aws-ecs

components:
  twenty:
    source: ./twenty
    variables:
      server_url: "https://crm.example.com"
      app_secret: "${APP_SECRET}"
```

### With Google Integration

```yaml
# environment.yml
name: twenty-production
datacenter: aws-ecs

components:
  twenty:
    source: ./twenty
    variables:
      server_url: "https://crm.example.com"
      app_secret: "${APP_SECRET}"
      # Google OAuth
      auth_google_enabled: "true"
      auth_google_client_id: "${GOOGLE_CLIENT_ID}"
      auth_google_client_secret: "${GOOGLE_CLIENT_SECRET}"
      # Gmail integration
      messaging_provider_gmail_enabled: "true"
      # Google Calendar
      calendar_provider_google_enabled: "true"
```

### With Email Configuration

```yaml
# environment.yml
name: twenty-production
datacenter: aws-ecs

components:
  twenty:
    source: ./twenty
    variables:
      server_url: "https://crm.example.com"
      app_secret: "${APP_SECRET}"
      # SMTP configuration
      email_driver: "smtp"
      email_from_address: "crm@example.com"
      email_from_name: "Example CRM"
      email_system_address: "system@example.com"
      email_smtp_host: "smtp.sendgrid.net"
      email_smtp_port: "587"
      email_smtp_user: "apikey"
      email_smtp_password: "${SENDGRID_API_KEY}"
```

### Multi-Workspace Mode (SaaS)

```yaml
# environment.yml
name: twenty-saas
datacenter: aws-ecs

components:
  twenty:
    source: ./twenty
    variables:
      server_url: "https://crm.example.com"
      app_secret: "${APP_SECRET}"
      is_multiworkspace_enabled: "true"
      default_subdomain: "app"
```

**Note**: Multi-workspace mode requires wildcard DNS configuration (`*.crm.example.com`).

## Optional Features

### Google Integration

To enable Google SSO and integrations:

1. Create a project in [Google Cloud Console](https://console.cloud.google.com/)
2. Enable Gmail API, Google Calendar API, and People API
3. Create OAuth 2.0 credentials with these redirect URIs:
   - `https://{your-domain}/auth/google/redirect` (for SSO)
   - `https://{your-domain}/auth/google-apis/get-access-token` (for integrations)

| Variable | Description |
|----------|-------------|
| `auth_google_enabled` | Enable Google SSO (`true`/`false`) |
| `auth_google_client_id` | Google OAuth client ID |
| `auth_google_client_secret` | Google OAuth client secret |
| `messaging_provider_gmail_enabled` | Enable Gmail integration |
| `calendar_provider_google_enabled` | Enable Google Calendar integration |

### Microsoft 365 Integration

To enable Microsoft SSO and integrations:

1. Create an app in [Microsoft Azure](https://portal.azure.com/)
2. Enable required Microsoft Graph permissions
3. Add redirect URIs:
   - `https://{your-domain}/auth/microsoft/redirect`
   - `https://{your-domain}/auth/microsoft-apis/get-access-token`

| Variable | Description |
|----------|-------------|
| `auth_microsoft_enabled` | Enable Microsoft SSO (`true`/`false`) |
| `auth_microsoft_client_id` | Microsoft OAuth client ID |
| `auth_microsoft_client_secret` | Microsoft OAuth client secret |
| `messaging_provider_microsoft_enabled` | Enable Microsoft 365 email integration |
| `calendar_provider_microsoft_enabled` | Enable Microsoft Calendar integration |

### Email Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `email_driver` | `logger` | Email driver (`smtp` or `logger`) |
| `email_from_address` | - | Sender email address |
| `email_from_name` | `Twenty` | Sender display name |
| `email_system_address` | - | System email address |
| `email_smtp_host` | - | SMTP server hostname |
| `email_smtp_port` | `587` | SMTP server port |
| `email_smtp_user` | - | SMTP username |
| `email_smtp_password` | - | SMTP password |

### Storage Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `storage_type` | `local` | Storage driver (`local` or `s3`) |
| `storage_s3_region` | - | S3 region |
| `storage_s3_name` | - | S3 bucket name |
| `storage_s3_endpoint` | - | S3 endpoint (for S3-compatible storage) |

### Logic Functions (Serverless)

| Variable | Default | Description |
|----------|---------|-------------|
| `serverless_type` | `LOCAL` | Execution driver (`LOCAL`, `LAMBDA`, or `DISABLED`) |

**Security Note**: The `LOCAL` driver runs code without sandboxing. Use `LAMBDA` or `DISABLED` for production with untrusted code.

## Getting Started

After deployment:

1. Access Twenty at your configured `server_url`
2. Create your first workspace and admin account
3. Invite team members
4. Configure integrations in **Settings → Admin Panel → Configuration Variables**

## SSL Requirements

SSL (HTTPS) is strongly recommended for production deployments. Some browser features (like clipboard API) require a secure context to work properly.

## Documentation

- [Twenty Documentation](https://docs.twenty.com/)
- [Self-Hosting Guide](https://docs.twenty.com/developers/self-host/self-host)
- [Setup Guide](https://docs.twenty.com/developers/self-host/capabilities/setup)
- [API Reference](https://docs.twenty.com/developers/extend/capabilities/rest-api)
