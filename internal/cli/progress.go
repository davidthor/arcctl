package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

// ResourceStatus represents the current status of a resource.
type ResourceStatus string

const (
	StatusPending    ResourceStatus = "pending"
	StatusWaiting    ResourceStatus = "waiting"
	StatusInProgress ResourceStatus = "in_progress"
	StatusCompleted  ResourceStatus = "completed"
	StatusFailed     ResourceStatus = "failed"
	StatusSkipped    ResourceStatus = "skipped"
)

// ResourceInfo holds information about a resource for progress tracking.
type ResourceInfo struct {
	Name         string
	Type         string
	Component    string
	Status       ResourceStatus
	Dependencies []string
	StartTime    time.Time
	EndTime      time.Time
	Error        error
	Message      string
}

// ProgressTable displays and updates deployment progress.
type ProgressTable struct {
	mu           sync.Mutex
	resources    map[string]*ResourceInfo
	order        []string // Maintains insertion order for display
	writer       io.Writer
	isTTY        bool
	lastRender   int // Number of lines in last render
	startTime    time.Time
	headerPrinted bool
}

// NewProgressTable creates a new progress table.
func NewProgressTable(w io.Writer) *ProgressTable {
	isTTY := false
	if f, ok := w.(*os.File); ok {
		isTTY = term.IsTerminal(int(f.Fd()))
	}

	return &ProgressTable{
		resources: make(map[string]*ResourceInfo),
		order:     []string{},
		writer:    w,
		isTTY:     isTTY,
		startTime: time.Now(),
	}
}

// AddResource adds a resource to track.
func (p *ProgressTable) AddResource(id, name, resourceType, component string, dependencies []string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.resources[id]; !exists {
		p.order = append(p.order, id)
	}

	status := StatusPending
	if len(dependencies) > 0 {
		status = StatusWaiting
	}

	p.resources[id] = &ResourceInfo{
		Name:         name,
		Type:         resourceType,
		Component:    component,
		Status:       status,
		Dependencies: dependencies,
	}
}

// UpdateStatus updates the status of a resource.
func (p *ProgressTable) UpdateStatus(id string, status ResourceStatus, message string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if res, ok := p.resources[id]; ok {
		res.Status = status
		res.Message = message

		if status == StatusInProgress && res.StartTime.IsZero() {
			res.StartTime = time.Now()
		}
		if status == StatusCompleted || status == StatusFailed || status == StatusSkipped {
			res.EndTime = time.Now()
		}
	}

	p.render()
}

// SetError sets an error for a resource.
func (p *ProgressTable) SetError(id string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if res, ok := p.resources[id]; ok {
		res.Status = StatusFailed
		res.Error = err
		res.EndTime = time.Now()
	}

	p.render()
}

// Render displays the current state of all resources.
func (p *ProgressTable) Render() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.render()
}

func (p *ProgressTable) render() {
	if p.isTTY {
		p.renderTTY()
	} else {
		p.renderStream()
	}
}

func (p *ProgressTable) renderTTY() {
	// Clear previous output
	if p.lastRender > 0 {
		// Move cursor up and clear lines
		for i := 0; i < p.lastRender; i++ {
			fmt.Fprintf(p.writer, "\033[A\033[K")
		}
	}

	lines := p.buildTable()
	for _, line := range lines {
		fmt.Fprintln(p.writer, line)
	}
	p.lastRender = len(lines)
}

func (p *ProgressTable) renderStream() {
	// For non-TTY, only print changes (status updates)
	// This is handled by UpdateStatus printing individual lines
}

// PrintInitial prints the initial table state (for both TTY and non-TTY).
func (p *ProgressTable) PrintInitial() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.headerPrinted {
		fmt.Fprintln(p.writer)
		fmt.Fprintln(p.writer, "Resource Deployment Progress:")
		fmt.Fprintln(p.writer, strings.Repeat("─", 80))
		p.headerPrinted = true
	}

	if p.isTTY {
		lines := p.buildTable()
		for _, line := range lines {
			fmt.Fprintln(p.writer, line)
		}
		p.lastRender = len(lines)
	} else {
		// For non-TTY, print the resource list with initial status
		for _, id := range p.order {
			res := p.resources[id]
			deps := ""
			if len(res.Dependencies) > 0 {
				depNames := p.getDependencyNames(res.Dependencies)
				deps = fmt.Sprintf(" (depends on: %s)", strings.Join(depNames, ", "))
			}
			fmt.Fprintf(p.writer, "  %-12s %-30s %s%s\n",
				formatResourceType(res.Type),
				res.Name,
				p.statusIcon(res.Status),
				deps)
		}
	}
}

