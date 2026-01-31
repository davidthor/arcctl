package iac

import (
	"bytes"
	"io"
	"testing"
)

func TestChangeAction_Constants(t *testing.T) {
	// Verify change action constants have expected values
	tests := []struct {
		action   ChangeAction
		expected string
	}{
		{ActionCreate, "create"},
		{ActionUpdate, "update"},
		{ActionDelete, "delete"},
		{ActionReplace, "replace"},
		{ActionNoop, "noop"},
	}

	for _, tc := range tests {
		if string(tc.action) != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, tc.action)
		}
	}
}

func TestRunOptions_Defaults(t *testing.T) {
	opts := RunOptions{}

	// Verify zero values
	if opts.ModuleSource != "" {
		t.Errorf("expected empty ModuleSource, got %q", opts.ModuleSource)
	}
	if opts.ModulePath != "" {
		t.Errorf("expected empty ModulePath, got %q", opts.ModulePath)
	}
	if opts.Inputs != nil {
		t.Error("expected nil Inputs")
	}
	if opts.StateReader != nil {
		t.Error("expected nil StateReader")
	}
	if opts.StateWriter != nil {
		t.Error("expected nil StateWriter")
	}
	if opts.WorkDir != "" {
		t.Errorf("expected empty WorkDir, got %q", opts.WorkDir)
	}
	if opts.Environment != nil {
		t.Error("expected nil Environment")
	}
	if opts.Volumes != nil {
		t.Error("expected nil Volumes")
	}
	if opts.Stdout != nil {
		t.Error("expected nil Stdout")
	}
	if opts.Stderr != nil {
		t.Error("expected nil Stderr")
	}
}

func TestRunOptions_WithValues(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stateReader := bytes.NewReader([]byte("state"))
	stateWriter := &bytes.Buffer{}

	opts := RunOptions{
		ModuleSource: "test-module",
		ModulePath:   "/path/to/module",
		Inputs: map[string]interface{}{
			"input1": "value1",
			"input2": 42,
		},
		StateReader: stateReader,
		StateWriter: stateWriter,
		WorkDir:     "/work/dir",
		Environment: map[string]string{
			"ENV_VAR": "value",
		},
		Volumes: []VolumeMount{
			{HostPath: "/host", MountPath: "/container", ReadOnly: true},
		},
		Stdout: &stdout,
		Stderr: &stderr,
	}

	if opts.ModuleSource != "test-module" {
		t.Errorf("expected ModuleSource 'test-module', got %q", opts.ModuleSource)
	}
	if opts.ModulePath != "/path/to/module" {
		t.Errorf("expected ModulePath '/path/to/module', got %q", opts.ModulePath)
	}
	if len(opts.Inputs) != 2 {
		t.Errorf("expected 2 inputs, got %d", len(opts.Inputs))
	}
	if opts.Inputs["input1"] != "value1" {
		t.Errorf("expected input1='value1', got %v", opts.Inputs["input1"])
	}
	if opts.Inputs["input2"] != 42 {
		t.Errorf("expected input2=42, got %v", opts.Inputs["input2"])
	}
	if opts.WorkDir != "/work/dir" {
		t.Errorf("expected WorkDir '/work/dir', got %q", opts.WorkDir)
	}
	if opts.Environment["ENV_VAR"] != "value" {
		t.Errorf("expected ENV_VAR='value', got %v", opts.Environment["ENV_VAR"])
	}
	if len(opts.Volumes) != 1 {
		t.Errorf("expected 1 volume, got %d", len(opts.Volumes))
	}
	if opts.Volumes[0].HostPath != "/host" {
		t.Errorf("expected HostPath '/host', got %q", opts.Volumes[0].HostPath)
	}
	if opts.Volumes[0].MountPath != "/container" {
		t.Errorf("expected MountPath '/container', got %q", opts.Volumes[0].MountPath)
	}
	if !opts.Volumes[0].ReadOnly {
		t.Error("expected ReadOnly to be true")
	}
}

