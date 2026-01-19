package search

import (
	"testing"

	"github.com/sha1n/mcp-acdc-server-go/internal/config"
)

func TestSearchService(t *testing.T) {
	settings := config.SearchSettings{
		MaxResults: 10,
		HeapSizeMB: 10,
	}
	service := NewService(settings)
	defer service.Close()

	docs := []Document{
		{
			URI:     "acdc://doc1",
			Name:    "Document One",
			Content: "This is the content of document one. It talks about testing.",
		},
		{
			URI:     "acdc://doc2",
			Name:    "Document Two",
			Content: "This is the content of document two. It talks about development.",
		},
	}

	if err := service.IndexDocuments(docs); err != nil {
		t.Fatalf("IndexDocuments failed: %v", err)
	}

	// Search for "testing"
	results, err := service.Search("testing", nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
	if results[0].URI != "acdc://doc1" {
		t.Errorf("Expected URI 'acdc://doc1', got '%s'", results[0].URI)
	}
	
	// Search for "document"
	results, err = service.Search("document", nil)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
}
