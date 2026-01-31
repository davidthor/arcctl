package backend

import (
	"fmt"
	"sync"
)

// Registry manages backend factories.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// DefaultRegistry is the global backend registry.
var DefaultRegistry = &Registry{
	factories: make(map[string]Factory),
}

// Register adds a backend factory to the registry.
func (r *Registry) Register(backendType string, factory Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[backendType] = factory
}

// Create instantiates a backend from configuration.
func (r *Registry) Create(config Config) (Backend, error) {
	r.mu.RLock()
	factory, ok := r.factories[config.Type]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("unknown backend type: %s", config.Type)
	}

	return factory(config.Config)
}

// List returns all registered backend types.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.factories))
	for t := range r.factories {
		types = append(types, t)
	}
	return types
}

// Register registers a backend factory with the default registry.
func Register(backendType string, factory Factory) {
	DefaultRegistry.Register(backendType, factory)
}

// Create creates a backend using the default registry.
func Create(config Config) (Backend, error) {
	return DefaultRegistry.Create(config)
}