func TestVolumeMount(t *testing.T) {
	vm := VolumeMount{
		HostPath:  "/host/path",
		MountPath: "/container/path",
		ReadOnly:  false,
	}

	if vm.HostPath != "/host/path" {
		t.Errorf("expected HostPath '/host/path', got %q", vm.HostPath)
	}
	if vm.MountPath != "/container/path" {
		t.Errorf("expected MountPath '/container/path', got %q", vm.MountPath)
	}
	if vm.ReadOnly {
		t.Error("expected ReadOnly to be false")
	}
}

func TestPreviewResult(t *testing.T) {
	result := PreviewResult{
		Changes: []ResourceChange{
			{
				ResourceID:   "resource1",
				ResourceType: "type1",
				Action:       ActionCreate,
			},
			{
				ResourceID:   "resource2",
				ResourceType: "type2",
				Action:       ActionUpdate,
			},
		},
		Summary: ChangeSummary{
			Create:  1,
			Update:  1,
			Delete:  0,
			Replace: 0,
		},
	}

	if len(result.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(result.Changes))
	}
	if result.Summary.Create != 1 {
		t.Errorf("expected Create=1, got %d", result.Summary.Create)
	}
	if result.Summary.Update != 1 {
		t.Errorf("expected Update=1, got %d", result.Summary.Update)
	}
}

func TestResourceChange(t *testing.T) {
	change := ResourceChange{
		ResourceID:   "test-resource",
		ResourceType: "test-type",
		Action:       ActionReplace,
		Before:       map[string]interface{}{"key": "old"},
		After:        map[string]interface{}{"key": "new"},
		Diff: []PropertyDiff{
			{Path: "key", OldValue: "old", NewValue: "new", Sensitive: false},
		},
	}

	if change.ResourceID != "test-resource" {
		t.Errorf("expected ResourceID 'test-resource', got %q", change.ResourceID)
	}
	if change.ResourceType != "test-type" {
		t.Errorf("expected ResourceType 'test-type', got %q", change.ResourceType)
	}
	if change.Action != ActionReplace {
		t.Errorf("expected Action 'replace', got %q", change.Action)
	}
	if len(change.Diff) != 1 {
		t.Errorf("expected 1 diff, got %d", len(change.Diff))
	}
}

func TestPropertyDiff(t *testing.T) {
	diff := PropertyDiff{
		Path:      "config.password",
		OldValue:  "old-secret",
		NewValue:  "new-secret",
		Sensitive: true,
	}

	if diff.Path != "config.password" {
		t.Errorf("expected Path 'config.password', got %q", diff.Path)
	}
	if diff.OldValue != "old-secret" {
		t.Errorf("expected OldValue 'old-secret', got %v", diff.OldValue)
	}
	if diff.NewValue != "new-secret" {
		t.Errorf("expected NewValue 'new-secret', got %v", diff.NewValue)
	}
	if !diff.Sensitive {
		t.Error("expected Sensitive to be true")
	}
}

func TestChangeSummary(t *testing.T) {
	summary := ChangeSummary{
		Create:  5,
		Update:  3,
		Delete:  2,
		Replace: 1,
	}

	if summary.Create != 5 {
		t.Errorf("expected Create=5, got %d", summary.Create)
	}
	if summary.Update != 3 {
		t.Errorf("expected Update=3, got %d", summary.Update)
	}
	if summary.Delete != 2 {
		t.Errorf("expected Delete=2, got %d", summary.Delete)
	}
	if summary.Replace != 1 {
		t.Errorf("expected Replace=1, got %d", summary.Replace)
	}
}

