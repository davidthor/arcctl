package expression

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name         string
		input        string
		wantSegments int
		wantLiteral  bool
		wantRefs     int
	}{
		{
			name:         "literal string",
			input:        "hello world",
			wantSegments: 1,
			wantLiteral:  true,
			wantRefs:     0,
		},
		{
			name:         "single expression",
			input:        "${{ databases.main.url }}",
			wantSegments: 1,
			wantLiteral:  false,
			wantRefs:     1,
		},
		{
			name:         "expression with surrounding text",
			input:        "postgresql://${{ databases.main.host }}:5432",
			wantSegments: 3,
			wantLiteral:  false,
			wantRefs:     1,
		},
		{
			name:         "multiple expressions",
			input:        "${{ variables.host }}:${{ variables.port }}",
			wantSegments: 3,
			wantLiteral:  false,
			wantRefs:     2,
		},
		{
			name:         "expression with pipe",
			input:        "${{ dependents.*.routes.*.url | join ',' }}",
			wantSegments: 1,
			wantLiteral:  false,
			wantRefs:     1,
		},
		{
			name:         "empty string",
			input:        "",
			wantSegments: 1,
			wantLiteral:  true,
			wantRefs:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(expr.Segments) != tt.wantSegments {
				t.Errorf("expected %d segments, got %d", tt.wantSegments, len(expr.Segments))
			}

			if expr.IsLiteral() != tt.wantLiteral {
				t.Errorf("expected IsLiteral=%v, got %v", tt.wantLiteral, expr.IsLiteral())
			}

			refs := expr.References()
			if len(refs) != tt.wantRefs {
				t.Errorf("expected %d references, got %d", tt.wantRefs, len(refs))
			}
		})
	}
}

func TestParser_ParseReference(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		input    string
		wantPath []string
		wantPipe int
	}{
		{
			name:     "simple path",
			input:    "${{ databases.main.url }}",
			wantPath: []string{"databases", "main", "url"},
			wantPipe: 0,
		},
		{
			name:     "path with wildcard",
			input:    "${{ dependents.*.routes.*.url }}",
			wantPath: []string{"dependents", "*", "routes", "*", "url"},
			wantPipe: 0,
		},
		{
			name:     "path with pipe function",
			input:    "${{ dependents.*.routes.*.url | join }}",
			wantPath: []string{"dependents", "*", "routes", "*", "url"},
			wantPipe: 1,
		},
		{
			name:     "path with multiple pipes",
			input:    "${{ values | first | upper }}",
			wantPath: []string{"values"},
			wantPipe: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(expr.Segments) != 1 {
				t.Fatalf("expected 1 segment, got %d", len(expr.Segments))
			}

			ref, ok := expr.Segments[0].(ReferenceSegment)
			if !ok {
				t.Fatal("expected ReferenceSegment")
			}

			if len(ref.Path) != len(tt.wantPath) {
				t.Errorf("expected path length %d, got %d", len(tt.wantPath), len(ref.Path))
			}

			for i, p := range tt.wantPath {
				if i < len(ref.Path) && ref.Path[i] != p {
					t.Errorf("path[%d]: expected %q, got %q", i, p, ref.Path[i])
				}
			}

			if len(ref.Pipe) != tt.wantPipe {
				t.Errorf("expected %d pipe functions, got %d", tt.wantPipe, len(ref.Pipe))
			}
		})
	}
}

func TestContainsExpression(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"hello", false},
		{"${{ foo }}", true},
		{"prefix ${{ bar }} suffix", true},
		{"${ not an expression }", false},
		{"{{ also not }}", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ContainsExpression(tt.input); got != tt.want {
				t.Errorf("ContainsExpression(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