// PrintUpdate prints a status update for a resource (used in non-TTY mode).
func (p *ProgressTable) PrintUpdate(id string) {
	if p.isTTY {
		return // TTY mode uses full table re-render
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	res, ok := p.resources[id]
	if !ok {
		return
	}

	var statusStr string
	switch res.Status {
	case StatusInProgress:
		statusStr = fmt.Sprintf("  %s Starting %s %q...", p.statusIcon(res.Status), res.Type, res.Name)
	case StatusCompleted:
		duration := ""
		if !res.EndTime.IsZero() && !res.StartTime.IsZero() {
			duration = fmt.Sprintf(" (%s)", res.EndTime.Sub(res.StartTime).Round(time.Millisecond))
		}
		statusStr = fmt.Sprintf("  %s Completed %s %q%s", p.statusIcon(res.Status), res.Type, res.Name, duration)
		if res.Message != "" {
			statusStr += fmt.Sprintf(" - %s", res.Message)
		}
	case StatusFailed:
		statusStr = fmt.Sprintf("  %s Failed %s %q", p.statusIcon(res.Status), res.Type, res.Name)
		if res.Error != nil {
			statusStr += fmt.Sprintf(": %v", res.Error)
		}
	case StatusSkipped:
		statusStr = fmt.Sprintf("  %s Skipped %s %q", p.statusIcon(res.Status), res.Type, res.Name)
		if res.Message != "" {
			statusStr += fmt.Sprintf(" - %s", res.Message)
		}
	default:
		return // Don't print pending/waiting updates in stream mode
	}

	fmt.Fprintln(p.writer, statusStr)
}

func (p *ProgressTable) buildTable() []string {
	var lines []string

	// Calculate column widths
	maxType := 10
	maxName := 20
	maxStatus := 12

	for _, id := range p.order {
		res := p.resources[id]
		if len(res.Type) > maxType {
			maxType = len(res.Type)
		}
		if len(res.Name) > maxName {
			maxName = len(res.Name)
		}
	}

	// Cap maximums for readability
	if maxType > 15 {
		maxType = 15
	}
	if maxName > 35 {
		maxName = 35
	}

	// Header
	header := fmt.Sprintf("  %-*s  %-*s  %-*s  %s",
		maxType, "TYPE",
		maxName, "NAME",
		maxStatus, "STATUS",
		"DEPENDENCIES")
	lines = append(lines, header)
	lines = append(lines, "  "+strings.Repeat("─", maxType+maxName+maxStatus+30))

	// Resources
	for _, id := range p.order {
		res := p.resources[id]

		// Truncate if needed
		typeName := res.Type
		if len(typeName) > maxType {
			typeName = typeName[:maxType-2] + ".."
		}
		name := res.Name
		if len(name) > maxName {
			name = name[:maxName-2] + ".."
		}

		statusStr := p.formatStatus(res)

		// Dependencies column
		deps := "-"
		if len(res.Dependencies) > 0 {
			depNames := p.getDependencyNames(res.Dependencies)
			deps = strings.Join(depNames, ", ")
			if len(deps) > 30 {
				deps = deps[:27] + "..."
			}
		}

		line := fmt.Sprintf("  %-*s  %-*s  %-*s  %s",
			maxType, typeName,
			maxName, name,
			maxStatus, statusStr,
			deps)
		lines = append(lines, line)
	}

	// Summary line
	lines = append(lines, "  "+strings.Repeat("─", maxType+maxName+maxStatus+30))
	summary := p.buildSummary()
	lines = append(lines, "  "+summary)

	return lines
}

func (p *ProgressTable) formatStatus(res *ResourceInfo) string {
	icon := p.statusIcon(res.Status)
	var status string

	switch res.Status {
	case StatusPending:
		status = "Pending"
	case StatusWaiting:
		status = "Waiting"
	case StatusInProgress:
		elapsed := time.Since(res.StartTime).Round(time.Second)
		status = fmt.Sprintf("Running %s", elapsed)
	case StatusCompleted:
		status = "Done"
	case StatusFailed:
		status = "Failed"
	case StatusSkipped:
		status = "Skipped"
	}

	return fmt.Sprintf("%s %s", icon, status)
}

func (p *ProgressTable) statusIcon(status ResourceStatus) string {
	switch status {
	case StatusPending:
		return "○"
	case StatusWaiting:
		return "◔"
	case StatusInProgress:
		return "◐"
	case StatusCompleted:
		return "●"
	case StatusFailed:
		return "✗"
	case StatusSkipped:
		return "◌"
	default:
		return "?"
	}
}

func (p *ProgressTable) getDependencyNames(depIDs []string) []string {
	names := make([]string, 0, len(depIDs))
	for _, depID := range depIDs {
		if res, ok := p.resources[depID]; ok {
			names = append(names, res.Name)
		} else {
			// Extract name from ID (format: component/type/name)
			parts := strings.Split(depID, "/")
			if len(parts) >= 3 {
				names = append(names, parts[len(parts)-1])
			} else {
				names = append(names, depID)
			}
		}
	}
	return names
}

func (p *ProgressTable) buildSummary() string {
	var pending, waiting, running, completed, failed, skipped int

	for _, res := range p.resources {
		switch res.Status {
		case StatusPending:
			pending++
		case StatusWaiting:
			waiting++
		case StatusInProgress:
			running++
		case StatusCompleted:
			completed++
		case StatusFailed:
			failed++
		case StatusSkipped:
			skipped++
		}
	}

	total := len(p.resources)
	elapsed := time.Since(p.startTime).Round(time.Second)

	parts := []string{}
	if completed > 0 {
		parts = append(parts, fmt.Sprintf("● %d completed", completed))
	}
	if running > 0 {
		parts = append(parts, fmt.Sprintf("◐ %d running", running))
	}
	if waiting > 0 {
		parts = append(parts, fmt.Sprintf("◔ %d waiting", waiting))
	}
	if pending > 0 {
		parts = append(parts, fmt.Sprintf("○ %d pending", pending))
	}
	if failed > 0 {
		parts = append(parts, fmt.Sprintf("✗ %d failed", failed))
	}
	if skipped > 0 {
		parts = append(parts, fmt.Sprintf("◌ %d skipped", skipped))
	}

	return fmt.Sprintf("%s  |  %d/%d resources  |  %s elapsed",
		strings.Join(parts, "  "),
		completed+failed+skipped,
		total,
		elapsed)
}

// PrintFinalSummary prints the final deployment summary.
func (p *ProgressTable) PrintFinalSummary() {
	p.mu.Lock()
	defer p.mu.Unlock()

	var completed, failed, skipped int
	for _, res := range p.resources {
		switch res.Status {
		case StatusCompleted:
			completed++
		case StatusFailed:
			failed++
		case StatusSkipped:
			skipped++
		}
	}

	elapsed := time.Since(p.startTime).Round(time.Millisecond)

	fmt.Fprintln(p.writer)
	fmt.Fprintln(p.writer, strings.Repeat("─", 80))

	if failed > 0 {
		fmt.Fprintf(p.writer, "Deployment completed with errors in %s\n", elapsed)
		fmt.Fprintf(p.writer, "  ● %d succeeded, ✗ %d failed, ◌ %d skipped\n", completed, failed, skipped)

		// List failed resources
		fmt.Fprintln(p.writer, "\nFailed resources:")
		for _, id := range p.order {
			res := p.resources[id]
			if res.Status == StatusFailed {
				fmt.Fprintf(p.writer, "  ✗ %s %q", res.Type, res.Name)
				if res.Error != nil {
					fmt.Fprintf(p.writer, ": %v", res.Error)
				}
				fmt.Fprintln(p.writer)
			}
		}
	} else {
		fmt.Fprintf(p.writer, "Deployment completed successfully in %s\n", elapsed)
		fmt.Fprintf(p.writer, "  ● %d resources deployed\n", completed)
	}
}

// GetCompletedCount returns the number of completed resources.
func (p *ProgressTable) GetCompletedCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	count := 0
	for _, res := range p.resources {
		if res.Status == StatusCompleted {
			count++
		}
	}
	return count
}

