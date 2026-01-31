# Clerk Authentication Component

A configuration passthrough component for [Clerk.com](https://clerk.com) - a modern authentication platform.

## Overview

Clerk is a SaaS authentication service (no self-hosted option). This component centralizes Clerk configuration across your environment, allowing multiple applications to share consistent auth settings.

### Why Use This Pattern?

1. **Single Source of Truth**: Configure Clerk credentials once, use everywhere
2. **Environment Separation**: Different Clerk instances for staging vs production
3. **Dependency Visibility**: Clearly shows which apps use Clerk authentication
4. **Consistent Configuration**: Ensures all apps use the same Clerk instance

## Configuration

### Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `domain` | Yes | Your Clerk frontend API domain (e.g., `my-app.clerk.accounts.dev`) |
| `publishable_key` | Yes | Clerk publishable key (`pk_test_...` or `pk_live_...`) |
| `secret_key` | Yes | Clerk secret key (`sk_test_...` or `sk_live_...`) - **sensitive** |
| `webhook_secret` | No | Webhook signing secret - **sensitive** |

### Outputs

| Output | Description |
|--------|-------------|
| `domain` | Clerk frontend API domain |
| `publishable_key` | Publishable key for frontend SDKs |
| `secret_key` | Secret key for backend verification (**sensitive**) |
| `webhook_secret` | Webhook signing secret (**sensitive**) |

## Usage

### 1. Environment Configuration

Configure Clerk in your `environment.yml`:

```yaml
# environment.yml
name: production
datacenter: ghcr.io/myorg/aws-datacenter:v1

components:
  clerk:
    component: ghcr.io/myorg/clerk:v1
    variables:
      domain: my-app.clerk.accounts.dev
      publishable_key: pk_live_xxxxxxxxxxxxx
      secret_key: ${{ secrets.clerk_secret_key }}
      webhook_secret: ${{ secrets.clerk_webhook_secret }}

  my-app:
    component: ghcr.io/myorg/my-app:v1
```

### 2. Dependent Application

In your application's `architect.yml`, declare Clerk as a dependency and access outputs:

```yaml
# my-app/architect.yml

dependencies:
  clerk:
    component: ghcr.io/myorg/clerk:v1

deployments:
  api:
    build:
      context: ./api
    environment:
      CLERK_DOMAIN: ${{ dependencies.clerk.outputs.domain }}
      CLERK_SECRET_KEY: ${{ dependencies.clerk.outputs.secret_key }}
      CLERK_WEBHOOK_SECRET: ${{ dependencies.clerk.outputs.webhook_secret }}

functions:
  web:
    build:
      context: ./web
    framework: nextjs
    environment:
      NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY: ${{ dependencies.clerk.outputs.publishable_key }}
      # App-specific settings configured here, not in the Clerk component
      NEXT_PUBLIC_CLERK_SIGN_IN_URL: /sign-in
      NEXT_PUBLIC_CLERK_SIGN_UP_URL: /sign-up

services:
  api:
    deployment: api
    port: 3000

  web:
    function: web

routes:
  main:
    type: http
    rules:
      - name: api
        matches:
          - path:
              type: PathPrefix
              value: /api
        backendRefs:
          - service: api
            port: 3000

      - name: web
        matches:
          - path:
              type: PathPrefix
              value: /
        backendRefs:
          - service: web
```

## Clerk Endpoints Reference

| Endpoint | URL Pattern |
|----------|-------------|
| Frontend API | `https://{domain}` |
| Backend API | `https://api.clerk.com` |
| JWKS | `https://{domain}/.well-known/jwks.json` |
| OpenID Config | `https://{domain}/.well-known/openid-configuration` |

## Design Notes

### Configuration Passthrough Pattern

This component demonstrates the "configuration passthrough" pattern:

1. No deployments, services, or infrastructure
2. Variables capture Clerk credentials
3. Outputs expose those values to dependent components
4. Apps access via `${{ dependencies.clerk.outputs.<field> }}`

### Test vs Live Mode

Use separate environments for test and live Clerk keys:

```yaml
# environments/staging/environment.yml
components:
  clerk:
    variables:
      domain: my-app-staging.clerk.accounts.dev
      publishable_key: pk_test_xxx
      secret_key: ${{ secrets.clerk_secret_key_test }}

# environments/production/environment.yml
components:
  clerk:
    variables:
      domain: my-app.clerk.accounts.dev
      publishable_key: pk_live_xxx
      secret_key: ${{ secrets.clerk_secret_key_live }}
```
