// Package environment provides parsing and validation for environment configurations.
package environment

import (
	"github.com/architect-io/arcctl/pkg/schema/environment/internal"
)

// Environment represents a parsed and validated environment configuration.
type Environment interface {
	// Metadata
	Name() string
	Datacenter() string

	// Locals
	Locals() map[string]interface{}

	// Components
	Components() map[string]ComponentConfig

	// Version information
	SchemaVersion() string

	// Source information
	SourcePath() string

	// Internal access (for engine use)
	Internal() *internal.InternalEnvironment
}

// ComponentConfig represents a component's configuration within an environment.
// The component key (map key) is the registry address (e.g., ghcr.io/org/my-app).
// Source is either a version tag (e.g., v1.0.0) or a file path (e.g., ./path/to/component).
type ComponentConfig interface {
	// Source returns the version tag (e.g., "v1.0.0") or file path (e.g., "./path/to/component")
	Source() string

	// Variable values
	Variables() map[string]interface{}

	// Scaling configurations per deployment
	Scaling() map[string]ScalingConfig

	// Function configurations per function
	Functions() map[string]FunctionConfig

	// Environment variable overrides per deployment
	Environment() map[string]map[string]string

	// Route configurations per route
	Routes() map[string]RouteConfig
}

// ScalingConfig represents scaling configuration for a deployment.
type ScalingConfig interface {
	Replicas() int
	CPU() string
	Memory() string
	MinReplicas() int
	MaxReplicas() int
}

// FunctionConfig represents configuration for a serverless function.
type FunctionConfig interface {
	Regions() []string
	Memory() string
	Timeout() int
}

// RouteConfig represents route configuration.
type RouteConfig interface {
	Hostnames() []Hostname
	TLS() TLSConfig
}

// Hostname represents a hostname configuration.
type Hostname interface {
	Subdomain() string
	Host() string
}

// TLSConfig represents TLS configuration.
type TLSConfig interface {
	Enabled() bool
	SecretName() string
}

// Loader loads and parses environment configurations.
type Loader interface {
	// Load parses an environment from the given path
	Load(path string) (Environment, error)

	// LoadFromBytes parses an environment from raw bytes
	LoadFromBytes(data []byte, sourcePath string) (Environment, error)

	// Validate validates an environment without fully parsing
	Validate(path string) error
}
