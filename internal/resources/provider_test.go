package resources

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sha1n/mcp-acdc-server-go/internal/content"
)

func TestDiscoverResources(t *testing.T) {
	tempDir := t.TempDir()
	resourcesDir := filepath.Join(tempDir, "mcp-resources")
	os.MkdirAll(resourcesDir, 0755)

	// Create a test resource
	resourcePath := filepath.Join(resourcesDir, "test.md")
	contentStr := "---\nname: Test Resource\ndescription: A test resource\n---\n# Test Content"
	os.WriteFile(resourcePath, []byte(contentStr), 0644)

	cp := content.NewContentProvider(tempDir)
	defs, err := DiscoverResources(cp)
	if err != nil {
		t.Fatalf("DiscoverResources failed: %v", err)
	}

	if len(defs) != 1 {
		t.Fatalf("Expected 1 resource, got %d", len(defs))
	}

	def := defs[0]
	if def.Name != "Test Resource" {
		t.Errorf("Expected name 'Test Resource', got '%s'", def.Name)
	}
	if def.URI != "acdc://test" {
		t.Errorf("Expected URI 'acdc://test', got '%s'", def.URI)
	}
}

func TestResourceProvider_ReadResource(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	resourcesDir := filepath.Join(tempDir, "mcp-resources")
	os.MkdirAll(resourcesDir, 0755)
	resourcePath := filepath.Join(resourcesDir, "test.md")
	contentStr := "---\nname: Test\ndescription: Desc\n---\nContent"
	os.WriteFile(resourcePath, []byte(contentStr), 0644)

	defs := []ResourceDefinition{
		{
			URI:         "acdc://test",
			Name:        "Test",
			Description: "Desc",
			MIMEType:    "text/markdown",
			FilePath:    resourcePath,
		},
	}

	p := NewResourceProvider(defs)
	
	// Test Read
	content, err := p.ReadResource("acdc://test")
	if err != nil {
		t.Fatalf("ReadResource failed: %v", err)
	}
	if content != "Content" {
		t.Errorf("Expected 'Content', got '%s'", content)
	}
}
