package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProgressTable(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	assert.NotNil(t, pt)
	assert.NotNil(t, pt.resources)
	assert.Equal(t, 0, len(pt.order))
}

func TestProgressTable_AddResource(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)

	assert.Equal(t, 1, len(pt.resources))
	assert.Equal(t, "main", pt.resources["comp/database/main"].Name)
	assert.Equal(t, "database", pt.resources["comp/database/main"].Type)
	assert.Equal(t, StatusPending, pt.resources["comp/database/main"].Status)
}

func TestProgressTable_AddResourceWithDependencies(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.AddResource("comp/deployment/api", "api", "deployment", "comp", []string{"comp/database/main"})

	assert.Equal(t, 2, len(pt.resources))
	// Resources with dependencies start in "waiting" status
	assert.Equal(t, StatusWaiting, pt.resources["comp/deployment/api"].Status)
	assert.Equal(t, []string{"comp/database/main"}, pt.resources["comp/deployment/api"].Dependencies)
}

func TestProgressTable_UpdateStatus(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.UpdateStatus("comp/database/main", StatusInProgress, "provisioning")

	assert.Equal(t, StatusInProgress, pt.resources["comp/database/main"].Status)
	assert.Equal(t, "provisioning", pt.resources["comp/database/main"].Message)
	assert.False(t, pt.resources["comp/database/main"].StartTime.IsZero())
}

func TestProgressTable_SetError(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.SetError("comp/database/main", assert.AnError)

	assert.Equal(t, StatusFailed, pt.resources["comp/database/main"].Status)
	assert.Equal(t, assert.AnError, pt.resources["comp/database/main"].Error)
}

func TestProgressTable_GetCompletedCount(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.AddResource("comp/database/cache", "cache", "database", "comp", nil)
	pt.UpdateStatus("comp/database/main", StatusCompleted, "")

	assert.Equal(t, 1, pt.GetCompletedCount())
}

func TestProgressTable_GetFailedCount(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.AddResource("comp/database/cache", "cache", "database", "comp", nil)
	pt.SetError("comp/database/main", assert.AnError)

	assert.Equal(t, 1, pt.GetFailedCount())
}

func TestProgressTable_HasPending(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	assert.True(t, pt.HasPending())

	pt.UpdateStatus("comp/database/main", StatusCompleted, "")
	assert.False(t, pt.HasPending())
}

func TestProgressTable_CheckDependencies(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.AddResource("comp/deployment/api", "api", "deployment", "comp", []string{"comp/database/main"})

	// api should be waiting, not ready
	assert.Equal(t, StatusWaiting, pt.resources["comp/deployment/api"].Status)
	ready := pt.CheckDependencies()
	assert.Empty(t, ready)

	// After database completes, api should become pending
	pt.UpdateStatus("comp/database/main", StatusCompleted, "")
	ready = pt.CheckDependencies()
	assert.Contains(t, ready, "comp/deployment/api")
	assert.Equal(t, StatusPending, pt.resources["comp/deployment/api"].Status)
}

func TestProgressTable_PrintInitial(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)
	pt.isTTY = false // Force non-TTY mode for predictable output

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.AddResource("comp/deployment/api", "api", "deployment", "comp", []string{"comp/database/main"})
	pt.PrintInitial()

	output := buf.String()
	assert.Contains(t, output, "Resource Deployment Progress")
	assert.Contains(t, output, "main")
	assert.Contains(t, output, "api")
}

func TestProgressTable_PrintFinalSummary_Success(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)
	pt.isTTY = false

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.UpdateStatus("comp/database/main", StatusCompleted, "")
	pt.PrintFinalSummary()

	output := buf.String()
	assert.Contains(t, output, "successfully")
	assert.Contains(t, output, "1 resources deployed")
}

func TestProgressTable_PrintFinalSummary_WithFailures(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)
	pt.isTTY = false

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.AddResource("comp/deployment/api", "api", "deployment", "comp", nil)
	pt.UpdateStatus("comp/database/main", StatusCompleted, "")
	pt.SetError("comp/deployment/api", assert.AnError)
	pt.PrintFinalSummary()

	output := buf.String()
	assert.Contains(t, output, "errors")
	assert.Contains(t, output, "1 succeeded")
	assert.Contains(t, output, "1 failed")
	assert.Contains(t, output, "Failed resources")
}

func TestStatusIcon(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	tests := []struct {
		status ResourceStatus
		want   string
	}{
		{StatusPending, "○"},
		{StatusWaiting, "◔"},
		{StatusInProgress, "◐"},
		{StatusCompleted, "●"},
		{StatusFailed, "✗"},
		{StatusSkipped, "◌"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			got := pt.statusIcon(tt.status)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatResourceType(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"database", "Database"},
		{"DEPLOYMENT", "Deployment"},
		{"function", "Function"},
		{"", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := formatResourceType(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResourceID(t *testing.T) {
	id := resourceID("mycomp", "database", "main")
	assert.Equal(t, "mycomp/database/main", id)
}

func TestProgressTable_BuildTable(t *testing.T) {
	buf := &bytes.Buffer{}
	pt := NewProgressTable(buf)

	pt.AddResource("comp/database/main", "main", "database", "comp", nil)
	pt.AddResource("comp/deployment/api", "api", "deployment", "comp", []string{"comp/database/main"})

	lines := pt.buildTable()

	// Should have header, separator, 2 resources, separator, summary
	assert.GreaterOrEqual(t, len(lines), 6)

	// Check header contains expected columns
	header := lines[0]
	assert.True(t, strings.Contains(header, "TYPE"))
	assert.True(t, strings.Contains(header, "NAME"))
	assert.True(t, strings.Contains(header, "STATUS"))
	assert.True(t, strings.Contains(header, "DEPENDENCIES"))
}
