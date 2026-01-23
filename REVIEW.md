# Remediation Plan for Commit `a917bf76`

## Objective
Address critical feedback from the code review regarding argument validation, error handling resilience, and performance optimization for the MCP Prompts feature.

## Proposed Changes

### 1. Robustness & Error Handling
**Target:** `internal/prompts/provider.go` - `DiscoverPrompts`

*   **Fix `os.Stat` check:** Do not swallow all errors. Only return nil/empty if `os.IsNotExist(err)` is true. For other errors (e.g., permissions), log a warning and return the error.
*   **Resilient Walking:** Modify the `filepath.WalkDir` callback.
    *   If an error occurs for a specific file/dir, log it (`slog.Error`) and return `nil` to continue walking instead of halting the entire server startup.

### 2. Argument Validation
**Target:** `internal/prompts/provider.go` - `GetPrompt`

*   **Implement Validation:** Iterate through the prompt's defined arguments.
*   **Check Required:** If an argument is marked `required: true` (or defaults to true) and is missing from the input map (or empty), return a structured error immediately.
*   **Template Strictness:** Consider using `tmpl.Option("missingkey=error")` to catch any other undefined variables referenced in the template.

### 3. Performance Optimization
**Target:** `internal/prompts/provider.go`

*   **Cache Templates:** instead of parsing the markdown file on every `GetPrompt` call:
    *   Store the parsed `*template.Template` inside `PromptDefinition` (or a parallel map in the provider).
    *   Parse the template *once* during `DiscoverPrompts`.
    *   This removes file I/O and parsing overhead from the hot path.

### 4. Test Coverage
**Target:** `internal/prompts/provider_test.go`

*   **Validation Tests:** Add a test case asserting that `GetPrompt` fails when a required argument is missing.
*   **Resilience Tests:**
    *   Simulate a "partial failure" (e.g., one unreadable file among valid ones) to ensure discovery continues.
    *   Verify `os.Stat` behavior for non-existent vs. permission-denied directories (if feasible with testkit).

## Execution Steps

1.  Refactor `DiscoverPrompts` to pre-load templates and handle FS errors gracefully.
2.  Update `PromptDefinition` to hold the parsed template.
3.  Update `GetPrompt` to use the cached template and validate arguments.
4.  Add comprehensive unit tests covering these failure modes.
5.  Run `make test` to verify.
