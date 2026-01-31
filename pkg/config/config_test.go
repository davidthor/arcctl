package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.State.Backend != "local" {
		t.Errorf("expected state backend 'local', got %q", cfg.State.Backend)
	}

	if cfg.Secrets.Provider != "env" {
		t.Errorf("expected secrets provider 'env', got %q", cfg.Secrets.Provider)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("expected logging level 'info', got %q", cfg.Logging.Level)
	}

	if cfg.Logging.Format != "text" {
		t.Errorf("expected logging format 'text', got %q", cfg.Logging.Format)
	}

	if cfg.Plugins.Default != "native" {
		t.Errorf("expected plugins default 'native', got %q", cfg.Plugins.Default)
	}

	if cfg.Profiles == nil {
		t.Error("expected profiles map to be initialized")
	}

	if cfg.Aliases == nil {
		t.Error("expected aliases map to be initialized")
	}
}

func TestLoadFromFile_NonExistent(t *testing.T) {
	cfg, err := LoadFromFile("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("unexpected error for non-existent file: %v", err)
	}

	// Should return default config
	if cfg.State.Backend != "local" {
		t.Errorf("expected default state backend 'local', got %q", cfg.State.Backend)
	}
}

func TestLoadFromFile_ValidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configData := `
default_datacenter: my-dc
default_environment: staging
state:
  backend: s3
  config:
    bucket: my-bucket
    region: us-east-1
logging:
  level: debug
  format: json
plugins:
  default: opentofu
profiles:
  prod:
    datacenter: prod-dc
    environment: production
    variables:
      API_URL: https://api.example.com
aliases:
  d: datacenter
  e: environment
`

	if err := os.WriteFile(configPath, []byte(configData), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.DefaultDatacenter != "my-dc" {
		t.Errorf("expected default_datacenter 'my-dc', got %q", cfg.DefaultDatacenter)
	}

	if cfg.DefaultEnvironment != "staging" {
		t.Errorf("expected default_environment 'staging', got %q", cfg.DefaultEnvironment)
	}

	if cfg.State.Backend != "s3" {
		t.Errorf("expected state backend 's3', got %q", cfg.State.Backend)
	}

	if cfg.State.Config["bucket"] != "my-bucket" {
		t.Errorf("expected state config bucket 'my-bucket', got %q", cfg.State.Config["bucket"])
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("expected logging level 'debug', got %q", cfg.Logging.Level)
	}

	if cfg.Plugins.Default != "opentofu" {
		t.Errorf("expected plugins default 'opentofu', got %q", cfg.Plugins.Default)
	}

	if cfg.Aliases["d"] != "datacenter" {
		t.Errorf("expected alias 'd' to be 'datacenter', got %q", cfg.Aliases["d"])
	}

	// Check profile
	profile, ok := cfg.Profiles["prod"]
	if !ok {
		t.Fatal("expected 'prod' profile to exist")
	}
	if profile.Datacenter != "prod-dc" {
		t.Errorf("expected profile datacenter 'prod-dc', got %q", profile.Datacenter)
	}
}

func TestLoadFromFile_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Use truly invalid YAML - tabs mixed with spaces in a way that breaks parsing
	invalidYAML := `default_datacenter: [unclosed bracket`

	if err := os.WriteFile(configPath, []byte(invalidYAML), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := LoadFromFile(configPath)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadFromFile_EnvVarExpansion(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Set environment variables
	os.Setenv("TEST_USERNAME", "testuser")
	os.Setenv("TEST_PASSWORD", "secret123")
	os.Setenv("TEST_BUCKET", "env-bucket")
	defer os.Unsetenv("TEST_USERNAME")
	defer os.Unsetenv("TEST_PASSWORD")
	defer os.Unsetenv("TEST_BUCKET")

	configData := `
registry:
  default: docker.io
  auth:
    docker.io:
      username: $TEST_USERNAME
      password: ${TEST_PASSWORD}
state:
  backend: s3
  config:
    bucket: ${TEST_BUCKET}
`

	if err := os.WriteFile(configPath, []byte(configData), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	auth, ok := cfg.Registry.Auth["docker.io"]
	if !ok {
		t.Fatal("expected docker.io auth to exist")
	}

	if auth.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %q", auth.Username)
	}

	if auth.Password != "secret123" {
		t.Errorf("expected password 'secret123', got %q", auth.Password)
	}

	if cfg.State.Config["bucket"] != "env-bucket" {
		t.Errorf("expected bucket 'env-bucket', got %q", cfg.State.Config["bucket"])
	}
}

func TestConfig_SaveToFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.yaml")

	cfg := DefaultConfig()
	cfg.DefaultDatacenter = "test-dc"
	cfg.DefaultEnvironment = "test-env"

	if err := cfg.SaveToFile(configPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("expected config file to exist")
	}

	// Load and verify
	loaded, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("failed to load saved config: %v", err)
	}

	if loaded.DefaultDatacenter != "test-dc" {
		t.Errorf("expected default_datacenter 'test-dc', got %q", loaded.DefaultDatacenter)
	}

	if loaded.DefaultEnvironment != "test-env" {
		t.Errorf("expected default_environment 'test-env', got %q", loaded.DefaultEnvironment)
	}
}

