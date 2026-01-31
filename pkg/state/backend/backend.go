// Package backend defines the interface for state storage backends.
package backend

import (
	"context"
	"errors"
	"io"
	"time"
)

// ErrNotFound is returned when a requested state file doesn't exist.
var ErrNotFound = errors.New("state not found")

// ErrLocked is returned when state is already locked.
var ErrLocked = errors.New("state is locked")

// Backend defines the interface for state storage backends.
type Backend interface {
	// Type returns the backend type identifier (e.g., "s3", "local", "gcs")
	Type() string

	// Read reads state data from the given path.
	// Returns ErrNotFound if the path doesn't exist.
	Read(ctx context.Context, path string) (io.ReadCloser, error)

	// Write writes state data to the given path.
	// Creates parent directories/prefixes as needed.
	Write(ctx context.Context, path string, data io.Reader) error

	// Delete removes state data at the given path.
	// Returns nil if path doesn't exist (idempotent).
	Delete(ctx context.Context, path string) error

	// List lists state files under the given prefix.
	// Returns relative paths from the prefix.
	List(ctx context.Context, prefix string) ([]string, error)

	// Exists checks if a state file exists.
	Exists(ctx context.Context, path string) (bool, error)

	// Lock acquires a lock for the given path.
	// Blocks until lock is acquired or context is cancelled.
	Lock(ctx context.Context, path string, info LockInfo) (Lock, error)
}

// Lock represents an acquired lock.
type Lock interface {
	// ID returns the lock identifier.
	ID() string

	// Unlock releases the lock.
	Unlock(ctx context.Context) error

	// Info returns lock metadata.
	Info() LockInfo
}

// LockInfo contains metadata about a lock.
type LockInfo struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Who       string    `json:"who"`       // User or CI job identity
	Operation string    `json:"operation"` // What operation holds the lock
	Created   time.Time `json:"created"`
	Expires   time.Time `json:"expires,omitempty"` // Optional expiration
}

// LockError is returned when locking fails because state is already locked.
type LockError struct {
	Info LockInfo
	Err  error
}

func (e *LockError) Error() string {
	return e.Err.Error()
}

func (e *LockError) Unwrap() error {
	return e.Err
}

// Config holds configuration for creating a backend.
type Config struct {
	Type   string            `json:"type"`   // Backend type
	Config map[string]string `json:"config"` // Backend-specific configuration
}

// Factory creates a backend from configuration.
type Factory func(config map[string]string) (Backend, error)
