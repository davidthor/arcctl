package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateStateManagerWithConfig_Defaults(t *testing.T) {
	// Clear any env vars that might interfere
	os.Unsetenv(EnvStateBackend)

	// Test with no CLI flags - should use local backend with ~/.cldctl/state
	mgr, err := createStateManagerWithConfig("", nil)
	require.NoError(t, err)
	assert.NotNil(t, mgr)
}

func TestCreateStateManagerWithConfig_ExplicitLocal(t *testing.T) {
	// Test with explicit local backend type
	mgr, err := createStateManagerWithConfig("local", nil)
	require.NoError(t, err)
	assert.NotNil(t, mgr)
}

func TestCreateStateManagerWithConfig_CustomPath(t *testing.T) {
	// Create a temp directory for the state
	tempDir := t.TempDir()
	statePath := filepath.Join(tempDir, "test-state")

	// Test with custom path via CLI backend config
	mgr, err := createStateManagerWithConfig("local", []string{"path=" + statePath})
	require.NoError(t, err)
	assert.NotNil(t, mgr)

	// Verify the directory was created
	_, err = os.Stat(statePath)
	assert.NoError(t, err)
}

func TestCreateStateManagerWithConfig_EnvVarBackend(t *testing.T) {
	// Set backend via environment variable
	os.Setenv(EnvStateBackend, "local")
	defer os.Unsetenv(EnvStateBackend)

	mgr, err := createStateManagerWithConfig("", nil)
	require.NoError(t, err)
	assert.NotNil(t, mgr)
}

func TestCreateStateManagerWithConfig_EnvVarPath(t *testing.T) {
	// Create a temp directory for the state
	tempDir := t.TempDir()
	statePath := filepath.Join(tempDir, "env-state")

	// Set path via environment variable
	os.Setenv(EnvStatePrefix+"PATH", statePath)
	defer os.Unsetenv(EnvStatePrefix + "PATH")

	mgr, err := createStateManagerWithConfig("", nil)
	require.NoError(t, err)
	assert.NotNil(t, mgr)

	// Verify the directory was created
	_, err = os.Stat(statePath)
	assert.NoError(t, err)
}

func TestCreateStateManagerWithConfig_CLIOverridesEnvVar(t *testing.T) {
	// Create temp directories
	tempDir := t.TempDir()
	envPath := filepath.Join(tempDir, "env-path")
	cliPath := filepath.Join(tempDir, "cli-path")

	// Set path via environment variable
	os.Setenv(EnvStatePrefix+"PATH", envPath)
	defer os.Unsetenv(EnvStatePrefix + "PATH")

	// CLI flag should override env var
	mgr, err := createStateManagerWithConfig("local", []string{"path=" + cliPath})
	require.NoError(t, err)
	assert.NotNil(t, mgr)

	// CLI path should be used, not env path
	_, err = os.Stat(cliPath)
	assert.NoError(t, err)
}
