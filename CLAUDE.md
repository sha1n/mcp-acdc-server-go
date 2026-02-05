# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ACDC (Agent Content Discovery Companion) is an MCP server that provides AI agents with full-text search and resource discovery capabilities over local Markdown content. It uses the official MCP Go SDK and Bleve for search indexing.

## Build & Development Commands

```bash
# Build for current platform
make go-build-current

# Run all tests
make test

# Run a single test
go test -run TestFunctionName ./path/to/package

# Run tests with coverage
make coverage

# Lint and format
make lint
make format

# Build Docker image
make build-docker

# Add dev server to Claude Code (uses examples/sample-content)
make mcp-add-claude-dev
```

## Architecture

### Package Structure

- `cmd/acdc-mcp/` - CLI entrypoint using Cobra
- `internal/app/` - Application wiring: CLI flags, server factory, SSE/stdio runners
- `internal/mcp/` - MCP server creation and tool registration (search, read)
- `internal/content/` - Content directory abstraction
- `internal/resources/` - Resource discovery and provider (scans `mcp-resources/`)
- `internal/prompts/` - Prompt discovery and provider (scans `mcp-prompts/`)
- `internal/search/` - Bleve-based search service
- `internal/config/` - Settings loading (viper-based, supports env vars and flags)
- `internal/auth/` - Authentication middleware (basic, apikey)
- `internal/domain/` - Core types (metadata, search results)

### Key Flow

1. `cmd/acdc-mcp/main.go` → `app.RunWithDeps()` handles CLI and starts server
2. `app.CreateMCPServer()` initializes all components:
   - Loads `mcp-metadata.yaml` for server identity
   - Discovers resources and prompts from content directory
   - Indexes resources into Bleve search
   - Creates MCP server with tools and resources
3. Server runs in either `stdio` mode (default) or `sse` mode (HTTP)

### Content Structure

The server expects this structure in the content directory:
```
content/
├── mcp-metadata.yaml    # Server identity (required)
├── mcp-resources/       # Markdown resources with YAML frontmatter
└── mcp-prompts/         # Prompt templates (optional)
```

### Integration Tests

Integration tests in `tests/integration/` use a testkit that:
- Creates temporary content directories
- Spawns real MCP servers
- Uses the typed MCP SDK client for assertions

See `tests/integration/testkit/` for helpers like `CreateTestContentDir()` and `NewTestFlags()`.
