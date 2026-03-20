# Development

The project includes a `justfile` for common build tasks. Install [just](https://github.com/casey/just) to use it.

Available tasks:

```bash
just              # Show available recipes
just build        # Build the binary
just test         # Run unit tests
just e2e          # Run E2E tests
just test-all     # Run unit + E2E tests
just clean        # Clean build artifacts
just build-all    # Cross-compile for multiple platforms
just dev-shellenv # Print shell integration for running from source
```

## Running from source with shell integration

When hacking on `wt`, you can use `dev-shellenv` to get a shell function that
runs directly from source instead of an installed binary. This means every `wt`
command instantly reflects your code changes — no rebuild or reinstall needed.

```bash
eval "$(just dev-shellenv)"
```

This replaces the `wt` shell function so it calls `go run` against your local
checkout. Auto-cd and tab completion work as normal.

Re-run the `eval` line after changing the shell completion code in `shellenv`.

## Requirements

- Git
- `gh` CLI (optional, only needed for `wt pr` command to checkout GitHub PRs)
- `glab` CLI (optional, only needed for `wt mr` command to checkout GitLab MRs)

### For Building from Source

- Go 1.24+ (we support and test the latest two Go releases: 1.24 and 1.25)
- `just` (optional, for using the justfile)
