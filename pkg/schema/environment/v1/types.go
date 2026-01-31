// Package v1 implements the v1 environment schema.
package v1

// SchemaV1 represents the v1 environment schema.
type SchemaV1 struct {
	Version    string `yaml:"version,omitempty" json:"version,omitempty"`
	Name       string `yaml:"name,omitempty" json:"name,omitempty"`
	Datacenter string `yaml:"datacenter,omitempty" json:"datacenter,omitempty"`

	// Reusable values
	Locals map[string]interface{} `yaml:"locals,omitempty" json:"locals,omitempty"`

	// Component configurations
	Components map[string]ComponentConfigV1 `yaml:"components,omitempty" json:"components,omitempty"`
}

// ComponentConfigV1 represents a component configuration in v1 schema.
// The component key (map key) is the registry address (e.g., ghcr.io/org/my-app).
// Source is either a version tag (e.g., v1.0.0) or a file path (e.g., ./path/to/component).
type ComponentConfigV1 struct {
	// Source is the version tag (e.g., "v1.0.0") or file path (e.g., "./path/to/component")
	Source string `yaml:"source" json:"source"`

	// Variable values
	Variables map[string]interface{} `yaml:"variables,omitempty" json:"variables,omitempty"`

	// Scaling configuration per deployment
	Scaling map[string]ScalingConfigV1 `yaml:"scaling,omitempty" json:"scaling,omitempty"`

	// Function configuration per function
	Functions map[string]FunctionConfigV1 `yaml:"functions,omitempty" json:"functions,omitempty"`

	// Environment variable overrides per deployment
	Environment map[string]map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`

	// Route configuration per route
	Routes map[string]RouteConfigV1 `yaml:"routes,omitempty" json:"routes,omitempty"`
}

// ScalingConfigV1 represents scaling configuration in v1 schema.
type ScalingConfigV1 struct {
	Replicas    int    `yaml:"replicas,omitempty" json:"replicas,omitempty"`
	CPU         string `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Memory      string `yaml:"memory,omitempty" json:"memory,omitempty"`
	MinReplicas int    `yaml:"min_replicas,omitempty" json:"min_replicas,omitempty"`
	MaxReplicas int    `yaml:"max_replicas,omitempty" json:"max_replicas,omitempty"`
}

// FunctionConfigV1 represents function configuration in v1 schema.
type FunctionConfigV1 struct {
	Regions []string `yaml:"regions,omitempty" json:"regions,omitempty"`
	Memory  string   `yaml:"memory,omitempty" json:"memory,omitempty"`
	Timeout int      `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

// RouteConfigV1 represents route configuration in v1 schema.
type RouteConfigV1 struct {
	Hostnames []HostnameV1  `yaml:"hostnames,omitempty" json:"hostnames,omitempty"`
	TLS       *TLSConfigV1  `yaml:"tls,omitempty" json:"tls,omitempty"`
}

// HostnameV1 represents a hostname in v1 schema.
type HostnameV1 struct {
	Subdomain string `yaml:"subdomain,omitempty" json:"subdomain,omitempty"`
	Host      string `yaml:"host,omitempty" json:"host,omitempty"`
}

// TLSConfigV1 represents TLS configuration in v1 schema.
type TLSConfigV1 struct {
	Enabled    bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	SecretName string `yaml:"secretName,omitempty" json:"secretName,omitempty"`
}
