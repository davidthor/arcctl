package v1

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Parser parses v1 environment schemas.
type Parser struct{}

// NewParser creates a new v1 environment parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses an environment configuration from the given file path.
func (p *Parser) Parse(path string) (*SchemaV1, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return p.ParseBytes(data)
}

// ParseBytes parses an environment configuration from raw bytes.
func (p *Parser) ParseBytes(data []byte) (*SchemaV1, error) {
	var schema SchemaV1
	if err := yaml.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &schema, nil
}
