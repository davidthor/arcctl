package expression

import (
	"testing"
)

func TestJoinFunc(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		args  []string
		want  string
	}{
		{
			name:  "join strings with comma",
			value: []string{"a", "b", "c"},
			args:  []string{","},
			want:  "a,b,c",
		},
		{
			name:  "join strings with space",
			value: []string{"hello", "world"},
			args:  []string{" "},
			want:  "hello world",
		},
		{
			name:  "join interfaces",
			value: []interface{}{"a", "b", "c"},
			args:  []string{"-"},
			want:  "a-b-c",
		},
		{
			name:  "default separator",
			value: []string{"a", "b"},
			args:  []string{},
			want:  "a,b",
		},
		{
			name:  "single value",
			value: "single",
			args:  []string{","},
			want:  "single",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := joinFunc(tt.value, tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestFirstFunc(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  interface{}
	}{
		{
			name:  "first of strings",
			value: []string{"a", "b", "c"},
			want:  "a",
		},
		{
			name:  "first of interfaces",
			value: []interface{}{1, 2, 3},
			want:  1,
		},
		{
			name:  "empty slice",
			value: []string{},
			want:  "",
		},
		{
			name:  "single value",
			value: "single",
			want:  "single",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := firstFunc(tt.value, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestLastFunc(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  interface{}
	}{
		{
			name:  "last of strings",
			value: []string{"a", "b", "c"},
			want:  "c",
		},
		{
			name:  "last of interfaces",
			value: []interface{}{1, 2, 3},
			want:  3,
		},
		{
			name:  "empty slice",
			value: []string{},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lastFunc(tt.value, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestLengthFunc(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  int
	}{
		{
			name:  "length of slice",
			value: []string{"a", "b", "c"},
			want:  3,
		},
		{
			name:  "length of string",
			value: "hello",
			want:  5,
		},
		{
			name:  "length of empty",
			value: []string{},
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := lengthFunc(tt.value, nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestDefaultFunc(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		args  []string
		want  interface{}
	}{
		{
			name:  "use value when present",
			value: "actual",
			args:  []string{"default"},
			want:  "actual",
		},
		{
			name:  "use default when nil",
			value: nil,
			args:  []string{"default"},
			want:  "default",
		},
		{
			name:  "use default when empty string",
			value: "",
			args:  []string{"default"},
			want:  "default",
		},
		{
			name:  "use default when empty slice",
			value: []interface{}{},
			args:  []string{"default"},
			want:  "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defaultFunc(tt.value, tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestUpperFunc(t *testing.T) {
	got, _ := upperFunc("hello", nil)
	if got != "HELLO" {
		t.Errorf("expected 'HELLO', got %q", got)
	}
}

func TestLowerFunc(t *testing.T) {
	got, _ := lowerFunc("HELLO", nil)
	if got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestTrimFunc(t *testing.T) {
	got, _ := trimFunc("  hello  ", nil)
	if got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}
