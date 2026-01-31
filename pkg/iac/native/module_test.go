package native

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadModule_Success(t *testing.T) {
	// Create a temporary directory with a module file
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	moduleContent := `
plugin: native
type: docker:container
inputs:
  image:
    type: string
    required: true
    description: Docker image to use
  port:
    type: number
    default: 8080
    description: Port to expose
resources:
  container:
    type: docker:container
    properties:
      image: ${inputs.image}
      ports:
        - container: ${inputs.port}
          host: auto
outputs:
  container_id:
    value: ${resources.container.outputs.container_id}
    description: The container ID
  port:
    value: ${resources.container.outputs.ports}
    sensitive: false
`

	moduleFile := filepath.Join(tmpDir, "module.yml")
	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	module, err := LoadModule(tmpDir)
	if err != nil {
		t.Fatalf("failed to load module: %v", err)
	}

	if module.Plugin != "native" {
		t.Errorf("expected plugin 'native', got %q", module.Plugin)
	}

	if module.Type != "docker:container" {
		t.Errorf("expected type 'docker:container', got %q", module.Type)
	}

	// Check inputs
	if len(module.Inputs) != 2 {
		t.Errorf("expected 2 inputs, got %d", len(module.Inputs))
	}

	imageInput := module.Inputs["image"]
	if imageInput.Type != "string" {
		t.Errorf("expected image type 'string', got %q", imageInput.Type)
	}
	if !imageInput.Required {
		t.Error("expected image to be required")
	}

	portInput := module.Inputs["port"]
	if portInput.Type != "number" {
		t.Errorf("expected port type 'number', got %q", portInput.Type)
	}
	if portInput.Default != 8080 {
		t.Errorf("expected port default 8080, got %v", portInput.Default)
	}

	// Check resources
	if len(module.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(module.Resources))
	}

	container := module.Resources["container"]
	if container.Type != "docker:container" {
		t.Errorf("expected resource type 'docker:container', got %q", container.Type)
	}

	// Check outputs
	if len(module.Outputs) != 2 {
		t.Errorf("expected 2 outputs, got %d", len(module.Outputs))
	}
}

func TestLoadModule_YamlExtension(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	moduleContent := `
plugin: native
inputs: {}
resources: {}
outputs: {}
`

	moduleFile := filepath.Join(tmpDir, "module.yaml")
	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	module, err := LoadModule(tmpDir)
	if err != nil {
		t.Fatalf("failed to load module: %v", err)
	}

	if module.Plugin != "native" {
		t.Errorf("expected plugin 'native', got %q", module.Plugin)
	}
}

func TestLoadModule_DirectFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	moduleContent := `
plugin: native
inputs: {}
resources: {}
outputs: {}
`

	moduleFile := filepath.Join(tmpDir, "custom-module.yml")
	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	module, err := LoadModule(moduleFile)
	if err != nil {
		t.Fatalf("failed to load module: %v", err)
	}

	if module.Plugin != "native" {
		t.Errorf("expected plugin 'native', got %q", module.Plugin)
	}
}

func TestLoadModule_NotFound(t *testing.T) {
	_, err := LoadModule("/nonexistent/path")
	if err == nil {
		t.Error("expected error for non-existent module")
	}
}

