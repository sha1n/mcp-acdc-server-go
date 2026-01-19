package config

import (
	"os"
	"testing"
)

func TestLoadSettings_Defaults(t *testing.T) {
	_ = os.Unsetenv("ACDC_MCP_HOST")
	_ = os.Unsetenv("ACDC_MCP_SEARCH_MAX_RESULTS")

	settings, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings failed: %v", err)
	}

	if settings.Host != "0.0.0.0" {
		t.Errorf("Expected default host 0.0.0.0, got %s", settings.Host)
	}
	if settings.Search.MaxResults != 10 {
		t.Errorf("Expected default max results 10, got %d", settings.Search.MaxResults)
	}
}

func TestLoadSettings_EnvVars(t *testing.T) {
	_ = os.Setenv("ACDC_MCP_HOST", "127.0.0.1")
	_ = os.Setenv("ACDC_MCP_SEARCH_MAX_RESULTS", "50")
	defer func() { _ = os.Unsetenv("ACDC_MCP_HOST") }()
	defer func() { _ = os.Unsetenv("ACDC_MCP_SEARCH_MAX_RESULTS") }()

	settings, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings failed: %v", err)
	}

	if settings.Host != "127.0.0.1" {
		t.Errorf("Expected host 127.0.0.1, got %s", settings.Host)
	}
	if settings.Search.MaxResults != 50 {
		t.Errorf("Expected max results 50, got %d", settings.Search.MaxResults)
	}
}
