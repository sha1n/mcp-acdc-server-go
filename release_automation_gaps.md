# Release Automation Gap Analysis

## Required Adjustments

### `.goreleaser.yaml` Configuration
The current `.goreleaser.yaml` references an older naming convention and repository structure. The following changes are required to match the current codebase:

*   **Binary Name & Path:**
    *   **Change:** `binary: mcp-acdc` &rarr; `binary: acdc-mcp`
    *   **Change:** `main: ./cmd/mcp-acdc` &rarr; `main: ./cmd/acdc-mcp`
    *   **Reason:** The main application directory is `cmd/acdc-mcp` and the project convention (Makefile/Dockerfile) uses `acdc-mcp`.

*   **Build Flags (ldflags):**
    *   **Change:** `-X main.commit={{.Commit}}` &rarr; `-X main.Build={{.ShortCommit}}`
    *   **Add:** `-X main.ProgramName=acdc-mcp`
    *   **Remove:** `-X main.date={{.Date}}` (Not used in `main.go`)
    *   **Reason:** `cmd/acdc-mcp/main.go` uses `Build` variable for the commit hash, not `commit`.

*   **Homebrew Tap:**
    *   **Change:** `bin.install "mcp-acdc"` &rarr; `bin.install "acdc-mcp"`
    *   **Change:** `homepage: "https://github.com/sha1n/mcp-acdc-server-go"` &rarr; `homepage: "https://github.com/sha1n/mcp-acdc-server"`
    *   **Reason:** Matches new binary name and actual repository URL found in `go.mod`.

### Docker (Potential Gap)
*   **Observation:** The repository contains a `Dockerfile`, but `.goreleaser.yaml` does not currently include a `dockers` section to build and push images automatically.
*   **Recommendation:** If Docker image publication is desired as part of the release, a `dockers` section should be added to `.goreleaser.yaml`.

## Required Repo Changes

### Secrets Configuration
For the release workflow to succeed, the following GitHub Secrets must be set in the repository:
*   `HOMEBREW_TAP_GITHUB_TOKEN`: A Personal Access Token (PAT) with `repo` scope (or sufficient permissions) to push to the `sha1n/homebrew-tap` repository.

### Repository Consistency
*   **Repo Name:** Ensure the repository name on GitHub matches `mcp-acdc-server` as used in `go.mod` to avoid module path issues during `go install` or `go get`.
