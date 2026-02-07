// Package output provides live output streaming for cldctl operations.
package output

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Level indicates the log level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Event represents a streaming event.
type Event struct {
	// Time is when the event occurred
	Time time.Time

	// Level is the log level
	Level Level

	// Component is the component name
	Component string

	// Resource is the resource name
	Resource string

	// Action is the current action
	Action string

	// Message is the event message
	Message string

	// Progress is the completion percentage (0-100, -1 for unknown)
	Progress int

	// Metadata contains additional data
	Metadata map[string]interface{}
}

// Handler processes streaming events.
type Handler interface {
	// HandleEvent processes a single event
	HandleEvent(event Event)

	// Close closes the handler
	Close() error
}

// Stream manages output streaming.
type Stream struct {
	mu       sync.RWMutex
	handlers []Handler
	events   chan Event
	done     chan struct{}
	wg       sync.WaitGroup
}

// NewStream creates a new output stream.
func NewStream() *Stream {
	s := &Stream{
		handlers: []Handler{},
		events:   make(chan Event, 100),
		done:     make(chan struct{}),
	}

	// Start event processor
	s.wg.Add(1)
	go s.processEvents()

	return s
}

// AddHandler adds an event handler.
func (s *Stream) AddHandler(h Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers = append(s.handlers, h)
}

// Emit emits an event to all handlers.
func (s *Stream) Emit(event Event) {
	if event.Time.IsZero() {
		event.Time = time.Now()
	}

	select {
	case s.events <- event:
	case <-s.done:
	}
}

// EmitInfo emits an info-level event.
func (s *Stream) EmitInfo(component, resource, message string) {
	s.Emit(Event{
		Level:     LevelInfo,
		Component: component,
		Resource:  resource,
		Message:   message,
		Progress:  -1,
	})
}

// EmitProgress emits a progress event.
func (s *Stream) EmitProgress(component, resource, action string, progress int) {
	s.Emit(Event{
		Level:     LevelInfo,
		Component: component,
		Resource:  resource,
		Action:    action,
		Progress:  progress,
	})
}

// EmitError emits an error-level event.
func (s *Stream) EmitError(component, resource string, err error) {
	s.Emit(Event{
		Level:     LevelError,
		Component: component,
		Resource:  resource,
		Message:   err.Error(),
		Progress:  -1,
	})
}

// Close closes the stream.
func (s *Stream) Close() error {
	close(s.done)
	s.wg.Wait()

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, h := range s.handlers {
		h.Close()
	}

	return nil
}

func (s *Stream) processEvents() {
	defer s.wg.Done()

	for {
		select {
		case event := <-s.events:
			s.mu.RLock()
			for _, h := range s.handlers {
				h.HandleEvent(event)
			}
			s.mu.RUnlock()

		case <-s.done:
			// Drain remaining events
			for {
				select {
				case event := <-s.events:
					s.mu.RLock()
					for _, h := range s.handlers {
						h.HandleEvent(event)
					}
					s.mu.RUnlock()
				default:
					return
				}
			}
		}
	}
}

// Writer returns an io.Writer that converts writes to events.
func (s *Stream) Writer(component, resource string, level Level) io.Writer {
	return &streamWriter{
		stream:    s,
		component: component,
		resource:  resource,
		level:     level,
	}
}

type streamWriter struct {
	stream    *Stream
	component string
	resource  string
	level     Level
	buf       []byte
}

func (w *streamWriter) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)

	// Process complete lines
	for {
		idx := indexOf(w.buf, '\n')
		if idx == -1 {
			break
		}

		line := string(w.buf[:idx])
		w.buf = w.buf[idx+1:]

		if line != "" {
			w.stream.Emit(Event{
				Level:     w.level,
				Component: w.component,
				Resource:  w.resource,
				Message:   line,
				Progress:  -1,
			})
		}
	}

	return len(p), nil
}

func indexOf(data []byte, b byte) int {
	for i, c := range data {
		if c == b {
			return i
		}
	}
	return -1
}

// ConsoleHandler writes events to the console with formatting.
type ConsoleHandler struct {
	writer    io.Writer
	useColors bool
	verbose   bool
	mu        sync.Mutex
}

// ConsoleOptions configures the console handler.
type ConsoleOptions struct {
	// Writer is the output writer (default: os.Stdout)
	Writer io.Writer

	// UseColors enables colored output
	UseColors bool

	// Verbose enables verbose output
	Verbose bool
}

// NewConsoleHandler creates a new console handler.
func NewConsoleHandler(opts ConsoleOptions) *ConsoleHandler {
	writer := opts.Writer
	if writer == nil {
		writer = os.Stdout
	}

	return &ConsoleHandler{
		writer:    writer,
		useColors: opts.UseColors,
		verbose:   opts.Verbose,
	}
}

