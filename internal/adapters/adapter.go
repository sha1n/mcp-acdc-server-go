// Package adapters provides a flexible system for discovering resources and prompts
// from different content directory structures. The adapter pattern allows ACDC to
// support multiple content layouts (native, legacy, Claude Code plugins, etc.)
// through a unified interface.
package adapters

import (
	"github.com/sha1n/mcp-acdc-server/internal/content"
	"github.com/sha1n/mcp-acdc-server/internal/prompts"
	"github.com/sha1n/mcp-acdc-server/internal/resources"
)

// Adapter defines the interface for discovering resources and prompts from a content location.
// Implementations provide support for different directory structures and naming conventions.
//
// Example implementations:
//   - ACDC adapter: uses resources/ and prompts/ directories
//   - Legacy adapter: uses mcp-resources/ and mcp-prompts/ directories
//   - Claude Code adapter: uses skills/ and commands/ directories
type Adapter interface {
	// Name returns the unique identifier for this adapter (e.g., "acdc-mcp", "legacy").
	Name() string

	// CanHandle checks if this adapter can handle the content at the given base path.
	// It should inspect the directory structure to determine compatibility.
	// Returns true if the adapter recognizes the structure, false otherwise.
	CanHandle(basePath string) bool

	// DiscoverResources scans the content location for resources and returns their definitions.
	// The ContentProvider is passed to enable reading file contents.
	// Returns a slice of ResourceDefinition or an error if discovery fails.
	DiscoverResources(location Location, cp *content.ContentProvider) ([]resources.ResourceDefinition, error)

	// DiscoverPrompts scans the content location for prompts and returns their definitions.
	// The ContentProvider is passed to enable reading file contents.
	// Returns a slice of PromptDefinition or an error if discovery fails.
	DiscoverPrompts(location Location, cp *content.ContentProvider) ([]prompts.PromptDefinition, error)
}

// Location represents a content location with its resolved adapter information.
// This is passed to adapter methods to provide context about where content is located
// and how it should be identified.
type Location struct {
	// Name is the content location identifier (e.g., "docs", "guides").
	// Used as the source prefix in URIs: acdc://<name>/path
	Name string

	// BasePath is the absolute path to the content root directory.
	BasePath string

	// AdapterType is the explicit adapter type specified in configuration.
	// Empty string means auto-detection should be used.
	AdapterType string
}
