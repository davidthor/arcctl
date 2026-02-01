package native

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDockerfileCmd_JSONFormat(t *testing.T) {
	// Create a temp Dockerfile
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")

	dockerfile := `FROM node:18
WORKDIR /app
COPY package.json .
RUN npm install
CMD ["npm", "start"]
`
	err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	require.NoError(t, err)

	cmd, err := ParseDockerfileCmd(dockerfilePath)
	require.NoError(t, err)
	assert.Equal(t, []string{"npm", "start"}, cmd)
}

func TestParseDockerfileCmd_ShellFormat(t *testing.T) {
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")

	dockerfile := `FROM node:18
WORKDIR /app
CMD npm run dev
`
	err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	require.NoError(t, err)

	cmd, err := ParseDockerfileCmd(dockerfilePath)
	require.NoError(t, err)
	assert.Equal(t, []string{"/bin/sh", "-c", "npm run dev"}, cmd)
}

func TestParseDockerfileCmd_LastCmdWins(t *testing.T) {
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")

	dockerfile := `FROM node:18
CMD ["npm", "install"]
CMD ["npm", "start"]
`
	err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	require.NoError(t, err)

	cmd, err := ParseDockerfileCmd(dockerfilePath)
	require.NoError(t, err)
	assert.Equal(t, []string{"npm", "start"}, cmd)
}

func TestParseDockerfileCmd_NoCmdFound(t *testing.T) {
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")

	dockerfile := `FROM node:18
WORKDIR /app
`
	err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	require.NoError(t, err)

	_, err = ParseDockerfileCmd(dockerfilePath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no CMD instruction found")
}

func TestParseDockerfileCmd_WithComments(t *testing.T) {
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")

	dockerfile := `FROM node:18
# This is a comment
# CMD ["npm", "test"]
CMD ["npm", "start"]
`
	err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	require.NoError(t, err)

	cmd, err := ParseDockerfileCmd(dockerfilePath)
	require.NoError(t, err)
	assert.Equal(t, []string{"npm", "start"}, cmd)
}

func TestExtractDockerfileCmdFromContext(t *testing.T) {
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")

	dockerfile := `FROM python:3.11
CMD ["python", "app.py"]
`
	err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	require.NoError(t, err)

	cmd, err := ExtractDockerfileCmdFromContext(tmpDir, "")
	require.NoError(t, err)
	assert.Equal(t, []string{"python", "app.py"}, cmd)
}

func TestExtractDockerfileCmdFromContext_CustomDockerfile(t *testing.T) {
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile.dev")

	dockerfile := `FROM node:18
CMD ["npm", "run", "dev"]
`
	err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	require.NoError(t, err)

	cmd, err := ExtractDockerfileCmdFromContext(tmpDir, "Dockerfile.dev")
	require.NoError(t, err)
	assert.Equal(t, []string{"npm", "run", "dev"}, cmd)
}

func TestEvaluateFunction_Coalesce(t *testing.T) {
	ctx := &EvalContext{
		Inputs: map[string]interface{}{
			"command": []interface{}{"npm", "start"},
			"empty":   "",
		},
	}

	// Test coalesce with valid value
	result, err := evaluateFunction("coalesce(inputs.command, inputs.empty)", ctx)
	require.NoError(t, err)
	assert.Equal(t, []interface{}{"npm", "start"}, result)

	// Test coalesce with empty string falling back
	result, err = evaluateFunction("coalesce(inputs.empty, inputs.command)", ctx)
	require.NoError(t, err)
	assert.Equal(t, []interface{}{"npm", "start"}, result)
}

func TestEvaluateFunction_DockerfileCmd(t *testing.T) {
	tmpDir := t.TempDir()
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")

	dockerfile := `FROM node:18
CMD ["node", "server.js"]
`
	err := os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	require.NoError(t, err)

	ctx := &EvalContext{
		Inputs: map[string]interface{}{
			"context": tmpDir,
		},
	}

	result, err := evaluateFunction("dockerfile_cmd(inputs.context)", ctx)
	require.NoError(t, err)
	assert.Equal(t, []string{"node", "server.js"}, result)
}

func TestSplitFunctionArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single arg",
			input:    "inputs.name",
			expected: []string{"inputs.name"},
		},
		{
			name:     "two args",
			input:    "inputs.command, inputs.default",
			expected: []string{"inputs.command", "inputs.default"},
		},
		{
			name:     "three args",
			input:    "a, b, c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "nested function",
			input:    "coalesce(a, b), c",
			expected: []string{"coalesce(a, b)", "c"},
		},
		{
			name:     "empty",
			input:    "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitFunctionArgs(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
