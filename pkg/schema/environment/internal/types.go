// Package internal contains the canonical internal representation for environments.
package internal

// InternalEnvironment is the canonical internal representation for environment configurations.
// Note: Name and Datacenter are not part of the config - they are provided via CLI
// when creating/updating an environment.
type InternalEnvironment struct {
	// Environment-level variable declarations
	Variables map[string]InternalEnvironmentVariable

	// Reusable values
	Locals map[string]interface{}

	// Component configurations
	Components map[string]InternalComponentConfig

	// Source information
	SourceVersion string
	SourcePath    string
}

// InternalEnvironmentVariable represents an environment-level variable declaration.
// Variables are resolved from OS environment variables, dotenv files, or defaults.
type InternalEnvironmentVariable struct {
	Name        string
	Description string
	Default     interface{}
	Required    bool
	Sensitive   bool
	Env         string // Explicit OS env var name override (defaults to UPPER_SNAKE_CASE of Name)
}

// InternalComponentConfig represents the configuration for a component in an environment.
// The component key (map key) is the registry address (e.g., ghcr.io/org/my-app).
// Source is either a version tag (e.g., v1.0.0) or a file path (e.g., ./path/to/component).
type InternalComponentConfig struct {
	// Source is the version tag (e.g., "v1.0.0") or file path (e.g., "./path/to/component")
	Source string

	// Variable values for the component
	Variables map[string]interface{}

	// Scaling configuration per deployment
	Scaling map[string]InternalScalingConfig

	// Function configuration per function
	Functions map[string]InternalFunctionConfig

	// Environment variable overrides per deployment
	Environment map[string]map[string]string

	// Route configuration per route
	Routes map[string]InternalRouteConfig
}

// InternalScalingConfig represents scaling configuration for a deployment.
type InternalScalingConfig struct {
	Replicas    int
	CPU         string
	Memory      string
	MinReplicas int
	MaxReplicas int
}

// InternalFunctionConfig represents configuration for a serverless function.
type InternalFunctionConfig struct {
	Regions []string
	Memory  string
	Timeout int
}

// InternalRouteConfig represents route configuration in an environment.
type InternalRouteConfig struct {
	Hostnames []InternalHostname
	TLS       *InternalTLSConfig
}

// InternalHostname represents a hostname configuration.
type InternalHostname struct {
	// One of these is set
	Subdomain string // Results in subdomain.<env-domain>
	Host      string // Explicit full hostname
}

// InternalTLSConfig represents TLS configuration.
type InternalTLSConfig struct {
	Enabled    bool
	SecretName string
}
