package v1

import (
	"testing"
)

func TestParser_ParseBytes(t *testing.T) {
	parser := &Parser{}

	yaml := `
databases:
  main:
    type: postgres:^15

deployments:
  api:
    image: nginx:latest
    environment:
      PORT: "8080"
    replicas: 2

services:
  api:
    deployment: api
    port: 8080

variables:
  log_level:
    default: info
`

	schema, err := parser.ParseBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	// Check databases
	if len(schema.Databases) != 1 {
		t.Errorf("expected 1 database, got %d", len(schema.Databases))
	}

	if db, ok := schema.Databases["main"]; ok {
		if db.Type != "postgres:^15" {
			t.Errorf("expected postgres type, got %s", db.Type)
		}
	} else {
		t.Error("expected 'main' database")
	}

	// Check deployments
	if len(schema.Deployments) != 1 {
		t.Errorf("expected 1 deployment, got %d", len(schema.Deployments))
	}

	if deploy, ok := schema.Deployments["api"]; ok {
		if deploy.Image != "nginx:latest" {
			t.Errorf("expected nginx image, got %s", deploy.Image)
		}
		if deploy.Replicas != 2 {
			t.Errorf("expected 2 replicas, got %d", deploy.Replicas)
		}
	} else {
		t.Error("expected 'api' deployment")
	}

	// Check services
	if len(schema.Services) != 1 {
		t.Errorf("expected 1 service, got %d", len(schema.Services))
	}

	// Check variables
	if len(schema.Variables) != 1 {
		t.Errorf("expected 1 variable, got %d", len(schema.Variables))
	}
}

func TestParser_ParseBytes_Invalid(t *testing.T) {
	parser := &Parser{}

	invalidYAML := `
databases:
  - invalid: list format
`

	_, err := parser.ParseBytes([]byte(invalidYAML))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParser_ParseBytes_Empty(t *testing.T) {
	parser := &Parser{}

	schema, err := parser.ParseBytes([]byte(""))
	if err != nil {
		t.Fatalf("failed to parse empty: %v", err)
	}

	// Empty schema should be valid (no required fields)
	if len(schema.Databases) != 0 {
		t.Error("expected empty databases")
	}
}
