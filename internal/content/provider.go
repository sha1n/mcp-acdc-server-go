package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// MarkdownWithFrontmatter parsed markdown file with YAML frontmatter
type MarkdownWithFrontmatter struct {
	Metadata map[string]interface{}
	Content  string
}

// ContentProvider provider for loading content files
type ContentProvider struct {
	ContentDir   string
	ResourcesDir string
}

// NewContentProvider creates a new ContentProvider
func NewContentProvider(contentDir string) *ContentProvider {
	return &ContentProvider{
		ContentDir:   contentDir,
		ResourcesDir: filepath.Join(contentDir, "mcp-resources"),
	}
}

// GetPath returns a path within the content directory
func (p *ContentProvider) GetPath(parts ...string) string {
	allParts := append([]string{p.ContentDir}, parts...)
	return filepath.Join(allParts...)
}

// LoadText loads a text file
func (p *ContentProvider) LoadText(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// LoadYAML loads and parses a YAML file
func (p *ContentProvider) LoadYAML(filePath string) (map[string]interface{}, error) {
	content, err := p.LoadText(filePath)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("invalid YAML in %s: %w", filePath, err)
	}

	return data, nil
}

// LoadMarkdownWithFrontmatter loads a markdown file with YAML frontmatter
func (p *ContentProvider) LoadMarkdownWithFrontmatter(filePath string) (*MarkdownWithFrontmatter, error) {
	content, err := p.LoadText(filePath)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(content, "---") {
		return nil, fmt.Errorf("file must start with YAML frontmatter (---) in %s", filePath)
	}

	// Find end of frontmatter
	// Python implementation looks for "\n---\n" starting from index 4
	endIndex := strings.Index(content[4:], "\n---\n")
	if endIndex == -1 {
		return nil, fmt.Errorf("invalid frontmatter format - missing closing --- in %s", filePath)
	}

	// Adjust index to be relative to start of string
	realEndIndex := endIndex + 4 

	frontmatterText := content[4:realEndIndex]
	markdownContent := content[realEndIndex+5:] // Skip "\n---\n" (length 5)

	var metadata map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontmatterText), &metadata); err != nil {
		return nil, fmt.Errorf("invalid YAML in frontmatter of %s: %w", filePath, err)
	}

	return &MarkdownWithFrontmatter{
		Metadata: metadata,
		Content:  markdownContent,
	}, nil
}
