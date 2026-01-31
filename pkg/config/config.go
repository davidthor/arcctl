// Package config provides configuration file support for arcctl.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the arcctl configuration.
type Config struct {
	// DefaultDatacenter is the default datacenter to use
	DefaultDatacenter string `yaml:"default_datacenter,omitempty"`

	// DefaultEnvironment is the default environment to use
	DefaultEnvironment string `yaml:"default_environment,omitempty"`

	// Registry configures OCI registry settings
	Registry RegistryConfig `yaml:"registry,omitempty"`

	// State configures state backend settings
	State StateConfig `yaml:"state,omitempty"`

	// Secrets configures secret provider settings
	Secrets SecretsConfig `yaml:"secrets,omitempty"`

	// Logging configures logging settings
	Logging LoggingConfig `yaml:"logging,omitempty"`

	// Plugins configures IaC plugin settings
	Plugins PluginsConfig `yaml:"plugins,omitempty"`

	// Profiles defines named configuration profiles
	Profiles map[string]ProfileConfig `yaml:"profiles,omitempty"`

	// ActiveProfile is the currently active profile
	ActiveProfile string `yaml:"active_profile,omitempty"`

	// Aliases defines command aliases
	Aliases map[string]string `yaml:"aliases,omitempty"`
}

// RegistryConfig configures OCI registry settings.
type RegistryConfig struct {
	// Default is the default registry to push/pull from
	Default string `yaml:"default,omitempty"`

	// Auth configures registry authentication
	Auth map[string]RegistryAuth `yaml:"auth,omitempty"`
}

// RegistryAuth contains registry authentication info.
type RegistryAuth struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"` // Can be a secret reference
	Token    string `yaml:"token,omitempty"`    // Can be a secret reference
}

// StateConfig configures state backend settings.
type StateConfig struct {
	// Backend is the state backend type (local, s3, gcs)
	Backend string `yaml:"backend,omitempty"`

	// Config contains backend-specific configuration
	Config map[string]string `yaml:"config,omitempty"`
}

// SecretsConfig configures secret provider settings.
type SecretsConfig struct {
	// Provider is the default secret provider
	Provider string `yaml:"provider,omitempty"`

	// Providers configures named secret providers
	Providers map[string]SecretProviderConfig `yaml:"providers,omitempty"`
}

// SecretProviderConfig configures a secret provider.
type SecretProviderConfig struct {
	Type   string            `yaml:"type"`
	Config map[string]string `yaml:"config,omitempty"`
}

// LoggingConfig configures logging settings.
type LoggingConfig struct {
	// Level is the log level (debug, info, warn, error)
	Level string `yaml:"level,omitempty"`

	// Format is the log format (text, json)
	Format string `yaml:"format,omitempty"`

	// File is an optional log file path
	File string `yaml:"file,omitempty"`

	// Color enables/disables colored output
	Color *bool `yaml:"color,omitempty"`
}

// PluginsConfig configures IaC plugin settings.
type PluginsConfig struct {
	// Default is the default IaC plugin to use
	Default string `yaml:"default,omitempty"`

	// Paths contains additional paths to search for plugins
	Paths []string `yaml:"paths,omitempty"`

	// Config contains plugin-specific configuration
	Config map[string]map[string]string `yaml:"config,omitempty"`
}

// ProfileConfig defines a configuration profile.
type ProfileConfig struct {
	Datacenter  string            `yaml:"datacenter,omitempty"`
	Environment string            `yaml:"environment,omitempty"`
	State       StateConfig       `yaml:"state,omitempty"`
	Secrets     SecretsConfig     `yaml:"secrets,omitempty"`
	Variables   map[string]string `yaml:"variables,omitempty"`
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		State: StateConfig{
			Backend: "local",
			Config:  map[string]string{},
		},
		Secrets: SecretsConfig{
			Provider: "env",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
		Plugins: PluginsConfig{
			Default: "native",
		},
		Profiles: map[string]ProfileConfig{},
		Aliases:  map[string]string{},
	}
}

// Load loads the configuration from the default location.
func Load() (*Config, error) {
	configPath, err := DefaultConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	return LoadFromFile(configPath)
}

// LoadFromFile loads configuration from a specific file.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Expand environment variables in values
	expandEnvVars(config)

	return config, nil
}

// Save saves the configuration to the default location.
func (c *Config) Save() error {
	configPath, err := DefaultConfigPath()
	if err != nil {
		return err
	}

	return c.SaveToFile(configPath)
}

