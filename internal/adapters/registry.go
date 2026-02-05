package adapters

import (
	"fmt"
	"sync"
)

// Registry manages adapter registration and selection.
// It provides both explicit adapter lookup and automatic detection based on
// directory structure inspection.
type Registry struct {
	mu       sync.RWMutex
	adapters map[string]Adapter
	priority []string // Ordered list for auto-detection priority
}

// NewRegistry creates a new adapter registry with default adapters pre-registered.
// Default adapters are registered in priority order:
//  1. acdc-mcp (native adapter)
//  2. legacy (backward compatibility)
func NewRegistry() *Registry {
	r := &Registry{
		adapters: make(map[string]Adapter),
		priority: []string{},
	}
	return r
}

// Register adds an adapter to the registry.
// If an adapter with the same name already exists, it will be replaced.
// The adapter is added to the end of the detection priority list.
func (r *Registry) Register(adapter Adapter) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := adapter.Name()
	r.adapters[name] = adapter

	// Add to priority list if not already present
	found := false
	for _, p := range r.priority {
		if p == name {
			found = true
			break
		}
	}
	if !found {
		r.priority = append(r.priority, name)
	}
}

// Get retrieves an adapter by name.
// Returns the adapter and true if found, nil and false otherwise.
func (r *Registry) Get(name string) (Adapter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adapter, ok := r.adapters[name]
	return adapter, ok
}

// AutoDetect selects the most appropriate adapter for the given base path.
// It checks adapters in priority order and returns the first one that can handle the path.
// Returns an error if no adapter can handle the path.
func (r *Registry) AutoDetect(basePath string) (Adapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, name := range r.priority {
		adapter, ok := r.adapters[name]
		if !ok {
			continue
		}

		if adapter.CanHandle(basePath) {
			return adapter, nil
		}
	}

	return nil, fmt.Errorf("no adapter found that can handle path: %s", basePath)
}

// List returns all registered adapter names in priority order.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]string, len(r.priority))
	copy(result, r.priority)
	return result
}

// Count returns the number of registered adapters.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.adapters)
}
