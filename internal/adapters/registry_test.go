package adapters

import (
	"strings"
	"sync"
	"testing"

	"github.com/sha1n/mcp-acdc-server/internal/content"
	"github.com/sha1n/mcp-acdc-server/internal/prompts"
	"github.com/sha1n/mcp-acdc-server/internal/resources"
)

// TestNewRegistry verifies registry initialization.
func TestNewRegistry(t *testing.T) {
	r := NewRegistry()

	if count := r.Count(); count != 0 {
		t.Errorf("new registry Count() = %d, want 0", count)
	}

	if len(r.List()) != 0 {
		t.Errorf("new registry List() length = %d, want 0", len(r.List()))
	}
}

// TestRegister verifies adapter registration.
func TestRegister(t *testing.T) {
	t.Run("register single adapter", func(t *testing.T) {
		r := NewRegistry()
		adapter := &mockAdapter{name: "test1"}

		r.Register(adapter)

		if count := r.Count(); count != 1 {
			t.Errorf("Count() = %d, want 1", count)
		}

		got, ok := r.Get("test1")
		if !ok {
			t.Error("Get(test1) returned false, want true")
		}
		if got != adapter {
			t.Error("Get(test1) returned different adapter instance")
		}
	})

	t.Run("register multiple adapters", func(t *testing.T) {
		r := NewRegistry()
		adapter1 := &mockAdapter{name: "test1"}
		adapter2 := &mockAdapter{name: "test2"}
		adapter3 := &mockAdapter{name: "test3"}

		r.Register(adapter1)
		r.Register(adapter2)
		r.Register(adapter3)

		if count := r.Count(); count != 3 {
			t.Errorf("Count() = %d, want 3", count)
		}

		list := r.List()
		if len(list) != 3 {
			t.Errorf("List() length = %d, want 3", len(list))
		}
	})

	t.Run("register replaces existing adapter", func(t *testing.T) {
		r := NewRegistry()
		adapter1 := &mockAdapter{name: "test", canHandle: true}
		adapter2 := &mockAdapter{name: "test", canHandle: false}

		r.Register(adapter1)
		r.Register(adapter2)

		if count := r.Count(); count != 1 {
			t.Errorf("Count() = %d, want 1", count)
		}

		got, _ := r.Get("test")
		if got.CanHandle("/test") {
			t.Error("expected second adapter to replace first")
		}
	})

	t.Run("register maintains priority order", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&mockAdapter{name: "first"})
		r.Register(&mockAdapter{name: "second"})
		r.Register(&mockAdapter{name: "third"})

		list := r.List()
		expected := []string{"first", "second", "third"}

		for i, name := range expected {
			if list[i] != name {
				t.Errorf("priority[%d] = %q, want %q", i, list[i], name)
			}
		}
	})

	t.Run("register same adapter twice preserves priority", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&mockAdapter{name: "first"})
		r.Register(&mockAdapter{name: "second"})
		r.Register(&mockAdapter{name: "first"}) // re-register

		list := r.List()
		if len(list) != 2 {
			t.Errorf("List() length = %d, want 2", len(list))
		}

		expected := []string{"first", "second"}
		for i, name := range expected {
			if list[i] != name {
				t.Errorf("priority[%d] = %q, want %q", i, list[i], name)
			}
		}
	})
}

// TestGet verifies adapter retrieval.
func TestGet(t *testing.T) {
	t.Run("get existing adapter", func(t *testing.T) {
		r := NewRegistry()
		adapter := &mockAdapter{name: "test"}
		r.Register(adapter)

		got, ok := r.Get("test")
		if !ok {
			t.Error("Get(test) returned false, want true")
		}
		if got != adapter {
			t.Error("Get(test) returned different adapter")
		}
	})

	t.Run("get non-existent adapter", func(t *testing.T) {
		r := NewRegistry()

		got, ok := r.Get("nonexistent")
		if ok {
			t.Error("Get(nonexistent) returned true, want false")
		}
		if got != nil {
			t.Error("Get(nonexistent) returned non-nil adapter")
		}
	})

	t.Run("get from empty registry", func(t *testing.T) {
		r := NewRegistry()

		_, ok := r.Get("test")
		if ok {
			t.Error("Get from empty registry returned true")
		}
	})
}