func TestLoadModule_InvalidYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	moduleFile := filepath.Join(tmpDir, "module.yml")
	if err := os.WriteFile(moduleFile, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	_, err = LoadModule(tmpDir)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadModule_InvalidPlugin(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	moduleContent := `
plugin: pulumi
inputs: {}
resources: {}
outputs: {}
`

	moduleFile := filepath.Join(tmpDir, "module.yml")
	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	_, err = LoadModule(tmpDir)
	if err == nil {
		t.Error("expected error for invalid plugin type")
	}
}

func TestLoadModule_EmptyPlugin(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Empty plugin is allowed (defaults to native)
	moduleContent := `
inputs: {}
resources: {}
outputs: {}
`

	moduleFile := filepath.Join(tmpDir, "module.yml")
	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	module, err := LoadModule(tmpDir)
	if err != nil {
		t.Fatalf("failed to load module: %v", err)
	}

	if module.Plugin != "" {
		t.Errorf("expected empty plugin, got %q", module.Plugin)
	}
}

func TestLoadModule_ComplexInputs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	moduleContent := `
plugin: native
inputs:
  string_input:
    type: string
    required: true
    description: A string input
  number_input:
    type: number
    default: 42
  bool_input:
    type: boolean
    default: true
  sensitive_input:
    type: string
    sensitive: true
resources: {}
outputs: {}
`

	moduleFile := filepath.Join(tmpDir, "module.yml")
	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	module, err := LoadModule(tmpDir)
	if err != nil {
		t.Fatalf("failed to load module: %v", err)
	}

	if len(module.Inputs) != 4 {
		t.Errorf("expected 4 inputs, got %d", len(module.Inputs))
	}

	stringInput := module.Inputs["string_input"]
	if stringInput.Type != "string" {
		t.Errorf("expected string_input type 'string', got %q", stringInput.Type)
	}
	if !stringInput.Required {
		t.Error("expected string_input to be required")
	}

	numberInput := module.Inputs["number_input"]
	if numberInput.Type != "number" {
		t.Errorf("expected number_input type 'number', got %q", numberInput.Type)
	}
	if numberInput.Default != 42 {
		t.Errorf("expected number_input default 42, got %v", numberInput.Default)
	}

	boolInput := module.Inputs["bool_input"]
	if boolInput.Type != "boolean" {
		t.Errorf("expected bool_input type 'boolean', got %q", boolInput.Type)
	}
	if boolInput.Default != true {
		t.Errorf("expected bool_input default true, got %v", boolInput.Default)
	}

	sensitiveInput := module.Inputs["sensitive_input"]
	if !sensitiveInput.Sensitive {
		t.Error("expected sensitive_input to be sensitive")
	}
}

func TestLoadModule_ComplexResources(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	moduleContent := `
plugin: native
inputs: {}
resources:
  network:
    type: docker:network
    properties:
      name: my-network
  container:
    type: docker:container
    depends_on:
      - network
    properties:
      image: nginx:latest
      name: my-container
      network: ${resources.network.outputs.name}
outputs: {}
`

	moduleFile := filepath.Join(tmpDir, "module.yml")
	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	module, err := LoadModule(tmpDir)
	if err != nil {
		t.Fatalf("failed to load module: %v", err)
	}

	if len(module.Resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(module.Resources))
	}

	network := module.Resources["network"]
	if network.Type != "docker:network" {
		t.Errorf("expected network type 'docker:network', got %q", network.Type)
	}

	container := module.Resources["container"]
	if container.Type != "docker:container" {
		t.Errorf("expected container type 'docker:container', got %q", container.Type)
	}
	if len(container.DependsOn) != 1 {
		t.Errorf("expected 1 dependency, got %d", len(container.DependsOn))
	}
	if container.DependsOn[0] != "network" {
		t.Errorf("expected dependency 'network', got %q", container.DependsOn[0])
	}
}

func TestLoadModule_ComplexOutputs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "module-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	moduleContent := `
plugin: native
inputs: {}
resources: {}
outputs:
  endpoint:
    value: ${resources.container.outputs.endpoint}
    description: The service endpoint
  password:
    value: ${resources.database.outputs.password}
    sensitive: true
