package content

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContentProvider_LoadText(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(filePath, []byte("hello world"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	p := NewContentProvider(tempDir)
	content, err := p.LoadText(filePath)
	if err != nil {
		t.Fatalf("LoadText failed: %v", err)
	}
	if content != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", content)
	}
}

func TestContentProvider_LoadYAML(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.yaml")
	err := os.WriteFile(filePath, []byte("key: value"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	p := NewContentProvider(tempDir)
	data, err := p.LoadYAML(filePath)
	if err != nil {
		t.Fatalf("LoadYAML failed: %v", err)
	}
	if data["key"] != "value" {
		t.Errorf("Expected value 'value', got '%v'", data["key"])
	}
}

func TestContentProvider_LoadMarkdownWithFrontmatter(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.md")
	content := "---\nname: Test\n---\nMarkdown content"
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	p := NewContentProvider(tempDir)
	md, err := p.LoadMarkdownWithFrontmatter(filePath)
	if err != nil {
		t.Fatalf("LoadMarkdownWithFrontmatter failed: %v", err)
	}

	if md.Metadata["name"] != "Test" {
		t.Errorf("Expected metadata name 'Test', got '%v'", md.Metadata["name"])
	}
	if md.Content != "Markdown content" {
		t.Errorf("Expected content 'Markdown content', got '%s'", md.Content)
	}
}