func TestApplyResult(t *testing.T) {
	result := ApplyResult{
		Outputs: map[string]OutputValue{
			"endpoint": {Value: "https://example.com", Sensitive: false},
			"password": {Value: "secret", Sensitive: true},
		},
		State: []byte(`{"version": 1}`),
	}

	if len(result.Outputs) != 2 {
		t.Errorf("expected 2 outputs, got %d", len(result.Outputs))
	}

	endpoint := result.Outputs["endpoint"]
	if endpoint.Value != "https://example.com" {
		t.Errorf("expected endpoint 'https://example.com', got %v", endpoint.Value)
	}
	if endpoint.Sensitive {
		t.Error("expected endpoint to not be sensitive")
	}

	password := result.Outputs["password"]
	if password.Value != "secret" {
		t.Errorf("expected password 'secret', got %v", password.Value)
	}
	if !password.Sensitive {
		t.Error("expected password to be sensitive")
	}

	if string(result.State) != `{"version": 1}` {
		t.Errorf("unexpected state: %s", result.State)
	}
}

func TestApplyResult_WithPartialError(t *testing.T) {
	result := ApplyResult{
		Outputs: map[string]OutputValue{
			"partial": {Value: "value", Sensitive: false},
		},
		State:        []byte(`{"partial": true}`),
		PartialError: io.EOF,
	}

	if result.PartialError == nil {
		t.Error("expected PartialError to be set")
	}
	if result.PartialError != io.EOF {
		t.Errorf("expected PartialError to be io.EOF, got %v", result.PartialError)
	}
}

func TestOutputValue(t *testing.T) {
	// Test non-sensitive output
	output := OutputValue{
		Value:     "test-value",
		Sensitive: false,
	}

	if output.Value != "test-value" {
		t.Errorf("expected Value 'test-value', got %v", output.Value)
	}
	if output.Sensitive {
		t.Error("expected Sensitive to be false")
	}

	// Test sensitive output
	sensitiveOutput := OutputValue{
		Value:     "secret-value",
		Sensitive: true,
	}

	if sensitiveOutput.Value != "secret-value" {
		t.Errorf("expected Value 'secret-value', got %v", sensitiveOutput.Value)
	}
	if !sensitiveOutput.Sensitive {
		t.Error("expected Sensitive to be true")
	}
}

func TestRefreshResult(t *testing.T) {
	result := RefreshResult{
		State: []byte(`{"refreshed": true}`),
		Drifts: []ResourceDrift{
			{
				ResourceID:   "drift-resource",
				ResourceType: "drift-type",
				Diffs: []PropertyDiff{
					{Path: "config", OldValue: "expected", NewValue: "actual", Sensitive: false},
				},
			},
		},
	}

	if string(result.State) != `{"refreshed": true}` {
		t.Errorf("unexpected state: %s", result.State)
	}

	if len(result.Drifts) != 1 {
		t.Errorf("expected 1 drift, got %d", len(result.Drifts))
	}

	drift := result.Drifts[0]
	if drift.ResourceID != "drift-resource" {
		t.Errorf("expected ResourceID 'drift-resource', got %q", drift.ResourceID)
	}
	if drift.ResourceType != "drift-type" {
		t.Errorf("expected ResourceType 'drift-type', got %q", drift.ResourceType)
	}
	if len(drift.Diffs) != 1 {
		t.Errorf("expected 1 diff in drift, got %d", len(drift.Diffs))
	}
}

func TestResourceDrift(t *testing.T) {
	drift := ResourceDrift{
		ResourceID:   "resource-id",
		ResourceType: "resource-type",
		Diffs: []PropertyDiff{
			{Path: "path1", OldValue: "old1", NewValue: "new1", Sensitive: false},
			{Path: "path2", OldValue: "old2", NewValue: "new2", Sensitive: true},
		},
	}

	if drift.ResourceID != "resource-id" {
		t.Errorf("expected ResourceID 'resource-id', got %q", drift.ResourceID)
	}
	if drift.ResourceType != "resource-type" {
		t.Errorf("expected ResourceType 'resource-type', got %q", drift.ResourceType)
	}
	if len(drift.Diffs) != 2 {
		t.Errorf("expected 2 diffs, got %d", len(drift.Diffs))
	}
}

func TestPlugin_Interface(t *testing.T) {
	// Verify the mock plugin implements the Plugin interface
	var _ Plugin = &mockPlugin{}
}
