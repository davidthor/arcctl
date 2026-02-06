package container

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/client"
)

// Builder builds container images for IaC modules.
type Builder struct {
	dockerClient *client.Client
}

// NewBuilder creates a new module builder.
func NewBuilder() (*Builder, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &Builder{dockerClient: cli}, nil
}

// BuildOptions configures module image building.
type BuildOptions struct {
	// ModuleDir is the directory containing the IaC module
	ModuleDir string

	// ModuleType is the IaC framework (auto-detected if empty)
	ModuleType ModuleType

	// Tag is the image tag
	Tag string

	// Output for build logs
	Output io.Writer
}

// BuildResult contains the result of a module build.
type BuildResult struct {
	// Image is the built image tag
	Image string

	// Digest is the image digest
	Digest string

	// ModuleType is the detected/specified module type
	ModuleType ModuleType
}

// Build builds a container image for an IaC module.
func (b *Builder) Build(ctx context.Context, opts BuildOptions) (*BuildResult, error) {
	// Detect module type if not specified
	moduleType := opts.ModuleType
	if moduleType == "" {
		detected, err := DetectModuleType(opts.ModuleDir)
		if err != nil {
			return nil, fmt.Errorf("failed to detect module type: %w", err)
		}
		moduleType = detected
	}

	// Generate Dockerfile
	dockerfile, err := generateDockerfile(moduleType, opts.ModuleDir)
	if err != nil {
		return nil, fmt.Errorf("failed to generate dockerfile: %w", err)
	}

	// Create build context (tar archive)
	buildContext, err := createBuildContext(opts.ModuleDir, dockerfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create build context: %w", err)
	}

	// Build the image
	buildResp, err := b.dockerClient.ImageBuild(ctx, buildContext, build.ImageBuildOptions{
		Tags:       []string{opts.Tag},
		Dockerfile: "Dockerfile",
		Remove:     true,
		// Use multi-platform if available
		Platform: "linux/amd64",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to build image: %w", err)
	}
	defer buildResp.Body.Close()

	// Stream build output
	output := opts.Output
	if output == nil {
		output = io.Discard
	}

	// Parse build output for errors and digest
	decoder := json.NewDecoder(buildResp.Body)
	var lastLine struct {
		Stream string `json:"stream"`
		Error  string `json:"error"`
		Aux    struct {
			ID string `json:"ID"`
		} `json:"aux"`
	}

	for {
		if err := decoder.Decode(&lastLine); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode build output: %w", err)
		}

		if lastLine.Error != "" {
			return nil, fmt.Errorf("build error: %s", lastLine.Error)
		}

		if lastLine.Stream != "" {
			fmt.Fprint(output, lastLine.Stream)
		}
	}

	return &BuildResult{
		Image:      opts.Tag,
		Digest:     lastLine.Aux.ID,
		ModuleType: moduleType,
	}, nil
}

// Close releases resources.
func (b *Builder) Close() error {
	return b.dockerClient.Close()
}

// generateDockerfile generates a Dockerfile for the given module type.
func generateDockerfile(moduleType ModuleType, moduleDir string) (string, error) {
	switch moduleType {
	case ModuleTypePulumi:
		return generatePulumiDockerfile(moduleDir)
	case ModuleTypeOpenTofu:
		return generateOpenTofuDockerfile(moduleDir)
	default:
		return "", fmt.Errorf("unsupported module type: %s", moduleType)
	}
}