// GetFailedCount returns the number of failed resources.
func (p *ProgressTable) GetFailedCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	count := 0
	for _, res := range p.resources {
		if res.Status == StatusFailed {
			count++
		}
	}
	return count
}

// HasPending returns true if there are pending or waiting resources.
func (p *ProgressTable) HasPending() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, res := range p.resources {
		if res.Status == StatusPending || res.Status == StatusWaiting || res.Status == StatusInProgress {
			return true
		}
	}
	return false
}

// GetResourcesByStatus returns all resource IDs with the given status.
func (p *ProgressTable) GetResourcesByStatus(status ResourceStatus) []string {
	p.mu.Lock()
	defer p.mu.Unlock()

	var ids []string
	for _, id := range p.order {
		if p.resources[id].Status == status {
			ids = append(ids, id)
		}
	}
	return ids
}

// CheckDependencies updates waiting resources to pending if their dependencies are met.
func (p *ProgressTable) CheckDependencies() []string {
	p.mu.Lock()
	defer p.mu.Unlock()

	var ready []string

	for _, id := range p.order {
		res := p.resources[id]
		if res.Status != StatusWaiting {
			continue
		}

		allMet := true
		for _, depID := range res.Dependencies {
			if dep, ok := p.resources[depID]; ok {
				if dep.Status != StatusCompleted {
					allMet = false
					break
				}
			}
		}

		if allMet {
			res.Status = StatusPending
			ready = append(ready, id)
		}
	}

	return ready
}

// formatResourceType returns a formatted resource type string.
func formatResourceType(t string) string {
	// Capitalize first letter and format nicely
	if t == "" {
		return "Unknown"
	}
	return strings.ToUpper(t[:1]) + strings.ToLower(t[1:])
}

// SortedResourceIDs returns resource IDs sorted by dependency order.
func (p *ProgressTable) SortedResourceIDs() []string {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Topological sort
	visited := make(map[string]bool)
	var result []string

	var visit func(id string)
	visit = func(id string) {
		if visited[id] {
			return
		}
		visited[id] = true

		res := p.resources[id]
		for _, depID := range res.Dependencies {
			visit(depID)
		}
		result = append(result, id)
	}

	// Sort order for determinism
	ids := make([]string, len(p.order))
	copy(ids, p.order)
	sort.Strings(ids)

	for _, id := range ids {
		visit(id)
	}

	return result
}
