package errors

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	err := New(ErrCodeValidation, "test message")

	if err.Code != ErrCodeValidation {
		t.Errorf("expected code %s, got %s", ErrCodeValidation, err.Code)
	}

	if err.Message != "test message" {
		t.Errorf("expected message 'test message', got '%s'", err.Message)
	}

	if err.Cause != nil {
		t.Error("expected nil cause")
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := Wrap(ErrCodeParse, "wrapped error", cause)

	if err.Code != ErrCodeParse {
		t.Errorf("expected code %s, got %s", ErrCodeParse, err.Code)
	}

	if err.Cause != cause {
		t.Error("expected cause to be set")
	}

	if !errors.Is(err, cause) {
		t.Error("Unwrap should return the cause")
	}
}

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name:     "without cause",
			err:      New(ErrCodeValidation, "validation failed"),
			expected: "[VALIDATION_ERROR] validation failed",
		},
		{
			name:     "with cause",
			err:      Wrap(ErrCodeParse, "parse failed", errors.New("syntax error")),
			expected: "[PARSE_ERROR] parse failed: syntax error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestWithDetails(t *testing.T) {
	err := New(ErrCodeValidation, "test").
		WithDetail("field", "name").
		WithDetail("value", 123)

	if err.Details["field"] != "name" {
		t.Error("expected field detail to be set")
	}

	if err.Details["value"] != 123 {
		t.Error("expected value detail to be set")
	}
}

func TestValidationError(t *testing.T) {
	err := ValidationError("invalid input", map[string]interface{}{
		"field": "email",
	})

	if err.Code != ErrCodeValidation {
		t.Errorf("expected code %s, got %s", ErrCodeValidation, err.Code)
	}

	if err.Details["field"] != "email" {
		t.Error("expected field detail to be set")
	}
}

func TestNotFoundError(t *testing.T) {
	err := NotFoundError("component", "my-app")

	if err.Code != ErrCodeNotFound {
		t.Errorf("expected code %s, got %s", ErrCodeNotFound, err.Code)
	}

	if err.Details["resource_type"] != "component" {
		t.Error("expected resource_type detail")
	}

	if err.Details["name"] != "my-app" {
		t.Error("expected name detail")
	}
}

func TestParseError(t *testing.T) {
	err := ParseError("config.yaml", errors.New("invalid yaml"))

	if err.Code != ErrCodeParse {
		t.Errorf("expected code %s, got %s", ErrCodeParse, err.Code)
	}

	if err.Details["file"] != "config.yaml" {
		t.Error("expected file detail")
	}
}
