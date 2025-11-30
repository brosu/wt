# Docker-based E2E Tests

This directory contains Docker-based end-to-end tests for `wt` on Linux, testing shell integration across bash, zsh, and fish.

## Overview

The E2E suite verifies that `wt` works correctly in different shell environments by:

1. Building the `wt` binary from source
2. Creating a temporary git repository with test branches
3. Loading the shell environment via `wt shellenv`
4. Running core operations: `checkout`, `create`, `list`, and `remove`
5. Verifying auto-cd functionality and worktree management

## Running Locally

### Prerequisites

- Docker installed and running
- From the repository root

### Build the test image

```bash
docker build -t wt-e2e -f e2e/docker/Dockerfile .
```

### Run tests for a specific shell

```bash
# Bash
docker run --rm -v $(pwd):/workspace -w /workspace wt-e2e bash e2e/docker/test-bash.sh

# Zsh
docker run --rm -v $(pwd):/workspace -w /workspace wt-e2e zsh e2e/docker/test-zsh.sh

# Fish
docker run --rm -v $(pwd):/workspace -w /workspace wt-e2e fish e2e/docker/test-fish.sh
```

### Run all tests

```bash
./e2e/docker/run-all.sh
```

## Test Structure

Each shell test script (`test-bash.sh`, `test-zsh.sh`, `test-fish.sh`) follows the same pattern:

1. **Setup**: Build `wt` binary and create temporary test repository
2. **Test 1**: Verify `wt checkout` with existing branch and auto-cd
3. **Test 2**: Verify `wt create` with new branch and auto-cd
4. **Test 3**: Verify `wt list` displays worktrees
5. **Test 4**: Verify `wt remove` deletes worktree
6. **Cleanup**: Remove temporary directories

## CI Integration

The tests run automatically in GitHub Actions on every push and pull request via the `linux-e2e.yml` workflow. The workflow:

- Uses a matrix strategy to run tests in parallel for bash, zsh, and fish
- Caches the Docker image using GitHub Actions cache
- Uploads test logs as artifacts on failure
- Keeps runtime under 5-6 minutes per shell

## Troubleshooting

If tests fail locally:

1. Ensure Docker is running
2. Check that you're running from the repository root
3. Verify the Docker image builds successfully
4. Run with verbose output: add `set -x` to the test script

If tests fail in CI:

1. Check the uploaded test logs in the Actions artifacts
2. Look for shell-specific issues in the shellenv output
3. Verify git configuration is set correctly
