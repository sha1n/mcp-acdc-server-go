package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sha1n/mcp-acdc-server-go/internal/domain"
	"github.com/sha1n/mcp-acdc-server-go/internal/resources"
	"github.com/sha1n/mcp-acdc-server-go/internal/search"
)

// CreateServer creates and configures the MCP server
func CreateServer(
	metadata domain.McpMetadata,
	resourceProvider *resources.ResourceProvider,
	searchService *search.Service,
) *server.MCPServer {
	// Create server
	s := server.NewMCPServer(
		metadata.Server.Name,
		metadata.Server.Version,
		server.WithInstructions(metadata.Server.Instructions),
	)

	// Register Resources
	for _, res := range resourceProvider.ListResources() {
		// Capture uri for closure
		uri := res.URI 
		
		s.AddResource(mcp.Resource{
			URI:         uri,
			Name:        res.Name,
			Description: res.Description,
			MIMEType:    res.MIMEType,
		}, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			content, err := resourceProvider.ReadResource(uri)
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      uri,
					MIMEType: "text/markdown",
					Text:     content,
				},
			}, nil
		})
	}
	
	// Register Tools
	toolsMap, _ := metadata.ToolsMap()
	if toolMeta, ok := toolsMap["search"]; ok {
		RegisterSearchTool(s, searchService, toolMeta)
	}

	return s
}
