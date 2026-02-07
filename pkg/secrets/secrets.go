// Package secrets provides secret management integration.
package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// Provider is the interface for secret providers.
type Provider interface {
	// Name returns the provider name
	Name() string

	// Get retrieves a secret value
	Get(ctx context.Context, key string) (string, error)

	// GetBatch retrieves multiple secrets
	GetBatch(ctx context.Context, keys []string) (map[string]string, error)

	// List lists available secret keys
	List(ctx context.Context, prefix string) ([]string, error)

	// Set stores a secret (if supported)
	Set(ctx context.Context, key, value string) error

	// Delete removes a secret (if supported)
	Delete(ctx context.Context, key string) error
}

// Manager manages multiple secret providers.
type Manager struct {
	mu        sync.RWMutex
	providers map[string]Provider
	priority  []string
	cache     *secretCache
}

// NewManager creates a new secret manager.
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]Provider),
		priority:  []string{},
		cache:     newSecretCache(),
	}
}

// RegisterProvider registers a secret provider.
func (m *Manager) RegisterProvider(p Provider) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.providers[p.Name()] = p
	m.priority = append(m.priority, p.Name())
}

// SetPriority sets the provider lookup priority.
func (m *Manager) SetPriority(providers []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.priority = providers
}

// Get retrieves a secret, trying providers in priority order.
func (m *Manager) Get(ctx context.Context, key string) (string, error) {
	// Check cache first
	if value, ok := m.cache.get(key); ok {
		return value, nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var lastErr error
	for _, name := range m.priority {
		provider, ok := m.providers[name]
		if !ok {
			continue
		}

		value, err := provider.Get(ctx, key)
		if err == nil {
			m.cache.set(key, value)
			return value, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return "", fmt.Errorf("secret %s not found: %w", key, lastErr)
	}
	return "", fmt.Errorf("secret %s not found: no providers configured", key)
}

// GetFromProvider retrieves a secret from a specific provider.
func (m *Manager) GetFromProvider(ctx context.Context, providerName, key string) (string, error) {
	m.mu.RLock()
	provider, ok := m.providers[providerName]
	m.mu.RUnlock()

	if !ok {
		return "", fmt.Errorf("unknown provider: %s", providerName)
	}

	return provider.Get(ctx, key)
}

// GetBatch retrieves multiple secrets.
func (m *Manager) GetBatch(ctx context.Context, keys []string) (map[string]string, error) {
	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()

			value, err := m.Get(ctx, k)
			if err == nil {
				mu.Lock()
				results[k] = value
				mu.Unlock()
			}
		}(key)
	}

	wg.Wait()
	return results, nil
}

// ResolveSecrets resolves secret references in a map.
// Format: ${secret:provider:key} or ${secret:key}
func (m *Manager) ResolveSecrets(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for k, v := range data {
		resolved, err := m.resolveValue(ctx, v)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve %s: %w", k, err)
		}
		result[k] = resolved
	}

	return result, nil
}

func (m *Manager) resolveValue(ctx context.Context, v interface{}) (interface{}, error) {
	switch val := v.(type) {
	case string:
		return m.resolveString(ctx, val)
	case map[string]interface{}:
		return m.ResolveSecrets(ctx, val)
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			resolved, err := m.resolveValue(ctx, item)
			if err != nil {
				return nil, err
			}
			result[i] = resolved
		}
		return result, nil
	default:
		return v, nil
	}
}

func (m *Manager) resolveString(ctx context.Context, s string) (string, error) {
	// Look for ${secret:...} patterns
	if !strings.Contains(s, "${secret:") {
		return s, nil
	}

	result := s
	for {
		start := strings.Index(result, "${secret:")
		if start == -1 {
			break
		}

		end := strings.Index(result[start:], "}")
		if end == -1 {
			return "", fmt.Errorf("unclosed secret reference in: %s", s)
		}
		end += start

		ref := result[start+9 : end]
		var value string
		var err error

		// Check for provider prefix
		if strings.Contains(ref, ":") {
			parts := strings.SplitN(ref, ":", 2)
			value, err = m.GetFromProvider(ctx, parts[0], parts[1])
		} else {
			value, err = m.Get(ctx, ref)
		}

		if err != nil {
			return "", err
		}

		result = result[:start] + value + result[end+1:]
	}

	return result, nil
}

