package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

// AWSProvider provides secrets from AWS Secrets Manager.
type AWSProvider struct {
	client *secretsmanager.Client
	prefix string
}

// AWSConfig configures the AWS Secrets Manager provider.
type AWSConfig struct {
	// Region is the AWS region
	Region string

	// AccessKeyID is the AWS access key (optional, uses default credentials chain)
	AccessKeyID string

	// SecretAccessKey is the AWS secret key
	SecretAccessKey string

	// Prefix is the prefix for secret names
	Prefix string

	// Endpoint is a custom endpoint URL (for localstack, etc.)
	Endpoint string
}

// NewAWSProvider creates a new AWS Secrets Manager provider.
func NewAWSProvider(ctx context.Context, cfg AWSConfig) (*AWSProvider, error) {
	// Build AWS config options
	var opts []func(*config.LoadOptions) error

	if cfg.Region != "" {
		opts = append(opts, config.WithRegion(cfg.Region))
	}

	// Support explicit credentials
	if cfg.AccessKeyID != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		))
	}

	// Load AWS config
	awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Secrets Manager client
	var clientOpts []func(*secretsmanager.Options)
	if cfg.Endpoint != "" {
		clientOpts = append(clientOpts, func(o *secretsmanager.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		})
	}

	client := secretsmanager.NewFromConfig(awsCfg, clientOpts...)

	return &AWSProvider{
		client: client,
		prefix: cfg.Prefix,
	}, nil
}

func (p *AWSProvider) Name() string {
	return "aws"
}

func (p *AWSProvider) Get(ctx context.Context, key string) (string, error) {
	// Parse key to get secret name and field
	// Format: secret-name#field or secret-name
	secretName, field := parseAWSKey(key)

	// Add prefix if configured
	if p.prefix != "" {
		secretName = p.prefix + secretName
	}

	// Get secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	output, err := p.client.GetSecretValue(ctx, input)
	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			return "", ErrSecretNotFound
		}
		return "", fmt.Errorf("failed to get secret: %w", err)
	}

	// Get secret string
	secretValue := ""
	if output.SecretString != nil {
		secretValue = *output.SecretString
	} else if output.SecretBinary != nil {
		secretValue = string(output.SecretBinary)
	}

	// If no field specified, return the entire secret
	if field == "" {
		return secretValue, nil
	}

	// Try to parse as JSON and extract field
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(secretValue), &jsonData); err != nil {
		return "", fmt.Errorf("secret is not JSON and field was specified")
	}

	value, ok := jsonData[field]
	if !ok {
		return "", fmt.Errorf("field %s not found in secret %s", field, secretName)
	}

	switch v := value.(type) {
	case string:
		return v, nil
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}
}

func (p *AWSProvider) GetBatch(ctx context.Context, keys []string) (map[string]string, error) {
	results := make(map[string]string)

	// AWS Secrets Manager supports batch retrieval
	// Build list of secret IDs
	secretIds := make([]string, len(keys))
	keyMap := make(map[string]string) // Maps full name to original key

	for i, key := range keys {
		secretName, _ := parseAWSKey(key)
		if p.prefix != "" {
			secretName = p.prefix + secretName
		}
		secretIds[i] = secretName
		keyMap[secretName] = key
	}

	input := &secretsmanager.BatchGetSecretValueInput{
		SecretIdList: secretIds,
	}

	output, err := p.client.BatchGetSecretValue(ctx, input)
	if err != nil {
		// Fall back to individual gets
		for _, key := range keys {
			value, err := p.Get(ctx, key)
			if err == nil {
				results[key] = value
			}
		}
		return results, nil
	}

	// Process results
	for _, secret := range output.SecretValues {
		if secret.SecretString == nil {
			continue
		}

		secretName := *secret.Name
		originalKey := keyMap[secretName]
		_, field := parseAWSKey(originalKey)

		if field == "" {
			results[originalKey] = *secret.SecretString
			continue
		}

		// Extract field from JSON
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(*secret.SecretString), &jsonData); err == nil {
			if value, ok := jsonData[field]; ok {
				switch v := value.(type) {
				case string:
					results[originalKey] = v
				default:
					if jsonBytes, err := json.Marshal(v); err == nil {
						results[originalKey] = string(jsonBytes)
					}
				}
			}
		}
	}

	return results, nil
}

func (p *AWSProvider) List(ctx context.Context, prefix string) ([]string, error) {
	fullPrefix := p.prefix + prefix

	input := &secretsmanager.ListSecretsInput{}
	if fullPrefix != "" {
		input.Filters = []types.Filter{
			{
				Key:    types.FilterNameStringTypeName,
				Values: []string{fullPrefix},
			},
		}
	}

	var keys []string
	paginator := secretsmanager.NewListSecretsPaginator(p.client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list secrets: %w", err)
		}

		for _, secret := range page.SecretList {
			if secret.Name != nil {
				// Remove prefix from name
				name := *secret.Name
				if p.prefix != "" {
					name = strings.TrimPrefix(name, p.prefix)
				}
				keys = append(keys, name)
			}
		}
	}

	return keys, nil
}

func (p *AWSProvider) Set(ctx context.Context, key, value string) error {
	secretName, _ := parseAWSKey(key)
	if p.prefix != "" {
		secretName = p.prefix + secretName
	}

	// Try to update existing secret
	_, err := p.client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(secretName),
		SecretString: aws.String(value),
	})

	if err != nil {
		// If secret doesn't exist, create it
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			_, err = p.client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
				Name:         aws.String(secretName),
				SecretString: aws.String(value),
			})
			if err != nil {
				return fmt.Errorf("failed to create secret: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to update secret: %w", err)
	}

	return nil
}

func (p *AWSProvider) Delete(ctx context.Context, key string) error {
	secretName, _ := parseAWSKey(key)
	if p.prefix != "" {
		secretName = p.prefix + secretName
	}

	_, err := p.client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
		SecretId:                   aws.String(secretName),
		ForceDeleteWithoutRecovery: aws.Bool(true),
	})

	if err != nil {
		if strings.Contains(err.Error(), "ResourceNotFoundException") {
			return nil // Already deleted
		}
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}

func parseAWSKey(key string) (secretName, field string) {
	if idx := strings.LastIndex(key, "#"); idx != -1 {
		return key[:idx], key[idx+1:]
	}
	return key, ""
}

// Ensure we implement the Provider interface
var _ Provider = (*AWSProvider)(nil)
