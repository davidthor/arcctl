package output

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockHandler is a test handler that records events
type mockHandler struct {
	mu     sync.Mutex
	events []Event
	closed bool
}

func (m *mockHandler) HandleEvent(event Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
}

func (m *mockHandler) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

func (m *mockHandler) Events() []Event {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]Event{}, m.events...)
}

func (m *mockHandler) IsClosed() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closed
}

func TestNewStream(t *testing.T) {
	stream := NewStream()
	if stream == nil {
		t.Fatal("NewStream() returned nil")
	}
	defer stream.Close()

	if stream.handlers == nil {
		t.Error("handlers is nil")
	}
	if stream.events == nil {
		t.Error("events channel is nil")
	}
	if stream.done == nil {
		t.Error("done channel is nil")
	}
}

func TestStreamAddHandler(t *testing.T) {
	stream := NewStream()
	defer stream.Close()

	handler := &mockHandler{}
	stream.AddHandler(handler)

	if len(stream.handlers) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(stream.handlers))
	}
}

func TestStreamEmit(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	stream.Emit(Event{
		Level:     LevelInfo,
		Component: "test",
		Resource:  "resource",
		Message:   "test message",
	})

	// Allow time for async processing
	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].Message != "test message" {
		t.Errorf("Message: got %q, want %q", events[0].Message, "test message")
	}
	if events[0].Time.IsZero() {
		t.Error("Event time should be set automatically")
	}
}

func TestStreamEmitInfo(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	stream.EmitInfo("mycomp", "myres", "info message")

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	e := events[0]
	if e.Level != LevelInfo {
		t.Errorf("Level: got %d, want %d", e.Level, LevelInfo)
	}
	if e.Component != "mycomp" {
		t.Errorf("Component: got %q, want %q", e.Component, "mycomp")
	}
	if e.Resource != "myres" {
		t.Errorf("Resource: got %q, want %q", e.Resource, "myres")
	}
	if e.Message != "info message" {
		t.Errorf("Message: got %q, want %q", e.Message, "info message")
	}
	if e.Progress != -1 {
		t.Errorf("Progress: got %d, want %d", e.Progress, -1)
	}
}

func TestStreamEmitProgress(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	stream.EmitProgress("comp", "res", "deploying", 50)

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	e := events[0]
	if e.Progress != 50 {
		t.Errorf("Progress: got %d, want %d", e.Progress, 50)
	}
	if e.Action != "deploying" {
		t.Errorf("Action: got %q, want %q", e.Action, "deploying")
	}
}

func TestStreamEmitError(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	testErr := errors.New("something went wrong")
	stream.EmitError("comp", "res", testErr)

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	e := events[0]
	if e.Level != LevelError {
		t.Errorf("Level: got %d, want %d", e.Level, LevelError)
	}
	if e.Message != "something went wrong" {
		t.Errorf("Message: got %q, want %q", e.Message, "something went wrong")
	}
}

func TestStreamClose(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	stream.EmitInfo("comp", "res", "before close")
	time.Sleep(50 * time.Millisecond)

	err := stream.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	if !handler.IsClosed() {
		t.Error("Handler was not closed")
	}
}

func TestStreamMultipleHandlers(t *testing.T) {
	stream := NewStream()
	handler1 := &mockHandler{}
	handler2 := &mockHandler{}
	stream.AddHandler(handler1)
	stream.AddHandler(handler2)

	stream.EmitInfo("comp", "res", "broadcast message")

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events1 := handler1.Events()
	events2 := handler2.Events()

	if len(events1) != 1 {
		t.Errorf("Handler1: expected 1 event, got %d", len(events1))
	}
	if len(events2) != 1 {
		t.Errorf("Handler2: expected 1 event, got %d", len(events2))
	}
}

func TestStreamWriter(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	writer := stream.Writer("comp", "res", LevelInfo)

	// Write data with newlines
	_, _ = writer.Write([]byte("line 1\nline 2\n"))

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(events))
	}

	if events[0].Message != "line 1" {
		t.Errorf("Event 0 message: got %q, want %q", events[0].Message, "line 1")
	}
	if events[1].Message != "line 2" {
		t.Errorf("Event 1 message: got %q, want %q", events[1].Message, "line 2")
	}
}

func TestStreamWriterPartialLines(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	writer := stream.Writer("comp", "res", LevelInfo)

	// Write partial line
	_, _ = writer.Write([]byte("partial "))
	_, _ = writer.Write([]byte("line\n"))

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].Message != "partial line" {
		t.Errorf("Message: got %q, want %q", events[0].Message, "partial line")
	}
}

func TestConsoleHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(ConsoleOptions{
		Writer:    &buf,
		UseColors: false,
		Verbose:   false,
	})

	event := Event{
		Time:      time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Level:     LevelInfo,
		Component: "api",
		Resource:  "deployment",
		Message:   "Creating deployment",
		Progress:  -1,
	}

	handler.HandleEvent(event)

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("Output should contain INFO: %q", output)
	}
	if !strings.Contains(output, "api") {
		t.Errorf("Output should contain component: %q", output)
	}
	if !strings.Contains(output, "deployment") {
		t.Errorf("Output should contain resource: %q", output)
	}
	if !strings.Contains(output, "Creating deployment") {
		t.Errorf("Output should contain message: %q", output)
	}
}

func TestConsoleHandlerVerbose(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(ConsoleOptions{
		Writer:    &buf,
		UseColors: false,
		Verbose:   true,
	})

	event := Event{
		Time:      time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC),
		Level:     LevelDebug,
		Component: "api",
		Message:   "Debug message",
	}

	handler.HandleEvent(event)

	output := buf.String()
	if !strings.Contains(output, "10:30:45") {
		t.Errorf("Verbose output should contain time: %q", output)
	}
	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Verbose output should contain DEBUG: %q", output)
	}
}

func TestConsoleHandlerDebugSkippedInNonVerbose(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(ConsoleOptions{
		Writer:    &buf,
		UseColors: false,
		Verbose:   false,
	})

	event := Event{
		Level:   LevelDebug,
		Message: "Debug message",
	}

	handler.HandleEvent(event)

	if buf.Len() > 0 {
		t.Errorf("Debug events should be skipped in non-verbose mode: %q", buf.String())
	}
}

func TestConsoleHandlerProgress(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(ConsoleOptions{
		Writer:    &buf,
		UseColors: false,
	})

	event := Event{
		Level:    LevelInfo,
		Progress: 75,
		Message:  "Processing",
	}

	handler.HandleEvent(event)

	output := buf.String()
	if !strings.Contains(output, "75%") {
		t.Errorf("Output should contain progress percentage: %q", output)
	}
}

func TestConsoleHandlerColors(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(ConsoleOptions{
		Writer:    &buf,
		UseColors: true,
	})

	event := Event{
		Level:     LevelError,
		Component: "api",
		Message:   "Error occurred",
	}

	handler.HandleEvent(event)

	output := buf.String()
	// Check for ANSI color codes
	if !strings.Contains(output, "\033[") {
		t.Errorf("Colored output should contain ANSI codes: %q", output)
	}
}

func TestJSONHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := NewJSONHandler(&buf)

	event := Event{
		Time:      time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Level:     LevelInfo,
		Component: "api",
		Resource:  "deployment",
		Action:    "create",
		Message:   "Creating resource",
		Progress:  50,
	}

	handler.HandleEvent(event)

	output := buf.String()

	// Parse as JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v\nOutput: %q", err, output)
	}

	if parsed["level"] != "info" {
		t.Errorf("level: got %v, want %v", parsed["level"], "info")
	}
	if parsed["component"] != "api" {
		t.Errorf("component: got %v, want %v", parsed["component"], "api")
	}
	if parsed["message"] != "Creating resource" {
		t.Errorf("message: got %v, want %v", parsed["message"], "Creating resource")
	}
	if parsed["progress"].(float64) != 50 {
		t.Errorf("progress: got %v, want %v", parsed["progress"], 50)
	}
}

func TestJSONHandlerEscaping(t *testing.T) {
	var buf bytes.Buffer
	handler := NewJSONHandler(&buf)

	event := Event{
		Level:   LevelInfo,
		Message: "Message with \"quotes\" and\nnewline",
	}

	handler.HandleEvent(event)

	output := buf.String()

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v\nOutput: %q", err, output)
	}

	expected := "Message with \"quotes\" and\nnewline"
	if parsed["message"] != expected {
		t.Errorf("message: got %q, want %q", parsed["message"], expected)
	}
}

func TestLevelToString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelDebug, "debug"},
		{LevelInfo, "info"},
		{LevelWarn, "warn"},
		{LevelError, "error"},
		{Level(99), "unknown"},
	}

	for _, tt := range tests {
		result := levelToString(tt.level)
		if result != tt.expected {
			t.Errorf("levelToString(%d): got %q, want %q", tt.level, result, tt.expected)
		}
	}
}

