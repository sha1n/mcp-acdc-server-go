package adapters

import (
	"testing"

	"github.com/sha1n/mcp-acdc-server/internal/content"
	"github.com/sha1n/mcp-acdc-server/internal/prompts"
	"github.com/sha1n/mcp-acdc-server/internal/resources"
)

// mockAdapter is a test implementation of the Adapter interface.
type mockAdapter struct {
	name               string
	canHandle          bool
	resourcesResult    []resources.ResourceDefinition
	resourcesError     error
	promptsResult      []prompts.PromptDefinition
	promptsError       error
	canHandleCalled    bool
	discoverResCalled  bool
	discoverPromCalled bool
}

func (m *mockAdapter) Name() string {
	return m.name
}

func (m *mockAdapter) CanHandle(basePath string) bool {
	m.canHandleCalled = true
	return m.canHandle
}

func (m *mockAdapter) DiscoverResources(location Location, cp *content.ContentProvider) ([]resources.ResourceDefinition, error) {
	m.discoverResCalled = true
	return m.resourcesResult, m.resourcesError
}

func (m *mockAdapter) DiscoverPrompts(location Location, cp *content.ContentProvider) ([]prompts.PromptDefinition, error) {
	m.discoverPromCalled = true
	return m.promptsResult, m.promptsError
}

// TestAdapterInterface verifies that the adapter interface contract can be implemented.
func TestAdapterInterface(t *testing.T) {
	t.Run("mock adapter implements interface", func(t *testing.T) {
		var _ Adapter = &mockAdapter{}
	})

	t.Run("adapter Name returns identifier", func(t *testing.T) {
		adapter := &mockAdapter{name: "test-adapter"}
		if got := adapter.Name(); got != "test-adapter" {
			t.Errorf("Name() = %q, want %q", got, "test-adapter")
		}
	})

	t.Run("adapter CanHandle returns boolean", func(t *testing.T) {
		tests := []struct {
			name      string
			canHandle bool
		}{
			{"returns true", true},
			{"returns false", false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				adapter := &mockAdapter{canHandle: tt.canHandle}
				if got := adapter.CanHandle("/test/path"); got != tt.canHandle {
					t.Errorf("CanHandle() = %v, want %v", got, tt.canHandle)
				}
				if !adapter.canHandleCalled {
					t.Error("CanHandle was not called")
				}
			})
		}
	})

	t.Run("adapter DiscoverResources returns results", func(t *testing.T) {
		expectedRes := []resources.ResourceDefinition{
			{Name: "test-resource"},
		}
		adapter := &mockAdapter{resourcesResult: expectedRes}

		location := Location{Name: "test", BasePath: "/test"}
		got, err := adapter.DiscoverResources(location, nil)

		if err != nil {
			t.Errorf("DiscoverResources() unexpected error: %v", err)
		}
		if len(got) != len(expectedRes) {
			t.Errorf("DiscoverResources() returned %d resources, want %d", len(got), len(expectedRes))
		}
		if !adapter.discoverResCalled {
			t.Error("DiscoverResources was not called")
		}
	})

	t.Run("adapter DiscoverPrompts returns results", func(t *testing.T) {
		expectedPrompts := []prompts.PromptDefinition{
			{Name: "test-prompt"},
		}
		adapter := &mockAdapter{promptsResult: expectedPrompts}

		location := Location{Name: "test", BasePath: "/test"}
		got, err := adapter.DiscoverPrompts(location, nil)

		if err != nil {
			t.Errorf("DiscoverPrompts() unexpected error: %v", err)
		}
		if len(got) != len(expectedPrompts) {
			t.Errorf("DiscoverPrompts() returned %d prompts, want %d", len(got), len(expectedPrompts))
		}
		if !adapter.discoverPromCalled {
			t.Error("DiscoverPrompts was not called")
		}
	})
}

// TestLocation verifies the Location type structure.
func TestLocation(t *testing.T) {
	t.Run("location with all fields", func(t *testing.T) {
		loc := Location{
			Name:        "docs",
			BasePath:    "/path/to/docs",
			AdapterType: "acdc-mcp",
		}

		if loc.Name != "docs" {
			t.Errorf("Name = %q, want %q", loc.Name, "docs")
		}
		if loc.BasePath != "/path/to/docs" {
			t.Errorf("BasePath = %q, want %q", loc.BasePath, "/path/to/docs")
		}
		if loc.AdapterType != "acdc-mcp" {
			t.Errorf("AdapterType = %q, want %q", loc.AdapterType, "acdc-mcp")
		}
	})

	t.Run("location with empty adapter type", func(t *testing.T) {
		loc := Location{
			Name:        "guides",
			BasePath:    "/path/to/guides",
			AdapterType: "",
		}

		if loc.AdapterType != "" {
			t.Errorf("AdapterType = %q, want empty string", loc.AdapterType)
		}
	})
}