`

	moduleFile := filepath.Join(tmpDir, "module.yml")
	if err := os.WriteFile(moduleFile, []byte(moduleContent), 0644); err != nil {
		t.Fatalf("failed to write module file: %v", err)
	}

	module, err := LoadModule(tmpDir)
	if err != nil {
		t.Fatalf("failed to load module: %v", err)
	}

	if len(module.Outputs) != 2 {
		t.Errorf("expected 2 outputs, got %d", len(module.Outputs))
	}

	endpoint := module.Outputs["endpoint"]
	if endpoint.Description != "The service endpoint" {
		t.Errorf("expected endpoint description 'The service endpoint', got %q", endpoint.Description)
	}
	if endpoint.Sensitive {
		t.Error("expected endpoint to not be sensitive")
	}

	password := module.Outputs["password"]
	if !password.Sensitive {
		t.Error("expected password to be sensitive")
	}
}

func TestModule_Struct(t *testing.T) {
	module := Module{
		Plugin: "native",
		Type:   "docker:container",
		Inputs: map[string]InputDef{
			"image": {
				Type:        "string",
				Required:    true,
				Default:     nil,
				Description: "Docker image",
				Sensitive:   false,
			},
		},
		Resources: map[string]Resource{
			"container": {
				Type: "docker:container",
				Properties: map[string]interface{}{
					"image": "${inputs.image}",
				},
				DependsOn: []string{},
			},
		},
		Outputs: map[string]OutputDef{
			"id": {
				Value:       "${resources.container.id}",
				Description: "Container ID",
				Sensitive:   false,
			},
		},
	}

	if module.Plugin != "native" {
		t.Errorf("expected plugin 'native', got %q", module.Plugin)
	}
	if module.Type != "docker:container" {
		t.Errorf("expected type 'docker:container', got %q", module.Type)
	}
}

func TestInputDef_Struct(t *testing.T) {
	input := InputDef{
		Type:        "string",
		Required:    true,
		Default:     "default-value",
		Description: "Test input",
		Sensitive:   true,
	}

	if input.Type != "string" {
		t.Errorf("expected type 'string', got %q", input.Type)
	}
	if !input.Required {
		t.Error("expected required to be true")
	}
	if input.Default != "default-value" {
		t.Errorf("expected default 'default-value', got %v", input.Default)
	}
	if input.Description != "Test input" {
		t.Errorf("expected description 'Test input', got %q", input.Description)
	}
	if !input.Sensitive {
		t.Error("expected sensitive to be true")
	}
}

func TestResource_Struct(t *testing.T) {
	resource := Resource{
		Type: "docker:container",
		Properties: map[string]interface{}{
			"image": "nginx:latest",
			"port":  8080,
		},
		DependsOn: []string{"network", "volume"},
	}

	if resource.Type != "docker:container" {
		t.Errorf("expected type 'docker:container', got %q", resource.Type)
	}
	if resource.Properties["image"] != "nginx:latest" {
		t.Errorf("expected image 'nginx:latest', got %v", resource.Properties["image"])
	}
	if len(resource.DependsOn) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(resource.DependsOn))
	}
}

func TestOutputDef_Struct(t *testing.T) {
	output := OutputDef{
		Value:       "${resources.container.id}",
		Description: "Container ID",
		Sensitive:   false,
	}

	if output.Value != "${resources.container.id}" {
		t.Errorf("expected value '${resources.container.id}', got %q", output.Value)
	}
	if output.Description != "Container ID" {
		t.Errorf("expected description 'Container ID', got %q", output.Description)
	}
	if output.Sensitive {
		t.Error("expected sensitive to be false")
	}
}

func TestState_Struct(t *testing.T) {
	state := State{
		ModulePath: "/path/to/module",
		Inputs: map[string]interface{}{
			"image": "nginx:latest",
		},
		Resources: map[string]*ResourceState{
			"container": {
				Type: "docker:container",
				ID:   "container-123",
				Properties: map[string]interface{}{
					"image": "nginx:latest",
				},
				Outputs: map[string]interface{}{
					"container_id": "container-123",
				},
			},
		},
		Outputs: map[string]interface{}{
			"id": "container-123",
		},
	}

	if state.ModulePath != "/path/to/module" {
		t.Errorf("expected module path '/path/to/module', got %q", state.ModulePath)
	}
	if state.Inputs["image"] != "nginx:latest" {
		t.Errorf("expected input image 'nginx:latest', got %v", state.Inputs["image"])
	}
	if state.Resources["container"].ID != "container-123" {
		t.Errorf("expected resource ID 'container-123', got %v", state.Resources["container"].ID)
	}
	if state.Outputs["id"] != "container-123" {
		t.Errorf("expected output id 'container-123', got %v", state.Outputs["id"])
	}
}

func TestResourceState_Struct(t *testing.T) {
	rs := ResourceState{
		Type: "docker:container",
		ID:   "container-id",
		Properties: map[string]interface{}{
			"image": "nginx:latest",
		},
		Outputs: map[string]interface{}{
			"container_id": "container-id",
			"ports": map[string]interface{}{
				"8080/tcp": 32000,
			},
		},
	}

	if rs.Type != "docker:container" {
		t.Errorf("expected type 'docker:container', got %q", rs.Type)
	}
	if rs.ID != "container-id" {
		t.Errorf("expected ID 'container-id', got %v", rs.ID)
	}
	if rs.Properties["image"] != "nginx:latest" {
		t.Errorf("expected property image 'nginx:latest', got %v", rs.Properties["image"])
	}
	if rs.Outputs["container_id"] != "container-id" {
		t.Errorf("expected output container_id 'container-id', got %v", rs.Outputs["container_id"])
	}
}
