# Git Repositories Indexing Specification

This specification describes the requirements and design for the Git repository indexing and search feature in the ACDC MCP Server.

## 1. Lifecycle and Resource Sharing (The "Coordinated Peer" Model)

Since this server is primarily used via `stdio` transport, multiple instances may run concurrently (one per agent/terminal). To ensure performance and low memory consumption:

### 1.1 Process Lifecycle
-   **Start**: Each instance is started independently by its parent MCP client (e.g., an AI agent or IDE).
-   **Stop**: Each instance is terminated by its parent client. There is no background "master" process to manage.
-   **Startup Optimization**: On startup, an instance checks for a `sync.lock` file in the cache directory.
    -   **If lock acquired**: This instance is the "Sync Leader". It performs `git pull` and updates the search index.
    -   **If lock exists**: This instance is a "Searcher". It skips the update phase and immediately opens the existing index in Read-Only mode.

### 1.2 Memory Management
-   **Shared Memory Mapping (mmap)**: All instances **must** use disk-based Bleve indexes (BoltDB). This allows the Operating System to share the same physical memory pages for index data across all running server processes.
-   **Heap Control**:
    -   `git-repos-max-file-size`: Files over 512KB are not indexed to prevent large buffers in the Go heap.
    -   **Lazy Indexing**: Git indexes are only opened when the `search_code` tool is invoked for the first time.

## 2. Configuration

| Environment Variable | CLI Flag | Description | Default |
| :--- | :--- | :--- | :--- |
| `ACDC_MCP_GIT_REPOS_CACHE_DIR` | `--git-repos-cache-dir` | Directory where git repositories are cloned. | `~/.acdc-mcp/git-repositories` |
| `ACDC_MCP_GIT_REPOS_URLS` | `--git-repos-urls` | Comma-separated list of Git repository URLs to index. | - |
| `ACDC_MCP_GIT_REPOS_AUTO_UPDATE` | `--git-repos-auto-update` | Periodically pull updates from repositories. | `true` |
| `ACDC_MCP_GIT_REPOS_UPDATE_INTERVAL` | `--git-repos-update-interval` | Interval between repository updates. | `30m` |
| `ACDC_MCP_GIT_REPOS_INDEX_DIR` | `--git-repos-index-dir` | Directory where indexes are stored. | `~/.acdc-mcp/indexes` |

## 3. Repository Management

### 3.1 Storage Layout
```text
~/.acdc-mcp/
├── git-repositories/
│   ├── github.com_user_repo1/  # Shallow clones
│   └── sync.lock               # Coordination lock
└── indexes/
    └── github.com_user_repo1/  # Bleve index files
```

### 3.2 Sync Workflow
1.  **Try Lock**: Attempt to create/acquire `sync.lock`.
2.  **Clone/Pull**: If locked, iterate through URLs. Use `git clone --depth 1` for new repos.
3.  **Index**: Update Bleve indexes.
4.  **Release**: Delete `sync.lock`.

## 4. MCP Search Tool: `search_code`

### 4.1 Tool Definition
-   **Name**: `search_code`
-   **Description**: "Searches across multiple git repositories for code patterns, documentation, and logic. Supports filtering by repository or file extension."

### 4.2 Parameters
```json
{
  "query": "string (Required) - Search query",
  "repo": "string (Optional) - Filter by repository name (substring match)",
  "extension": "string (Optional) - Filter by file extension (e.g., 'go', 'py')"
}
```

## 5. Implementation Roadmap

1.  **Config**: Add `GitSettings` to `internal/config/settings.go`.
2.  **Concurrency**: Implement a simple file-locker utility.
3.  **Git**: Use `go-git` to implement the Sync Workflow.
4.  **Search**: Extend the search service to support multiple Bleve indexes via `IndexAlias`.
5.  **Tool**: Register `search_code` in `internal/mcp/tools.go`.