func TestConfig_GetProfile(t *testing.T) {
	tests := []struct {
		name           string
		activeProfile  string
		profiles       map[string]ProfileConfig
		expectNil      bool
		expectDC       string
	}{
		{
			name:          "no active profile",
			activeProfile: "",
			profiles:      map[string]ProfileConfig{},
			expectNil:     true,
		},
		{
			name:          "active profile not found",
			activeProfile: "nonexistent",
			profiles:      map[string]ProfileConfig{},
			expectNil:     true,
		},
		{
			name:          "active profile found",
			activeProfile: "prod",
			profiles: map[string]ProfileConfig{
				"prod": {
					Datacenter:  "prod-dc",
					Environment: "production",
				},
			},
			expectNil: false,
			expectDC:  "prod-dc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				ActiveProfile: tt.activeProfile,
				Profiles:      tt.profiles,
			}

			profile := cfg.GetProfile()

			if tt.expectNil {
				if profile != nil {
					t.Error("expected nil profile")
				}
			} else {
				if profile == nil {
					t.Fatal("expected non-nil profile")
				}
				if profile.Datacenter != tt.expectDC {
					t.Errorf("expected datacenter %q, got %q", tt.expectDC, profile.Datacenter)
				}
			}
		})
	}
}

func TestConfig_Merge(t *testing.T) {
	t.Run("merge nil config", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Merge(nil) // Should not panic
	})

	t.Run("merge basic fields", func(t *testing.T) {
		cfg := DefaultConfig()
		other := &Config{
			DefaultDatacenter:  "other-dc",
			DefaultEnvironment: "other-env",
			ActiveProfile:      "prod",
		}

		cfg.Merge(other)

		if cfg.DefaultDatacenter != "other-dc" {
			t.Errorf("expected default_datacenter 'other-dc', got %q", cfg.DefaultDatacenter)
		}
		if cfg.DefaultEnvironment != "other-env" {
			t.Errorf("expected default_environment 'other-env', got %q", cfg.DefaultEnvironment)
		}
		if cfg.ActiveProfile != "prod" {
			t.Errorf("expected active_profile 'prod', got %q", cfg.ActiveProfile)
		}
	})

	t.Run("merge registry config", func(t *testing.T) {
		cfg := DefaultConfig()
		other := &Config{
			Registry: RegistryConfig{
				Default: "gcr.io",
				Auth: map[string]RegistryAuth{
					"gcr.io": {Username: "user", Password: "pass"},
				},
			},
		}

		cfg.Merge(other)

		if cfg.Registry.Default != "gcr.io" {
			t.Errorf("expected registry default 'gcr.io', got %q", cfg.Registry.Default)
		}
		if cfg.Registry.Auth["gcr.io"].Username != "user" {
			t.Error("expected registry auth to be merged")
		}
	})

	t.Run("merge state config", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.State.Config["existing"] = "value"

		other := &Config{
			State: StateConfig{
				Backend: "s3",
				Config: map[string]string{
					"bucket": "my-bucket",
				},
			},
		}

		cfg.Merge(other)

		if cfg.State.Backend != "s3" {
			t.Errorf("expected state backend 's3', got %q", cfg.State.Backend)
		}
		if cfg.State.Config["bucket"] != "my-bucket" {
			t.Error("expected bucket config to be merged")
		}
		if cfg.State.Config["existing"] != "value" {
			t.Error("expected existing config to be preserved")
		}
	})

	t.Run("merge logging config", func(t *testing.T) {
		cfg := DefaultConfig()
		colorTrue := true
		other := &Config{
			Logging: LoggingConfig{
				Level:  "debug",
				Format: "json",
				File:   "/var/log/arcctl.log",
				Color:  &colorTrue,
			},
		}

		cfg.Merge(other)

		if cfg.Logging.Level != "debug" {
			t.Errorf("expected logging level 'debug', got %q", cfg.Logging.Level)
		}
		if cfg.Logging.Format != "json" {
			t.Errorf("expected logging format 'json', got %q", cfg.Logging.Format)
		}
		if cfg.Logging.File != "/var/log/arcctl.log" {
			t.Errorf("expected logging file '/var/log/arcctl.log', got %q", cfg.Logging.File)
		}
		if cfg.Logging.Color == nil || *cfg.Logging.Color != true {
			t.Error("expected logging color to be true")
		}
	})

	t.Run("merge plugins config", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Plugins.Paths = []string{"/existing/path"}

		other := &Config{
			Plugins: PluginsConfig{
				Default: "opentofu",
				Paths:   []string{"/new/path"},
				Config: map[string]map[string]string{
					"opentofu": {"version": "1.0"},
				},
			},
		}

		cfg.Merge(other)

		if cfg.Plugins.Default != "opentofu" {
			t.Errorf("expected plugins default 'opentofu', got %q", cfg.Plugins.Default)
		}
		if len(cfg.Plugins.Paths) != 2 {
			t.Errorf("expected 2 plugin paths, got %d", len(cfg.Plugins.Paths))
		}
		if cfg.Plugins.Config["opentofu"]["version"] != "1.0" {
			t.Error("expected plugin config to be merged")
		}
	})

	t.Run("merge profiles", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Profiles["dev"] = ProfileConfig{Datacenter: "dev-dc"}

		other := &Config{
			Profiles: map[string]ProfileConfig{
				"prod": {Datacenter: "prod-dc"},
			},
		}

		cfg.Merge(other)

		if cfg.Profiles["dev"].Datacenter != "dev-dc" {
			t.Error("expected existing profile to be preserved")
		}
		if cfg.Profiles["prod"].Datacenter != "prod-dc" {
			t.Error("expected new profile to be merged")
		}
	})

	t.Run("merge aliases", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.Aliases["d"] = "datacenter"

		other := &Config{
			Aliases: map[string]string{
				"e": "environment",
			},
		}

		cfg.Merge(other)

		if cfg.Aliases["d"] != "datacenter" {
			t.Error("expected existing alias to be preserved")
		}
		if cfg.Aliases["e"] != "environment" {
			t.Error("expected new alias to be merged")
		}
	})

	t.Run("merge secrets config", func(t *testing.T) {
		cfg := DefaultConfig()
		other := &Config{
			Secrets: SecretsConfig{
				Provider: "vault",
				Providers: map[string]SecretProviderConfig{
					"vault": {
						Type:   "vault",
						Config: map[string]string{"addr": "https://vault.example.com"},
					},
				},
			},
		}

		cfg.Merge(other)

		if cfg.Secrets.Provider != "vault" {
			t.Errorf("expected secrets provider 'vault', got %q", cfg.Secrets.Provider)
		}
		if cfg.Secrets.Providers["vault"].Type != "vault" {
			t.Error("expected vault provider to be merged")
		}
	})
}