// TestAutoDetect verifies automatic adapter selection.
func TestAutoDetect(t *testing.T) {
	t.Run("detect first matching adapter", func(t *testing.T) {
		r := NewRegistry()
		adapter1 := &mockAdapter{name: "first", canHandle: false}
		adapter2 := &mockAdapter{name: "second", canHandle: true}
		adapter3 := &mockAdapter{name: "third", canHandle: true}

		r.Register(adapter1)
		r.Register(adapter2)
		r.Register(adapter3)

		got, err := r.AutoDetect("/test/path")
		if err != nil {
			t.Errorf("AutoDetect() unexpected error: %v", err)
		}
		if got.Name() != "second" {
			t.Errorf("AutoDetect() returned %q, want %q", got.Name(), "second")
		}
	})

	t.Run("detect with priority order", func(t *testing.T) {
		r := NewRegistry()
		// Register in specific order
		r.Register(&mockAdapter{name: "low-priority", canHandle: true})
		r.Register(&mockAdapter{name: "high-priority", canHandle: true})

		got, err := r.AutoDetect("/test/path")
		if err != nil {
			t.Errorf("AutoDetect() unexpected error: %v", err)
		}
		// Should return first in priority order
		if got.Name() != "low-priority" {
			t.Errorf("AutoDetect() returned %q, want %q", got.Name(), "low-priority")
		}
	})

	t.Run("no adapter can handle path", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&mockAdapter{name: "test1", canHandle: false})
		r.Register(&mockAdapter{name: "test2", canHandle: false})

		got, err := r.AutoDetect("/test/path")
		if err == nil {
			t.Error("AutoDetect() expected error, got nil")
		}
		if got != nil {
			t.Error("AutoDetect() expected nil adapter on error")
		}
		if !strings.Contains(err.Error(), "no adapter found") {
			t.Errorf("AutoDetect() error = %q, want 'no adapter found'", err.Error())
		}
		if !strings.Contains(err.Error(), "/test/path") {
			t.Errorf("AutoDetect() error should include path, got: %q", err.Error())
		}
	})

	t.Run("empty registry returns error", func(t *testing.T) {
		r := NewRegistry()

		_, err := r.AutoDetect("/test/path")
		if err == nil {
			t.Error("AutoDetect() on empty registry expected error")
		}
	})
}

// TestList verifies adapter listing.
func TestList(t *testing.T) {
	t.Run("list empty registry", func(t *testing.T) {
		r := NewRegistry()

		list := r.List()
		if len(list) != 0 {
			t.Errorf("List() length = %d, want 0", len(list))
		}
	})

	t.Run("list returns copy", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&mockAdapter{name: "test"})

		list1 := r.List()
		list1[0] = "modified"

		list2 := r.List()
		if list2[0] == "modified" {
			t.Error("List() should return a copy, not reference")
		}
	})

	t.Run("list returns priority order", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&mockAdapter{name: "a"})
		r.Register(&mockAdapter{name: "b"})
		r.Register(&mockAdapter{name: "c"})

		list := r.List()
		expected := []string{"a", "b", "c"}

		for i, name := range expected {
			if list[i] != name {
				t.Errorf("List()[%d] = %q, want %q", i, list[i], name)
			}
		}
	})
}

// TestCount verifies adapter counting.
func TestCount(t *testing.T) {
	t.Run("count empty registry", func(t *testing.T) {
		r := NewRegistry()

		if count := r.Count(); count != 0 {
			t.Errorf("Count() = %d, want 0", count)
		}
	})

	t.Run("count after registrations", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&mockAdapter{name: "test1"})
		r.Register(&mockAdapter{name: "test2"})

		if count := r.Count(); count != 2 {
			t.Errorf("Count() = %d, want 2", count)
		}
	})

	t.Run("count after replacement", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&mockAdapter{name: "test"})
		r.Register(&mockAdapter{name: "test"}) // replace

		if count := r.Count(); count != 1 {
			t.Errorf("Count() = %d, want 1", count)
		}
	})
}

// TestConcurrency verifies thread-safe operations.
func TestConcurrency(t *testing.T) {
	t.Run("concurrent register and get", func(t *testing.T) {
		r := NewRegistry()
		var wg sync.WaitGroup

		// Concurrent registrations
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				adapter := &mockAdapter{name: strings.Repeat("a", n+1)}
				r.Register(adapter)
			}(i)
		}

		// Concurrent gets
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(n int) {
				defer wg.Done()
				name := strings.Repeat("a", n+1)
				r.Get(name)
			}(i)
		}

		wg.Wait()

		if count := r.Count(); count != 10 {
			t.Errorf("after concurrent operations Count() = %d, want 10", count)
		}
	})

	t.Run("concurrent auto-detect", func(t *testing.T) {
		r := NewRegistry()
		r.Register(&mockAdapter{name: "test", canHandle: true})

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = r.AutoDetect("/test/path")
			}()
		}

		wg.Wait()
	})
}

// TestRegistryWithRealAdapterScenario tests a realistic usage scenario.
func TestRegistryWithRealAdapterScenario(t *testing.T) {
	r := NewRegistry()

	// Simulate registering adapters in priority order
	acdcAdapter := &mockAdapter{
		name:      "acdc-mcp",
		canHandle: false, // doesn't handle this path
	}
	legacyAdapter := &mockAdapter{
		name:            "legacy",
		canHandle:       true, // handles this path
		resourcesResult: []resources.ResourceDefinition{{Name: "test-resource"}},
		promptsResult:   []prompts.PromptDefinition{{Name: "test-prompt"}},
	}

	r.Register(acdcAdapter)
	r.Register(legacyAdapter)

	// Auto-detect should skip acdc-mcp and select legacy
	detected, err := r.AutoDetect("/legacy/path")
	if err != nil {
		t.Errorf("AutoDetect() unexpected error: %v", err)
	}
	if detected.Name() != "legacy" {
		t.Errorf("AutoDetect() = %q, want %q", detected.Name(), "legacy")
	}

	// Use the detected adapter
	location := Location{Name: "test", BasePath: "/legacy/path"}
	resources, err := detected.DiscoverResources(location, &content.ContentProvider{})
	if err != nil {
		t.Errorf("DiscoverResources() unexpected error: %v", err)
	}
	if len(resources) != 1 {
		t.Errorf("DiscoverResources() returned %d resources, want 1", len(resources))
	}
}
