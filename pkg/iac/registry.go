package iac

import (
	"fmt"
	"sync"
)

// Factory creates a plugin instance.
type Factory func() (Plugin, error)

// Registry manages plugin factories.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// DefaultRegistry is the global plugin registry.
var DefaultRegistry = &Registry{
	factories: make(map[string]Factory),
}

// Register adds a plugin factory to the registry.
func (r *Registry) Register(name string, factory Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[name] = factory
}

// Get retrieves a plugin by name.
func (r *Registry) Get(name string) (Plugin, error) {
	r.mu.RLock()
	factory, ok := r.factories[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown plugin: %s", name)
	}

	return factory()
}

// List returns all registered plugin names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}

// Register registers a plugin factory with the default registry.
func Register(name string, factory Factory) {
	DefaultRegistry.Register(name, factory)
}

// Get retrieves a plugin from the default registry.
func Get(name string) (Plugin, error) {
	return DefaultRegistry.Get(name)
}