// SaveToFile saves the configuration to a specific file.
func (c *Config) SaveToFile(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetProfile returns the active profile configuration.
func (c *Config) GetProfile() *ProfileConfig {
	if c.ActiveProfile == "" {
		return nil
	}

	profile, ok := c.Profiles[c.ActiveProfile]
	if !ok {
		return nil
	}

	return &profile
}

// Merge merges another config into this one.
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	if other.DefaultDatacenter != "" {
		c.DefaultDatacenter = other.DefaultDatacenter
	}
	if other.DefaultEnvironment != "" {
		c.DefaultEnvironment = other.DefaultEnvironment
	}
	if other.ActiveProfile != "" {
		c.ActiveProfile = other.ActiveProfile
	}

	// Merge registry config
	if other.Registry.Default != "" {
		c.Registry.Default = other.Registry.Default
	}
	if other.Registry.Auth != nil {
		if c.Registry.Auth == nil {
			c.Registry.Auth = make(map[string]RegistryAuth)
		}
		for k, v := range other.Registry.Auth {
			c.Registry.Auth[k] = v
		}
	}

	// Merge state config
	if other.State.Backend != "" {
		c.State.Backend = other.State.Backend
	}
	if other.State.Config != nil {
		if c.State.Config == nil {
			c.State.Config = make(map[string]string)
		}
		for k, v := range other.State.Config {
			c.State.Config[k] = v
		}
	}

	// Merge secrets config
	if other.Secrets.Provider != "" {
		c.Secrets.Provider = other.Secrets.Provider
	}
	if other.Secrets.Providers != nil {
		if c.Secrets.Providers == nil {
			c.Secrets.Providers = make(map[string]SecretProviderConfig)
		}
		for k, v := range other.Secrets.Providers {
			c.Secrets.Providers[k] = v
		}
	}

	// Merge logging config
	if other.Logging.Level != "" {
		c.Logging.Level = other.Logging.Level
	}
	if other.Logging.Format != "" {
		c.Logging.Format = other.Logging.Format
	}
	if other.Logging.File != "" {
		c.Logging.File = other.Logging.File
	}
	if other.Logging.Color != nil {
		c.Logging.Color = other.Logging.Color
	}

	// Merge plugins config
	if other.Plugins.Default != "" {
		c.Plugins.Default = other.Plugins.Default
	}
	if other.Plugins.Paths != nil {
		c.Plugins.Paths = append(c.Plugins.Paths, other.Plugins.Paths...)
	}
	if other.Plugins.Config != nil {
		if c.Plugins.Config == nil {
			c.Plugins.Config = make(map[string]map[string]string)
		}
		for k, v := range other.Plugins.Config {
			c.Plugins.Config[k] = v
		}
	}

	// Merge profiles
	if other.Profiles != nil {
		if c.Profiles == nil {
			c.Profiles = make(map[string]ProfileConfig)
		}
		for k, v := range other.Profiles {
			c.Profiles[k] = v
		}
	}

	// Merge aliases
	if other.Aliases != nil {
		if c.Aliases == nil {
			c.Aliases = make(map[string]string)
		}
		for k, v := range other.Aliases {
			c.Aliases[k] = v
		}
	}
}

// DefaultConfigPath returns the default configuration file path.
func DefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".arcctl", "config.yaml"), nil
}

// DefaultConfigDir returns the default configuration directory.
func DefaultConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".arcctl"), nil
}

// expandEnvVars expands environment variables in the config.
func expandEnvVars(config *Config) {
	// Expand in registry auth
	for name, auth := range config.Registry.Auth {
		auth.Username = os.ExpandEnv(auth.Username)
		auth.Password = os.ExpandEnv(auth.Password)
		auth.Token = os.ExpandEnv(auth.Token)
		config.Registry.Auth[name] = auth
	}

	// Expand in state config
	for k, v := range config.State.Config {
		config.State.Config[k] = os.ExpandEnv(v)
	}

	// Expand in secrets providers
	for name, provider := range config.Secrets.Providers {
		for k, v := range provider.Config {
			provider.Config[k] = os.ExpandEnv(v)
		}
		config.Secrets.Providers[name] = provider
	}

	// Expand in plugin config
	for name, pluginCfg := range config.Plugins.Config {
		for k, v := range pluginCfg {
			config.Plugins.Config[name][k] = os.ExpandEnv(v)
		}
	}

	// Expand profile variables
	for name, profile := range config.Profiles {
		for k, v := range profile.Variables {
			profile.Variables[k] = os.ExpandEnv(v)
		}
		config.Profiles[name] = profile
	}
}

// Manager manages configuration loading and access.
type Manager struct {
	config     *Config
	configPath string
}

// NewManager creates a new configuration manager.
func NewManager() (*Manager, error) {
	configPath, err := DefaultConfigPath()
	if err != nil {
		return nil, err
	}

	config, err := LoadFromFile(configPath)
	if err != nil {
		return nil, err
	}

	return &Manager{
		config:     config,
		configPath: configPath,
	}, nil
}

// NewManagerWithPath creates a manager with a specific config path.
func NewManagerWithPath(configPath string) (*Manager, error) {
	config, err := LoadFromFile(configPath)
	if err != nil {
		return nil, err
	}

	return &Manager{
		config:     config,
		configPath: configPath,
	}, nil
}

// Get returns the current configuration.
func (m *Manager) Get() *Config {
	return m.config
}

// Set updates the configuration.
func (m *Manager) Set(key, value string) error {
	parts := strings.Split(key, ".")

	switch parts[0] {
	case "default_datacenter":
		m.config.DefaultDatacenter = value
	case "default_environment":
		m.config.DefaultEnvironment = value
	case "active_profile":
		m.config.ActiveProfile = value
	case "state":
		if len(parts) > 1 {
			switch parts[1] {
			case "backend":
				m.config.State.Backend = value
			default:
				if m.config.State.Config == nil {
					m.config.State.Config = make(map[string]string)
				}
				m.config.State.Config[parts[1]] = value
			}
		}
	case "logging":
		if len(parts) > 1 {
			switch parts[1] {
			case "level":
				m.config.Logging.Level = value
			case "format":
				m.config.Logging.Format = value
			case "file":
				m.config.Logging.File = value
			}
		}
	case "plugins":
		if len(parts) > 1 && parts[1] == "default" {
			m.config.Plugins.Default = value
		}
	case "registry":
		if len(parts) > 1 && parts[1] == "default" {
			m.config.Registry.Default = value
		}
	case "secrets":
		if len(parts) > 1 && parts[1] == "provider" {
			m.config.Secrets.Provider = value
		}
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	return m.Save()
}

// Save saves the current configuration.
func (m *Manager) Save() error {
	return m.config.SaveToFile(m.configPath)
}

// Reload reloads the configuration from disk.
func (m *Manager) Reload() error {
	config, err := LoadFromFile(m.configPath)
	if err != nil {
		return err
	}

	m.config = config
	return nil
}
