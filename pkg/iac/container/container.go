// Package container implements container-based IaC module execution.
// This allows IaC modules (Pulumi, OpenTofu) to be packaged as self-contained
// container images that include both the IaC code and runtime, eliminating
// the need for Pulumi/OpenTofu to be installed on the deployment host.
package container

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// ModuleRequest represents the input contract for a containerized module.
// This is passed to the container via a mounted JSON file.
type ModuleRequest struct {
	// Action is the operation to perform: "preview", "apply", "destroy", "refresh"
	Action string `json:"action"`

	// Inputs are the module input values
	Inputs map[string]interface{} `json:"inputs"`

	// State is the current module state (for updates/destroys)
	State map[string]interface{} `json:"state,omitempty"`

	// Environment variables to set
	Environment map[string]string `json:"environment,omitempty"`

	// StackName for Pulumi or workspace name for OpenTofu
	StackName string `json:"stack_name,omitempty"`

	// Backend configuration for state storage
	Backend *BackendConfig `json:"backend,omitempty"`
}

// BackendConfig configures state storage for the module.
type BackendConfig struct {
	// Type is the backend type (e.g., "s3", "gcs", "azurerm", "local")
	Type string `json:"type"`

	// Config contains backend-specific configuration
	Config map[string]string `json:"config"`
}

// ModuleResponse represents the output contract from a containerized module.
// The container writes this as JSON to a mounted output file.
type ModuleResponse struct {
	// Success indicates whether the operation succeeded
	Success bool `json:"success"`

	// Action that was performed
	Action string `json:"action"`

	// Outputs from the module (after apply)
	Outputs map[string]OutputValue `json:"outputs,omitempty"`

	// State to persist (opaque to cldctl)
	State map[string]interface{} `json:"state,omitempty"`

	// Changes describes what changed (for preview)
	Changes []ResourceChange `json:"changes,omitempty"`

	// Error message if Success is false
	Error string `json:"error,omitempty"`

	// Logs from the operation
	Logs string `json:"logs,omitempty"`
}

// OutputValue represents a module output.
type OutputValue struct {
	Value     interface{} `json:"value"`
	Sensitive bool        `json:"sensitive,omitempty"`
}

// ResourceChange describes a planned or executed change.
type ResourceChange struct {
	// Resource identifier
	Resource string `json:"resource"`

	// Action: create, update, delete, replace, no-op
	Action string `json:"action"`

	// Before state (for updates/deletes)
	Before map[string]interface{} `json:"before,omitempty"`

	// After state (for creates/updates)
	After map[string]interface{} `json:"after,omitempty"`
}

// Executor runs containerized IaC modules.
type Executor struct {
	dockerClient *client.Client
}

// NewExecutor creates a new container executor.
func NewExecutor() (*Executor, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &Executor{dockerClient: cli}, nil
}

// ExecuteOptions configures module execution.
type ExecuteOptions struct {
	// Image is the container image to run
	Image string

	// Request is the module request
	Request *ModuleRequest

	// WorkDir is where to store temporary files
	WorkDir string

	// Credentials for cloud providers
	Credentials map[string]string

	// Stdout for streaming output
	Stdout io.Writer

	// Stderr for streaming errors
	Stderr io.Writer
}

// Execute runs a containerized module and returns the response.
func (e *Executor) Execute(ctx context.Context, opts ExecuteOptions) (*ModuleResponse, error) {
	// Create work directory if needed
	if opts.WorkDir == "" {
		tmpDir, err := os.MkdirTemp("", "cldctl-module-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp dir: %w", err)
		}
		opts.WorkDir = tmpDir
		defer os.RemoveAll(tmpDir)
	}

	// Write request to input file
	inputFile := filepath.Join(opts.WorkDir, "input.json")
	inputData, err := json.MarshalIndent(opts.Request, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	if err := os.WriteFile(inputFile, inputData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write input file: %w", err)
	}

	// Prepare output file path
	outputFile := filepath.Join(opts.WorkDir, "output.json")

	// Pull image if needed
	reader, err := e.dockerClient.ImagePull(ctx, opts.Image, image.PullOptions{})
	if err != nil {
		// Image might already exist locally, continue
	} else {
		_, _ = io.Copy(io.Discard, reader)
		reader.Close()
	}

	// Build environment variables
	env := []string{}
	for k, v := range opts.Request.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	for k, v := range opts.Credentials {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// Create container
	containerConfig := &container.Config{
		Image: opts.Image,
		Env:   env,
		Cmd:   []string{"/cldctl-entrypoint", "--input", "/workspace/input.json", "--output", "/workspace/output.json"},
	}

	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: opts.WorkDir,
				Target: "/workspace",
			},
		},
		AutoRemove: true,
	}

	resp, err := e.dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := e.dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Stream logs if writers provided
	if opts.Stdout != nil || opts.Stderr != nil {
		logReader, err := e.dockerClient.ContainerLogs(ctx, resp.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		})
		if err == nil {
			stdout := opts.Stdout
			stderr := opts.Stderr
			if stdout == nil {
				stdout = io.Discard
			}
			if stderr == nil {
				stderr = io.Discard
			}
			go func() { _, _ = stdcopy.StdCopy(stdout, stderr, logReader) }()
		}
	}

	// Wait for container to finish
	statusCh, errCh := e.dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, fmt.Errorf("container wait failed: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return &ModuleResponse{
				Success: false,
				Action:  opts.Request.Action,
				Error:   fmt.Sprintf("container exited with code %d", status.StatusCode),
			}, nil
		}
	}

	// Read output file
	outputData, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	var response ModuleResponse
	if err := json.Unmarshal(outputData, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// Close releases resources.
func (e *Executor) Close() error {
	return e.dockerClient.Close()
}

// ModuleType identifies the IaC framework.
type ModuleType string

const (
	ModuleTypePulumi   ModuleType = "pulumi"
	ModuleTypeOpenTofu ModuleType = "opentofu"
)

// DetectModuleType detects the IaC framework from a module directory.
func DetectModuleType(dir string) (ModuleType, error) {
	// Check for Pulumi
	if _, err := os.Stat(filepath.Join(dir, "Pulumi.yaml")); err == nil {
		return ModuleTypePulumi, nil
	}

	// Check for OpenTofu/Terraform
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tf") {
			return ModuleTypeOpenTofu, nil
		}
	}

	return "", fmt.Errorf("unable to detect module type in %s", dir)
}