func TestDefaultConfigPath(t *testing.T) {
	path, err := DefaultConfigPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path, got %q", path)
	}

	if filepath.Base(path) != "config.yaml" {
		t.Errorf("expected config.yaml, got %q", filepath.Base(path))
	}

	if filepath.Base(filepath.Dir(path)) != ".arcctl" {
		t.Errorf("expected parent dir .arcctl, got %q", filepath.Base(filepath.Dir(path)))
	}
}

func TestDefaultConfigDir(t *testing.T) {
	dir, err := DefaultConfigDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !filepath.IsAbs(dir) {
		t.Errorf("expected absolute path, got %q", dir)
	}

	if filepath.Base(dir) != ".arcctl" {
		t.Errorf("expected .arcctl, got %q", filepath.Base(dir))
	}
}

func TestNewManagerWithPath(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configData := `
default_datacenter: test-dc
logging:
  level: debug
`
	if err := os.WriteFile(configPath, []byte(configData), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	manager, err := NewManagerWithPath(configPath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	cfg := manager.Get()
	if cfg.DefaultDatacenter != "test-dc" {
		t.Errorf("expected default_datacenter 'test-dc', got %q", cfg.DefaultDatacenter)
	}
}

func TestManager_Set(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create empty config
	if err := os.WriteFile(configPath, []byte(""), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	manager, err := NewManagerWithPath(configPath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	tests := []struct {
		key      string
		value    string
		validate func(*Config) bool
	}{
		{
			key:   "default_datacenter",
			value: "my-dc",
			validate: func(c *Config) bool {
				return c.DefaultDatacenter == "my-dc"
			},
		},
		{
			key:   "default_environment",
			value: "staging",
			validate: func(c *Config) bool {
				return c.DefaultEnvironment == "staging"
			},
		},
		{
			key:   "active_profile",
			value: "prod",
			validate: func(c *Config) bool {
				return c.ActiveProfile == "prod"
			},
		},
		{
			key:   "state.backend",
			value: "s3",
			validate: func(c *Config) bool {
				return c.State.Backend == "s3"
			},
		},
		{
			key:   "state.bucket",
			value: "my-bucket",
			validate: func(c *Config) bool {
				return c.State.Config["bucket"] == "my-bucket"
			},
		},
		{
			key:   "logging.level",
			value: "debug",
			validate: func(c *Config) bool {
				return c.Logging.Level == "debug"
			},
		},
		{
			key:   "logging.format",
			value: "json",
			validate: func(c *Config) bool {
				return c.Logging.Format == "json"
			},
		},
		{
			key:   "logging.file",
			value: "/var/log/arcctl.log",
			validate: func(c *Config) bool {
				return c.Logging.File == "/var/log/arcctl.log"
			},
		},
		{
			key:   "plugins.default",
			value: "opentofu",
			validate: func(c *Config) bool {
				return c.Plugins.Default == "opentofu"
			},
		},
		{
			key:   "registry.default",
			value: "gcr.io",
			validate: func(c *Config) bool {
				return c.Registry.Default == "gcr.io"
			},
		},
		{
			key:   "secrets.provider",
			value: "vault",
			validate: func(c *Config) bool {
				return c.Secrets.Provider == "vault"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			err := manager.Set(tt.key, tt.value)
			if err != nil {
				t.Fatalf("failed to set %s: %v", tt.key, err)
			}

			if !tt.validate(manager.Get()) {
				t.Errorf("validation failed for key %s", tt.key)
			}
		})
	}
}

func TestManager_Set_UnknownKey(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(""), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	manager, err := NewManagerWithPath(configPath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	err = manager.Set("unknown_key", "value")
	if err == nil {
		t.Error("expected error for unknown key")
	}
}

func TestManager_Reload(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Initial config
	initialData := `default_datacenter: initial-dc`
	if err := os.WriteFile(configPath, []byte(initialData), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	manager, err := NewManagerWithPath(configPath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	if manager.Get().DefaultDatacenter != "initial-dc" {
		t.Errorf("expected initial default_datacenter 'initial-dc'")
	}

	// Modify config file externally
	updatedData := `default_datacenter: updated-dc`
	if err := os.WriteFile(configPath, []byte(updatedData), 0600); err != nil {
		t.Fatalf("failed to update test config: %v", err)
	}

	// Reload
	if err := manager.Reload(); err != nil {
		t.Fatalf("failed to reload config: %v", err)
	}

	if manager.Get().DefaultDatacenter != "updated-dc" {
		t.Errorf("expected updated default_datacenter 'updated-dc', got %q", manager.Get().DefaultDatacenter)
	}
}

func TestManager_Save(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(""), 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	manager, err := NewManagerWithPath(configPath)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}

	// Modify config
	manager.Get().DefaultDatacenter = "saved-dc"

	// Save
	if err := manager.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Create new manager to load saved config
	newManager, err := NewManagerWithPath(configPath)
	if err != nil {
		t.Fatalf("failed to create new manager: %v", err)
	}

	if newManager.Get().DefaultDatacenter != "saved-dc" {
		t.Errorf("expected saved default_datacenter 'saved-dc', got %q", newManager.Get().DefaultDatacenter)
	}
}

func TestConfig_SaveToFile_Permissions(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	cfg := DefaultConfig()
	cfg.Registry.Auth = map[string]RegistryAuth{
		"docker.io": {
			Username: "user",
			Password: "secret",
		},
	}

	if err := cfg.SaveToFile(configPath); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Verify file permissions (should be 0600 for security)
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat config file: %v", err)
	}

	// On Unix systems, check that file is not world-readable
	mode := info.Mode().Perm()
	if mode&0044 != 0 {
		t.Errorf("config file should not be world-readable, got permissions %o", mode)
	}
}