// generatePulumiDockerfile generates a Dockerfile for a Pulumi module.
func generatePulumiDockerfile(moduleDir string) (string, error) {
	// Read Pulumi.yaml to detect runtime
	pulumiYaml := filepath.Join(moduleDir, "Pulumi.yaml")
	data, err := os.ReadFile(pulumiYaml)
	if err != nil {
		return "", fmt.Errorf("failed to read Pulumi.yaml: %w", err)
	}

	// Simple runtime detection from Pulumi.yaml
	runtime := "nodejs" // default
	content := string(data)
	if strings.Contains(content, "runtime: python") || strings.Contains(content, "runtime:\n  name: python") {
		runtime = "python"
	} else if strings.Contains(content, "runtime: go") || strings.Contains(content, "runtime:\n  name: go") {
		runtime = "go"
	} else if strings.Contains(content, "runtime: dotnet") || strings.Contains(content, "runtime:\n  name: dotnet") {
		runtime = "dotnet"
	}

	// Check if package.json or requirements.txt exists
	hasPackageJson := fileExists(filepath.Join(moduleDir, "package.json"))
	hasRequirements := fileExists(filepath.Join(moduleDir, "requirements.txt"))
	hasGoMod := fileExists(filepath.Join(moduleDir, "go.mod"))

	var dockerfile strings.Builder

	dockerfile.WriteString(`# Auto-generated Dockerfile for Pulumi module
# This image bundles the Pulumi CLI with the module code

`)

	switch runtime {
	case "nodejs":
		dockerfile.WriteString(`FROM pulumi/pulumi-nodejs:latest

WORKDIR /app

# Copy module files
COPY . .

`)
		if hasPackageJson {
			dockerfile.WriteString(`# Install dependencies
RUN npm ci --production
`)
		}

	case "python":
		dockerfile.WriteString(`FROM pulumi/pulumi-python:latest

WORKDIR /app

# Copy module files
COPY . .

`)
		if hasRequirements {
			dockerfile.WriteString(`# Install dependencies
RUN pip install -r requirements.txt
`)
		}

	case "go":
		dockerfile.WriteString(`FROM pulumi/pulumi-go:latest

WORKDIR /app

# Copy module files
COPY . .

`)
		if hasGoMod {
			dockerfile.WriteString(`# Download dependencies
RUN go mod download

# Build the module
RUN go build -o /app/module .
`)
		}

	case "dotnet":
		dockerfile.WriteString(`FROM pulumi/pulumi-dotnet:latest

WORKDIR /app

# Copy module files
COPY . .

# Restore dependencies
RUN dotnet restore
`)
	}

	// Set entrypoint to the Pulumi CLI
	dockerfile.WriteString(`
ENTRYPOINT ["pulumi"]
`)

	return dockerfile.String(), nil
}

// generateOpenTofuDockerfile generates a Dockerfile for an OpenTofu module.
// Uses a multi-stage build: the minimal OpenTofu image provides the tofu binary,
// and Alpine provides the runtime environment (see https://opentofu.org/docs/intro/install/docker/).
func generateOpenTofuDockerfile(moduleDir string) (string, error) {
	return `# Auto-generated Dockerfile for OpenTofu module
# Uses multi-stage build per OpenTofu 1.10+ requirements

FROM ghcr.io/opentofu/opentofu:minimal AS tofu

FROM alpine:3.20

# Install the tofu binary from the minimal image
COPY --from=tofu /usr/local/bin/tofu /usr/local/bin/tofu

# Install common utilities needed by providers
RUN apk add --no-cache git curl ca-certificates

WORKDIR /app

# Copy module files
COPY . .

# Initialize the module (download providers and lock versions)
RUN tofu init -backend=false

ENTRYPOINT ["tofu"]
`, nil
}

// createBuildContext creates a tar archive for the Docker build context.
func createBuildContext(moduleDir string, dockerfile string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	// Add Dockerfile
	dockerfileBytes := []byte(dockerfile)
	if err := tw.WriteHeader(&tar.Header{
		Name: "Dockerfile",
		Mode: 0644,
		Size: int64(len(dockerfileBytes)),
	}); err != nil {
		return nil, err
	}
	if _, err := tw.Write(dockerfileBytes); err != nil {
		return nil, err
	}

	// Walk the module directory and add files
	err := filepath.Walk(moduleDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files/directories and common excludes
		name := info.Name()
		if strings.HasPrefix(name, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip node_modules, __pycache__, etc.
		if info.IsDir() {
			if name == "node_modules" || name == "__pycache__" || name == ".terraform" || name == ".pulumi" {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(moduleDir, path)
		if err != nil {
			return err
		}

		// Read file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Add to tar
		header := &tar.Header{
			Name: relPath,
			Mode: int64(info.Mode()),
			Size: int64(len(data)),
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if _, err := tw.Write(data); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return &buf, nil
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
