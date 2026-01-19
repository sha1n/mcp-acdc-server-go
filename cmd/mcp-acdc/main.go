package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sha1n/mcp-acdc-server-go/internal/config"
	"github.com/sha1n/mcp-acdc-server-go/internal/content"
	"github.com/sha1n/mcp-acdc-server-go/internal/domain"
	"github.com/sha1n/mcp-acdc-server-go/internal/mcp"
	"github.com/sha1n/mcp-acdc-server-go/internal/resources"
	"github.com/sha1n/mcp-acdc-server-go/internal/search"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	// Version is injected at build time
	Version = "dev"
	// Build is injected at build time
	Build = "unknown"
	// ProgramName is injected at build time
	ProgramName = "mcp-acdc"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   ProgramName,
		Short: "MCP ACDC Server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	// Load settings
	settings, err := config.LoadSettings()
	if err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}

	// Configure logging
	var handler slog.Handler
	if settings.Transport == "stdio" {
		// Log to stderr for stdio transport
		handler = slog.NewTextHandler(os.Stderr, nil)
	} else {
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(handler))

	slog.Info("Starting MCP Acdc server", "version", Version, "transport", settings.Transport)

	// Initialize content provider
	cp := content.NewContentProvider(settings.ContentDir)

	// Load metadata
	metadataPath := cp.GetPath("mcp-metadata.yaml")

	mdBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata file: %w", err)
	}

	var metadata domain.McpMetadata
	if err := yaml.Unmarshal(mdBytes, &metadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	if err := metadata.Validate(); err != nil {
		return fmt.Errorf("metadata validation failed: %w", err)
	}

	// Discover resources
	resourceDefinitions, err := resources.DiscoverResources(cp)
	if err != nil {
		return fmt.Errorf("failed to discover resources: %w", err)
	}

	resourceProvider := resources.NewResourceProvider(resourceDefinitions)

	// Initialize search service
	searchService := search.NewService(settings.Search)
	defer searchService.Close()

	// Index resources
	docsToIndex, err := resourceProvider.GetAllResourceContents()
	if err != nil {
		slog.Error("Failed to get resource contents for indexing", "error", err)
	} else {
		var docs []search.Document
		for _, d := range docsToIndex {
			docs = append(docs, search.Document{
				URI:     d["uri"],
				Name:    d["name"],
				Content: d["content"],
			})
		}

		if err := searchService.IndexDocuments(docs); err != nil {
			slog.Error("Failed to index documents", "error", err)
		} else {
			slog.Info("Indexed documents", "count", len(docs))
		}
	}

	// Create MCP server
	mcpServer := mcp.CreateServer(metadata, resourceProvider, searchService)

	// Start server
	if settings.Transport == "stdio" {
		return server.ServeStdio(mcpServer)
	} else {
		slog.Info("Starting SSE server", "host", settings.Host, "port", settings.Port)
		// Mark3labs SSE server implementation
		sseServer := server.NewSSEServer(mcpServer)
		return sseServer.Start(fmt.Sprintf("%s:%d", settings.Host, settings.Port))
	}
}
