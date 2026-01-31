package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// VaultProvider provides secrets from HashiCorp Vault.
type VaultProvider struct {
	address   string
	token     string
	namespace string
	mountPath string
	client    *http.Client
}

// VaultConfig configures the Vault provider.
type VaultConfig struct {
	// Address is the Vault server address
	Address string

	// Token is the authentication token
	Token string

	// Namespace is the Vault namespace (Enterprise feature)
	Namespace string

	// MountPath is the secrets engine mount path (default: "secret")
	MountPath string
}

// NewVaultProvider creates a new Vault provider.
func NewVaultProvider(cfg VaultConfig) (*VaultProvider, error) {
	address := cfg.Address
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if address == "" {
		return nil, fmt.Errorf("vault address required")
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		// Try to read from token file
		if tokenFile := os.Getenv("VAULT_TOKEN_FILE"); tokenFile != "" {
			data, err := os.ReadFile(tokenFile)
			if err == nil {
				token = strings.TrimSpace(string(data))
			}
		}
		// Try default token file
		if token == "" {
			homeDir, _ := os.UserHomeDir()
			data, err := os.ReadFile(homeDir + "/.vault-token")
			if err == nil {
				token = strings.TrimSpace(string(data))
			}
		}
	}
	if token == "" {
		return nil, fmt.Errorf("vault token required")
	}

	namespace := cfg.Namespace
	if namespace == "" {
		namespace = os.Getenv("VAULT_NAMESPACE")
	}

	mountPath := cfg.MountPath
	if mountPath == "" {
		mountPath = "secret"
	}

	return &VaultProvider{
		address:   strings.TrimSuffix(address, "/"),
		token:     token,
		namespace: namespace,
		mountPath: mountPath,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (p *VaultProvider) Name() string {
	return "vault"
}

func (p *VaultProvider) Get(ctx context.Context, key string) (string, error) {
	// Parse key to get path and field
	// Format: path/to/secret#field or path/to/secret (returns "value" field)
	path, field := parseVaultKey(key)

	// Build URL for KV v2
	url := fmt.Sprintf("%s/v1/%s/data/%s", p.address, p.mountPath, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Vault-Token", p.token)
	if p.namespace != "" {
		req.Header.Set("X-Vault-Namespace", p.namespace)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("vault request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrSecretNotFound
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("vault error: %s", string(body))
	}

	var result struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse vault response: %w", err)
	}

	// Get the requested field
	value, ok := result.Data.Data[field]
	if !ok {
		return "", fmt.Errorf("field %s not found in secret %s", field, path)
	}

	// Convert to string
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}
}

func (p *VaultProvider) GetBatch(ctx context.Context, keys []string) (map[string]string, error) {
	results := make(map[string]string)
	for _, key := range keys {
		value, err := p.Get(ctx, key)
		if err == nil {
			results[key] = value
		}
	}
	return results, nil
}

func (p *VaultProvider) List(ctx context.Context, prefix string) ([]string, error) {
	url := fmt.Sprintf("%s/v1/%s/metadata/%s", p.address, p.mountPath, prefix)

	req, err := http.NewRequestWithContext(ctx, "LIST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Vault-Token", p.token)
	if p.namespace != "" {
		req.Header.Set("X-Vault-Namespace", p.namespace)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return []string{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vault error: %s", string(body))
	}

	var result struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse vault response: %w", err)
	}

	return result.Data.Keys, nil
}

func (p *VaultProvider) Set(ctx context.Context, key, value string) error {
	path, field := parseVaultKey(key)

	// Read existing data first
	existingData := make(map[string]interface{})
	url := fmt.Sprintf("%s/v1/%s/data/%s", p.address, p.mountPath, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", p.token)
	if p.namespace != "" {
		req.Header.Set("X-Vault-Namespace", p.namespace)
	}

	resp, err := p.client.Do(req)
	if err == nil && resp.StatusCode == http.StatusOK {
		var result struct {
			Data struct {
				Data map[string]interface{} `json:"data"`
			} `json:"data"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		existingData = result.Data.Data
	}
	if resp != nil {
		resp.Body.Close()
	}

	// Update the field
	existingData[field] = value

	// Write back
	payload := map[string]interface{}{
		"data": existingData,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err = http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return err
	}

	req.Header.Set("X-Vault-Token", p.token)
	req.Header.Set("Content-Type", "application/json")
	if p.namespace != "" {
		req.Header.Set("X-Vault-Namespace", p.namespace)
	}

	resp, err = p.client.Do(req)
	if err != nil {
		return fmt.Errorf("vault request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("vault error: %s", string(body))
	}

	return nil
}

func (p *VaultProvider) Delete(ctx context.Context, key string) error {
	path, _ := parseVaultKey(key)
	url := fmt.Sprintf("%s/v1/%s/metadata/%s", p.address, p.mountPath, path)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Vault-Token", p.token)
	if p.namespace != "" {
		req.Header.Set("X-Vault-Namespace", p.namespace)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("vault request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("vault error: %s", string(body))
	}

	return nil
}

func parseVaultKey(key string) (path, field string) {
	if idx := strings.LastIndex(key, "#"); idx != -1 {
		return key[:idx], key[idx+1:]
	}
	return key, "value"
}

// Ensure we implement the Provider interface
var _ Provider = (*VaultProvider)(nil)
