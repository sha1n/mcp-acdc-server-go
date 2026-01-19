package domain

import (
	"fmt"
)

// ServerMetadata represents the server section of mcp-metadata.yaml
type ServerMetadata struct {
	Name         string `yaml:"name"`
	Version      string `yaml:"version"`
	Instructions string `yaml:"instructions"`
}

// ToolMetadata represents a tool definition in mcp-metadata.yaml
type ToolMetadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// McpMetadata represents the root of mcp-metadata.yaml
type McpMetadata struct {
	Server ServerMetadata `yaml:"server"`
	Tools  []ToolMetadata `yaml:"tools"`
}

// ToolsMap returns tools as a map for easy lookup
func (m *McpMetadata) ToolsMap() (map[string]ToolMetadata, error) {
	tools := make(map[string]ToolMetadata)
	for _, t := range m.Tools {
		if _, exists := tools[t.Name]; exists {
			return nil, fmt.Errorf("duplicate tool name: %s", t.Name)
		}
		tools[t.Name] = t
	}
	return tools, nil
}

// Validate checks for required fields
func (m *McpMetadata) Validate() error {
	if m.Server.Name == "" {
		return fmt.Errorf("server name is required")
	}
	if m.Server.Version == "" {
		return fmt.Errorf("server version is required")
	}
	if m.Server.Instructions == "" {
		return fmt.Errorf("server instructions are required")
	}

	for i, t := range m.Tools {
		if t.Name == "" {
			return fmt.Errorf("tool at index %d missing name", i)
		}
		if t.Description == "" {
			return fmt.Errorf("tool at index %d missing description", i)
		}
	}

	if _, err := m.ToolsMap(); err != nil {
		return err
	}

	return nil
}
