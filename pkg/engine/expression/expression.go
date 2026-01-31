// Package expression provides expression parsing and evaluation for arcctl.
package expression

// No imports needed for this file

// Expression represents a parsed expression.
type Expression struct {
	Raw      string    // Original expression text
	Segments []Segment // Parsed segments
}

// Segment is part of an expression.
type Segment interface {
	segment()
}

// LiteralSegment is a literal string.
type LiteralSegment struct {
	Value string
}

func (LiteralSegment) segment() {}

// ReferenceSegment is a reference to a value.
type ReferenceSegment struct {
	Path []string   // e.g., ["databases", "main", "url"]
	Pipe []PipeFunc // Optional pipe functions
}

func (ReferenceSegment) segment() {}

// PipeFunc represents a pipe function call.
type PipeFunc struct {
	Name string
	Args []string
}

// IsLiteral returns true if the expression has no dynamic references.
func (e *Expression) IsLiteral() bool {
	if len(e.Segments) == 1 {
		_, ok := e.Segments[0].(LiteralSegment)
		return ok
	}
	return len(e.Segments) == 0
}

// HasReferences returns true if the expression contains references.
func (e *Expression) HasReferences() bool {
	for _, seg := range e.Segments {
		if _, ok := seg.(ReferenceSegment); ok {
			return true
		}
	}
	return false
}

// References returns all reference paths in the expression.
func (e *Expression) References() [][]string {
	var refs [][]string
	for _, seg := range e.Segments {
		if ref, ok := seg.(ReferenceSegment); ok {
			refs = append(refs, ref.Path)
		}
	}
	return refs
}
