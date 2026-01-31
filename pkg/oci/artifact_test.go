package oci

import (
	"testing"
)

func TestParseReference(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Reference
	}{
		{
			name:  "simple image name",
			input: "nginx",
			expected: Reference{
				Registry:   "docker.io",
				Repository: "library/nginx",
				Tag:        "latest",
			},
		},
		{
			name:  "image with tag",
			input: "nginx:1.21",
			expected: Reference{
				Registry:   "docker.io",
				Repository: "library/nginx",
				Tag:        "1.21",
			},
		},
		{
			name:  "image with org",
			input: "myorg/myapp",
			expected: Reference{
				Registry:   "docker.io",
				Repository: "myorg/myapp",
				Tag:        "latest",
			},
		},
		{
			name:  "image with org and tag",
			input: "myorg/myapp:v1.0.0",
			expected: Reference{
				Registry:   "docker.io",
				Repository: "myorg/myapp",
				Tag:        "v1.0.0",
			},
		},
		{
			name:  "ghcr.io registry",
			input: "ghcr.io/myorg/myapp:latest",
			expected: Reference{
				Registry:   "ghcr.io",
				Repository: "myorg/myapp",
				Tag:        "latest",
			},
		},
		{
			name:  "registry with port",
			input: "localhost:5000/myapp:v1",
			expected: Reference{
				Registry:   "localhost:5000",
				Repository: "myapp",
				Tag:        "v1",
			},
		},
		{
			name:  "registry with nested path",
			input: "gcr.io/my-project/my-image:tag",
			expected: Reference{
				Registry:   "gcr.io",
				Repository: "my-project/my-image",
				Tag:        "tag",
			},
		},
		{
			name:  "image with digest",
			input: "nginx@sha256:abc123",
			expected: Reference{
				Registry:   "docker.io",
				Repository: "library/nginx",
				Digest:     "sha256:abc123",
			},
		},
		{
			name:  "image with tag and digest",
			input: "nginx:latest@sha256:abc123",
			expected: Reference{
				Registry:   "docker.io",
				Repository: "library/nginx",
				Tag:        "latest",
				Digest:     "sha256:abc123",
			},
		},
		{
			name:  "full reference with registry, tag and digest",
			input: "ghcr.io/myorg/myapp:v1.0.0@sha256:abc123def456",
			expected: Reference{
				Registry:   "ghcr.io",
				Repository: "myorg/myapp",
				Tag:        "v1.0.0",
				Digest:     "sha256:abc123def456",
			},
		},
		{
			name:  "localhost registry",
			input: "localhost/myapp:test",
			expected: Reference{
				Registry:   "localhost",
				Repository: "myapp",
				Tag:        "test",
			},
		},
		{
			name:  "AWS ECR reference",
			input: "123456789012.dkr.ecr.us-west-2.amazonaws.com/my-repo:latest",
			expected: Reference{
				Registry:   "123456789012.dkr.ecr.us-west-2.amazonaws.com",
				Repository: "my-repo",
				Tag:        "latest",
			},
		},
		{
			name:  "Azure ACR reference",
			input: "myregistry.azurecr.io/samples/hello-world:v1",
			expected: Reference{
				Registry:   "myregistry.azurecr.io",
				Repository: "samples/hello-world",
				Tag:        "v1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseReference(tt.input)
			if err != nil {
				t.Fatalf("ParseReference(%q) returned error: %v", tt.input, err)
			}

			if result.Registry != tt.expected.Registry {
				t.Errorf("Registry: got %q, want %q", result.Registry, tt.expected.Registry)
			}
			if result.Repository != tt.expected.Repository {
				t.Errorf("Repository: got %q, want %q", result.Repository, tt.expected.Repository)
			}
			if result.Tag != tt.expected.Tag {
				t.Errorf("Tag: got %q, want %q", result.Tag, tt.expected.Tag)
			}
			if result.Digest != tt.expected.Digest {
				t.Errorf("Digest: got %q, want %q", result.Digest, tt.expected.Digest)
			}
		})
	}
}

func TestReferenceString(t *testing.T) {
	tests := []struct {
		name     string
		ref      Reference
		expected string
	}{
		{
			name: "basic reference",
			ref: Reference{
				Registry:   "docker.io",
				Repository: "library/nginx",
				Tag:        "latest",
			},
			expected: "docker.io/library/nginx:latest",
		},
		{
			name: "reference with digest only",
			ref: Reference{
				Registry:   "ghcr.io",
				Repository: "myorg/myapp",
				Digest:     "sha256:abc123",
			},
			expected: "ghcr.io/myorg/myapp@sha256:abc123",
		},
		{
			name: "reference with tag and digest",
			ref: Reference{
				Registry:   "gcr.io",
				Repository: "project/image",
				Tag:        "v1.0.0",
				Digest:     "sha256:def456",
			},
			expected: "gcr.io/project/image:v1.0.0@sha256:def456",
		},
		{
			name: "reference with no tag or digest",
			ref: Reference{
				Registry:   "docker.io",
				Repository: "myapp",
			},
			expected: "docker.io/myapp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ref.String()
			if result != tt.expected {
				t.Errorf("String(): got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestArtifactTypes(t *testing.T) {
	// Verify artifact type constants
	if ArtifactTypeComponent != "component" {
		t.Errorf("ArtifactTypeComponent: got %q, want %q", ArtifactTypeComponent, "component")
	}
	if ArtifactTypeDatacenter != "datacenter" {
		t.Errorf("ArtifactTypeDatacenter: got %q, want %q", ArtifactTypeDatacenter, "datacenter")
	}
	if ArtifactTypeModule != "module" {
		t.Errorf("ArtifactTypeModule: got %q, want %q", ArtifactTypeModule, "module")
	}
}

func TestMediaTypes(t *testing.T) {
	// Verify media type constants follow expected pattern
	expectedPrefix := "application/vnd.architect."

	mediaTypes := []struct {
		name      string
		mediaType string
	}{
		{"ComponentConfig", MediaTypeComponentConfig},
		{"ComponentLayer", MediaTypeComponentLayer},
		{"DatacenterConfig", MediaTypeDatacenterConfig},
		{"DatacenterLayer", MediaTypeDatacenterLayer},
		{"ModuleConfig", MediaTypeModuleConfig},
		{"ModuleLayer", MediaTypeModuleLayer},
	}

	for _, mt := range mediaTypes {
		t.Run(mt.name, func(t *testing.T) {
			if len(mt.mediaType) < len(expectedPrefix) {
				t.Errorf("%s: media type too short: %q", mt.name, mt.mediaType)
				return
			}
			if mt.mediaType[:len(expectedPrefix)] != expectedPrefix {
				t.Errorf("%s: expected prefix %q, got %q", mt.name, expectedPrefix, mt.mediaType[:len(expectedPrefix)])
			}
		})
	}
}