func (h *ConsoleHandler) HandleEvent(event Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Build output line
	var line strings.Builder

	// Time prefix (only in verbose mode)
	if h.verbose {
		line.WriteString(event.Time.Format("15:04:05"))
		line.WriteString(" ")
	}

	// Level prefix
	levelStr := ""
	switch event.Level {
	case LevelDebug:
		if !h.verbose {
			return // Skip debug in non-verbose mode
		}
		levelStr = h.colorize("DEBUG", "cyan")
	case LevelInfo:
		levelStr = h.colorize("INFO ", "green")
	case LevelWarn:
		levelStr = h.colorize("WARN ", "yellow")
	case LevelError:
		levelStr = h.colorize("ERROR", "red")
	}
	line.WriteString(levelStr)
	line.WriteString(" ")

	// Component/Resource prefix
	if event.Component != "" {
		line.WriteString("[")
		line.WriteString(h.colorize(event.Component, "blue"))
		if event.Resource != "" {
			line.WriteString("/")
			line.WriteString(event.Resource)
		}
		line.WriteString("] ")
	}

	// Progress indicator
	if event.Progress >= 0 {
		line.WriteString(fmt.Sprintf("[%3d%%] ", event.Progress))
	}

	// Action
	if event.Action != "" {
		line.WriteString(event.Action)
		line.WriteString(": ")
	}

	// Message
	line.WriteString(event.Message)
	line.WriteString("\n")

	fmt.Fprint(h.writer, line.String())
}

func (h *ConsoleHandler) Close() error {
	return nil
}

func (h *ConsoleHandler) colorize(text, color string) string {
	if !h.useColors {
		return text
	}

	colors := map[string]string{
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"reset":   "\033[0m",
	}

	colorCode, ok := colors[color]
	if !ok {
		return text
	}

	return colorCode + text + colors["reset"]
}

// JSONHandler writes events as JSON lines.
type JSONHandler struct {
	writer io.Writer
	mu     sync.Mutex
}

// NewJSONHandler creates a new JSON handler.
func NewJSONHandler(writer io.Writer) *JSONHandler {
	return &JSONHandler{
		writer: writer,
	}
}

func (h *JSONHandler) HandleEvent(event Event) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Simple JSON formatting
	fmt.Fprintf(h.writer, `{"time":"%s","level":"%s","component":"%s","resource":"%s","action":"%s","message":"%s","progress":%d}`+"\n",
		event.Time.Format(time.RFC3339),
		levelToString(event.Level),
		event.Component,
		event.Resource,
		event.Action,
		escapeJSON(event.Message),
		event.Progress,
	)
}

func (h *JSONHandler) Close() error {
	return nil
}

func levelToString(level Level) string {
	switch level {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "unknown"
	}
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// MultiWriter creates a writer that duplicates writes to multiple writers.
func MultiWriter(writers ...io.Writer) io.Writer {
	return io.MultiWriter(writers...)
}

// LineScanner scans an io.Reader and emits events for each line.
type LineScanner struct {
	stream    *Stream
	component string
	resource  string
	level     Level
}

// NewLineScanner creates a new line scanner.
func NewLineScanner(stream *Stream, component, resource string, level Level) *LineScanner {
	return &LineScanner{
		stream:    stream,
		component: component,
		resource:  resource,
		level:     level,
	}
}

// Scan reads from the reader and emits events.
func (s *LineScanner) Scan(ctx context.Context, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			s.stream.Emit(Event{
				Level:     s.level,
				Component: s.component,
				Resource:  s.resource,
				Message:   scanner.Text(),
				Progress:  -1,
			})
		}
	}

	return scanner.Err()
}

// ProgressBar represents a text-based progress bar.
type ProgressBar struct {
	stream    *Stream
	component string
	resource  string
	total     int
	current   int
	mu        sync.Mutex
}

// NewProgressBar creates a new progress bar.
func NewProgressBar(stream *Stream, component, resource string, total int) *ProgressBar {
	return &ProgressBar{
		stream:    stream,
		component: component,
		resource:  resource,
		total:     total,
		current:   0,
	}
}

// Increment increments the progress by 1.
func (p *ProgressBar) Increment() {
	p.Add(1)
}

// Add adds to the current progress.
func (p *ProgressBar) Add(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current += n
	if p.current > p.total {
		p.current = p.total
	}

	progress := 0
	if p.total > 0 {
		progress = (p.current * 100) / p.total
	}

	p.stream.EmitProgress(p.component, p.resource, "", progress)
}

// SetCurrent sets the current progress.
func (p *ProgressBar) SetCurrent(current int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = current
	if p.current > p.total {
		p.current = p.total
	}

	progress := 0
	if p.total > 0 {
		progress = (p.current * 100) / p.total
	}

	p.stream.EmitProgress(p.component, p.resource, "", progress)
}

// Complete marks the progress as complete.
func (p *ProgressBar) Complete() {
	p.SetCurrent(p.total)
}
