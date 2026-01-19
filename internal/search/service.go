package search

import (
	"fmt"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/sha1n/mcp-acdc-server-go/internal/config"
)

// SearchResult represents a search result
type SearchResult struct {
	URI     string
	Name    string
	Snippet string
}

// Document represents a document to index
type Document struct {
	URI     string
	Name    string
	Content string
}

// Service search service using Bleve
type Service struct {
	settings config.SearchSettings
	index    bleve.Index
	indexDir string
}

// NewService creates a new search service
func NewService(settings config.SearchSettings) *Service {
	return &Service{
		settings: settings,
	}
}

// IndexDocuments indexes a list of documents
func (s *Service) IndexDocuments(documents []Document) error {
	// Close existing index if any
	if s.index != nil {
		s.index.Close()
		s.index = nil
	}
	if s.indexDir != "" {
		os.RemoveAll(s.indexDir)
	}

	// Create temp dir
	tempDir, err := os.MkdirTemp("", "acdc_search_")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	s.indexDir = tempDir

	// Define mapping
	indexMapping := buildMapping()

	// Create index
	index, err := bleve.New(s.indexDir, indexMapping)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	s.index = index

	// Batch index
	batch := index.NewBatch()
	for _, doc := range documents {
		// We map the document struct to the fields
		// We can't pass the struct directly if we want specific field control unless we tag it
		// But building a map or strict struct with tags is better.
		// Let's use a map for clarity with the manual mapping we built
		docMap := map[string]interface{}{
			"uri":     doc.URI,
			"name":    doc.Name,
			"content": doc.Content,
		}
		if err := batch.Index(doc.URI, docMap); err != nil {
			return fmt.Errorf("failed to add document to batch: %w", err)
		}
	}

	if err := index.Batch(batch); err != nil {
		return fmt.Errorf("failed to execute batch index: %w", err)
	}

	return nil
}

func buildMapping() mapping.IndexMapping {
	// URI field: Stored, Indexed
	uriMapping := bleve.NewTextFieldMapping()
	uriMapping.Store = true
	uriMapping.IncludeInAll = false

	// Name field: Stored, Indexed, Included in All (default)
	nameMapping := bleve.NewTextFieldMapping()
	nameMapping.Store = true
	nameMapping.IncludeInAll = true
	nameMapping.Analyzer = "standard"

	// Content field: Indexed, Not Stored, Included in All
	contentMapping := bleve.NewTextFieldMapping()
	contentMapping.Store = false
	contentMapping.IncludeInAll = true
	contentMapping.Analyzer = "standard"

	docMapping := bleve.NewDocumentMapping()
	docMapping.AddFieldMappingsAt("uri", uriMapping)
	docMapping.AddFieldMappingsAt("name", nameMapping)
	docMapping.AddFieldMappingsAt("content", contentMapping)

	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping = docMapping
	return mapping
}

// Search searches for resources
func (s *Service) Search(queryStr string, limit *int) ([]SearchResult, error) {
	if s.index == nil {
		return []SearchResult{}, nil
	}

	maxResults := s.settings.MaxResults
	if limit != nil {
		maxResults = *limit
	}

	// Match query (searches across all fields included in 'all')
	// Python uses parse_query with default fields ["name", "content"]
	// Bleve's QueryStringQuery is similar
	query := bleve.NewQueryStringQuery(queryStr)
	
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = maxResults
	searchRequest.Fields = []string{"uri", "name"} // Fields to retrieve
	
	searchResult, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	results := make([]SearchResult, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		uri, ok1 := hit.Fields["uri"].(string)
		name, ok2 := hit.Fields["name"].(string)

		if !ok1 || !ok2 {
			continue
		}

		// Replicate Python snippet behavior
		snippet := fmt.Sprintf("%s (relevance: %.2f)", name, hit.Score)

		results = append(results, SearchResult{
			URI:     uri,
			Name:    name,
			Snippet: snippet,
		})
	}

	return results, nil
}

// Close cleans up resources
func (s *Service) Close() {
	if s.index != nil {
		s.index.Close()
	}
	if s.indexDir != "" {
		os.RemoveAll(s.indexDir)
	}
}