func TestEscapeJSON(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello\"world", "hello\\\"world"},
		{"line1\nline2", "line1\\nline2"},
		{"tab\there", "tab\\there"},
		{"back\\slash", "back\\\\slash"},
		{"carriage\rreturn", "carriage\\rreturn"},
	}

	for _, tt := range tests {
		result := escapeJSON(tt.input)
		if result != tt.expected {
			t.Errorf("escapeJSON(%q): got %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestLineScanner(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	scanner := NewLineScanner(stream, "comp", "res", LevelInfo)

	reader := strings.NewReader("line 1\nline 2\nline 3\n")
	ctx := context.Background()

	err := scanner.Scan(ctx, reader)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 3 {
		t.Fatalf("Expected 3 events, got %d", len(events))
	}

	for i, expected := range []string{"line 1", "line 2", "line 3"} {
		if events[i].Message != expected {
			t.Errorf("Event %d message: got %q, want %q", i, events[i].Message, expected)
		}
	}
}

func TestLineScannerContextCancellation(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)
	defer stream.Close()

	scanner := NewLineScanner(stream, "comp", "res", LevelInfo)

	// Use a slow reader that checks context
	ctx, cancel := context.WithCancel(context.Background())

	// Create a reader that blocks until context is cancelled
	slowReader := &slowReader{ctx: ctx}

	// Start scanning in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- scanner.Scan(ctx, slowReader)
	}()

	// Give scanner time to start
	time.Sleep(10 * time.Millisecond)

	// Cancel context
	cancel()

	select {
	case err := <-errCh:
		// Scanner should return when reader stops (either nil or context error is acceptable)
		// The important thing is it doesn't hang
		_ = err
	case <-time.After(1 * time.Second):
		t.Error("Scanner did not return after context cancellation")
	}
}

// slowReader blocks until context is done
type slowReader struct {
	ctx context.Context
}

func (r *slowReader) Read(p []byte) (n int, err error) {
	<-r.ctx.Done()
	return 0, r.ctx.Err()
}

func TestProgressBar(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	bar := NewProgressBar(stream, "comp", "res", 100)

	bar.Increment()
	bar.Increment()
	bar.Add(8)

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 3 {
		t.Fatalf("Expected 3 events, got %d", len(events))
	}

	// Progress should be 1%, 2%, 10%
	expectedProgress := []int{1, 2, 10}
	for i, expected := range expectedProgress {
		if events[i].Progress != expected {
			t.Errorf("Event %d progress: got %d, want %d", i, events[i].Progress, expected)
		}
	}
}

func TestProgressBarSetCurrent(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	bar := NewProgressBar(stream, "comp", "res", 100)
	bar.SetCurrent(50)

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].Progress != 50 {
		t.Errorf("Progress: got %d, want %d", events[0].Progress, 50)
	}
}

func TestProgressBarComplete(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	bar := NewProgressBar(stream, "comp", "res", 100)
	bar.Complete()

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	if events[0].Progress != 100 {
		t.Errorf("Progress: got %d, want %d", events[0].Progress, 100)
	}
}

func TestProgressBarOverflow(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	bar := NewProgressBar(stream, "comp", "res", 100)
	bar.Add(150) // Exceeds total

	time.Sleep(50 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	// Should be capped at 100%
	if events[0].Progress != 100 {
		t.Errorf("Progress: got %d, want %d (should be capped)", events[0].Progress, 100)
	}
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		data     []byte
		b        byte
		expected int
	}{
		{[]byte("hello\nworld"), '\n', 5},
		{[]byte("no newline"), '\n', -1},
		{[]byte(""), '\n', -1},
		{[]byte("\n"), '\n', 0},
		{[]byte("abc"), 'b', 1},
	}

	for _, tt := range tests {
		result := indexOf(tt.data, tt.b)
		if result != tt.expected {
			t.Errorf("indexOf(%q, %q): got %d, want %d", tt.data, tt.b, result, tt.expected)
		}
	}
}

func TestMultiWriter(t *testing.T) {
	var buf1, buf2 bytes.Buffer

	mw := MultiWriter(&buf1, &buf2)
	_, _ = mw.Write([]byte("test data"))

	if buf1.String() != "test data" {
		t.Errorf("buf1: got %q, want %q", buf1.String(), "test data")
	}
	if buf2.String() != "test data" {
		t.Errorf("buf2: got %q, want %q", buf2.String(), "test data")
	}
}

func TestConcurrentEventEmission(t *testing.T) {
	stream := NewStream()
	handler := &mockHandler{}
	stream.AddHandler(handler)

	const numGoroutines = 10
	const eventsPerGoroutine = 100

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				stream.EmitInfo("comp", "res", "message")
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)
	stream.Close()

	events := handler.Events()
	expected := numGoroutines * eventsPerGoroutine
	if len(events) != expected {
		t.Errorf("Expected %d events, got %d", expected, len(events))
	}
}
