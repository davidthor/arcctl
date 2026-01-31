# config

Configuration management for arcctl. Handles loading, saving, and managing YAML configuration files with support for profiles, registries, state backends, secrets, and more.

## Overview

The `config` package provides a complete configuration management system for the arcctl CLI tool. Configuration files are stored in YAML format at `~/.arcctl/config.yaml` by default.

## Features

- Default datacenter and environment settings
- OCI registry configuration and authentication
- State backend configuration
- Secret provider settings
- Logging configuration
- IaC plugin settings
- Named configuration profiles
- Command aliases
- Environment variable expansion

## Types

### Config

Main configuration struct containing all settings.

```go
type Config struct {
    DefaultDatacenter string
    DefaultEnvironment string
    Registry          RegistryConfig
    State             StateConfig
    Secrets           SecretsConfig
    Logging           LoggingConfig
    Plugins           PluginsConfig
    Profiles          map[string]ProfileConfig
    ActiveProfile     string
    Aliases           map[string]string
}
```

### Manager

Configuration manager for loading and accessing config.

```go
type Manager struct {
    // ...
}
```

### Supporting Types

- **RegistryConfig** - OCI registry settings (default registry, authentication)
- **RegistryAuth** - Registry authentication (username, password, token)
- **StateConfig** - State backend configuration (backend type, config map)
- **SecretsConfig** - Secret provider settings (default provider, named providers)
- **SecretProviderConfig** - Individual secret provider configuration
- **LoggingConfig** - Logging settings (level, format, file, color)
- **PluginsConfig** - IaC plugin settings (default plugin, paths, config)
- **ProfileConfig** - Named profile configuration (datacenter, environment, state, secrets, variables)

## Functions

### Loading Configuration

```go
// Load config from the default location (~/.arcctl/config.yaml)
cfg, err := config.Load()

// Load from a specific file
cfg, err := config.LoadFromFile("/path/to/config.yaml")

// Create a config with defaults
cfg := config.DefaultConfig()
```

### Using the Manager

```go
// Create a new config manager
manager, err := config.NewManager()

// Or with a specific path
manager, err := config.NewManagerWithPath("/path/to/config.yaml")

// Get the current configuration
cfg := manager.Get()

// Update a configuration value using dot notation
err := manager.Set("state.backend", "s3")
err := manager.Set("logging.level", "debug")

// Save changes
err := manager.Save()

// Reload from disk
err := manager.Reload()
```

### Saving Configuration

```go
// Save to the default location
err := cfg.Save()

// Save to a specific file
err := cfg.SaveToFile("/path/to/config.yaml")
```

### Working with Profiles

```go
// Get the active profile configuration
profile := cfg.GetProfile()

// Merge another config into this one
cfg.Merge(otherConfig)
```

### Path Helpers

```go
// Get the default config file path
path, err := config.DefaultConfigPath()  // Returns ~/.arcctl/config.yaml

// Get the default config directory
dir, err := config.DefaultConfigDir()    // Returns ~/.arcctl
```

## Default Values

| Setting | Default Value |
|---------|---------------|
| State backend | `"local"` |
| Secrets provider | `"env"` |
| Logging level | `"info"` |
| Logging format | `"text"` |
| Default plugin | `"native"` |

## File Permissions

- Config file: `0600` (read/write for owner only)
- Config directory: `0755` (created if missing)

## Example Configuration

```yaml
defaultDatacenter: production
defaultEnvironment: staging

registry:
  default: ghcr.io/myorg
  auth:
    ghcr.io:
      username: myuser
      token: ${GITHUB_TOKEN}

state:
  backend: s3
  config:
    bucket: my-arcctl-state
    region: us-west-2

secrets:
  default: vault
  providers:
    vault:
      address: https://vault.example.com
      namespace: arcctl

logging:
  level: info
  format: text
  color: true

plugins:
  default: opentofu
  paths:
    - ~/.arcctl/plugins

profiles:
  production:
    datacenter: prod-dc
    environment: prod
  staging:
    datacenter: staging-dc
    environment: staging

activeProfile: staging
```
