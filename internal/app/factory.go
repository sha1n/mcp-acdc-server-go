package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sha1n/mcp-acdc-server/internal/config"
	"github.com/sha1n/mcp-acdc-server/internal/content"
	"github.com/sha1n/mcp-acdc-server/internal/domain"
	"github.com/sha1n/mcp-acdc-server/internal/mcp"
	"github.com/sha1n/mcp-acdc-server/internal/prompts"
	"github.com/sha1n/mcp-acdc-server/internal/resources"
	"github.com/sha1n/mcp-acdc-server/internal/search"
	"gopkg.in/yaml.v3"
)

// CreateMCPServer initializes the core MCP server components
func CreateMCPServer(settings *config.Settings) (*mcpsdk.Server, func(), error) {
	// Initialize content provider
	cp := content.NewContentProvider(settings.ContentDir)

	// Load metadata
	metadataPath := cp.GetPath("mcp-metadata.yaml")

	mdBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata domain.McpMetadata
	if err := yaml.Unmarshal(mdBytes, &metadata); err != nil {
		return nil, nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	if err := metadata.Validate(); err != nil {
		return nil, nil, fmt.Errorf("metadata validation failed: %w", err)
	}

	// Discover resources
	resourceDefinitions, err := resources.DiscoverResources(cp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to discover resources: %w", err)
	}

	resourceProvider := resources.NewResourceProvider(resourceDefinitions)

	// Discover prompts
	promptDefinitions, err := prompts.DiscoverPrompts(cp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to discover prompts: %w", err)
	}

	promptProvider := prompts.NewPromptProvider(promptDefinitions, cp)

	// Initialize search service
	searchService := search.NewService(settings.Search)
	cleanup := func() {
		searchService.Close()
	}

	// Index resources
	docsChan := make(chan domain.Document, 100)
	ctx := context.Background()
	go func() {
		defer close(docsChan)
		if err := resourceProvider.StreamResources(ctx, docsChan); err != nil {
			slog.Error("StreamResources failed", "error", err)
		}
	}()

	if err := searchService.Index(ctx, docsChan); err != nil {
		slog.Error("Failed to index documents", "error", err)
	} else {
		slog.Info("Indexed documents finished")
	}

	// Create MCP server
	mcpServer := mcp.CreateServer(metadata, resourceProvider, promptProvider, searchService)

	return mcpServer, cleanup, nil
}