// ClearCache clears the secret cache.
func (m *Manager) ClearCache() {
	m.cache.clear()
}

// secretCache provides a simple in-memory cache for secrets.
type secretCache struct {
	mu    sync.RWMutex
	items map[string]string
}

func newSecretCache() *secretCache {
	return &secretCache{
		items: make(map[string]string),
	}
}

func (c *secretCache) get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.items[key]
	return v, ok
}

func (c *secretCache) set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *secretCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]string)
}

// ErrSecretNotFound is returned when a secret is not found.
var ErrSecretNotFound = fmt.Errorf("secret not found")

// DefaultManager returns a manager with default providers.
func DefaultManager() *Manager {
	m := NewManager()
	m.RegisterProvider(NewEnvProvider())
	return m
}

// EnvProvider provides secrets from environment variables.
type EnvProvider struct {
	prefix string
}

// NewEnvProvider creates a new environment variable provider.
func NewEnvProvider() *EnvProvider {
	return &EnvProvider{
		prefix: "CLDCTL_SECRET_",
	}
}

// NewEnvProviderWithPrefix creates an env provider with a custom prefix.
func NewEnvProviderWithPrefix(prefix string) *EnvProvider {
	return &EnvProvider{
		prefix: prefix,
	}
}

func (p *EnvProvider) Name() string {
	return "env"
}

func (p *EnvProvider) Get(ctx context.Context, key string) (string, error) {
	// Try with prefix first
	envKey := p.prefix + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
	if value := os.Getenv(envKey); value != "" {
		return value, nil
	}

	// Try without prefix
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	return "", ErrSecretNotFound
}

func (p *EnvProvider) GetBatch(ctx context.Context, keys []string) (map[string]string, error) {
	results := make(map[string]string)
	for _, key := range keys {
		if value, err := p.Get(ctx, key); err == nil {
			results[key] = value
		}
	}
	return results, nil
}

func (p *EnvProvider) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	fullPrefix := p.prefix + strings.ToUpper(strings.ReplaceAll(prefix, "-", "_"))

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && strings.HasPrefix(parts[0], fullPrefix) {
			key := strings.TrimPrefix(parts[0], p.prefix)
			key = strings.ToLower(strings.ReplaceAll(key, "_", "-"))
			keys = append(keys, key)
		}
	}

	return keys, nil
}

func (p *EnvProvider) Set(ctx context.Context, key, value string) error {
	envKey := p.prefix + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
	return os.Setenv(envKey, value)
}

func (p *EnvProvider) Delete(ctx context.Context, key string) error {
	envKey := p.prefix + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
	return os.Unsetenv(envKey)
}

// FileProvider provides secrets from a file.
type FileProvider struct {
	secrets map[string]string
}

// NewFileProvider creates a new file-based secret provider.
func NewFileProvider(secrets map[string]string) *FileProvider {
	return &FileProvider{
		secrets: secrets,
	}
}

func (p *FileProvider) Name() string {
	return "file"
}

func (p *FileProvider) Get(ctx context.Context, key string) (string, error) {
	if value, ok := p.secrets[key]; ok {
		return value, nil
	}
	return "", ErrSecretNotFound
}

func (p *FileProvider) GetBatch(ctx context.Context, keys []string) (map[string]string, error) {
	results := make(map[string]string)
	for _, key := range keys {
		if value, ok := p.secrets[key]; ok {
			results[key] = value
		}
	}
	return results, nil
}

func (p *FileProvider) List(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	for key := range p.secrets {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func (p *FileProvider) Set(ctx context.Context, key, value string) error {
	p.secrets[key] = value
	return nil
}

func (p *FileProvider) Delete(ctx context.Context, key string) error {
	delete(p.secrets, key)
	return nil
}
