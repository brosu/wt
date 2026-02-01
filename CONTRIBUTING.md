# Contributing to wt

## Development Workflow

### For External Contributors (Fork-based workflow)

If you don't have write access to the repository:

1. **Fork the repository** on GitHub

2. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR-USERNAME/wt.git
   cd wt
   ```

3. **Add the upstream remote (recommended):**
   ```bash
   git remote add upstream https://github.com/timvw/wt.git
   git fetch upstream
   ```

4. **Create a feature branch:**
   ```bash
   # If you added upstream:
   git checkout -b feat/your-feature-name upstream/main
   # Otherwise:
   # git checkout -b feat/your-feature-name main
   ```

5. **Make your changes and commit:**
   ```bash
   git add .
   git commit -m "feat: your feature description"
   ```

6. **Push to your fork:**
   ```bash
   git push -u origin feat/your-feature-name
   ```

7. **Create a Pull Request** from your fork to `timvw/wt:main`
   - Go to https://github.com/timvw/wt
   - Click "New Pull Request"
   - Select "compare across forks"
   - Choose your fork and branch

8. **Wait for CI to pass** and respond to any review feedback

### For Maintainers (Branch-based workflow)

If you have write access to the repository:

1. **Create a feature branch:**
   ```bash
   git checkout -b feat/your-feature-name
   ```

2. **Make your changes and commit:**
   ```bash
   git add .
   git commit -m "feat: your feature description"
   ```

3. **Push the branch:**
   ```bash
   git push -u origin feat/your-feature-name
   ```

4. **Create a Pull Request:**
   ```bash
   gh pr create --title "feat: your feature" --body "Description of changes"
   ```

5. **Wait for CI to pass** - Branch protection requires all checks to pass

6. **Merge the PR** when CI is green

### Branching Workflow (personal `origin/main`, plus an `upstream` mirror branch)

Keep `origin/main` as your personal stable branch, and require every change to land via a PR.
To track the canonical repo, maintain a separate mirror branch in your fork (for example `origin/upstream-main`) that
fast-forwards from `upstream/main`.

In a fork setup:
- `origin` is your fork
- `upstream` is the canonical repo (e.g. `timvw/wt`)

If you haven't configured `upstream` yet:

```bash
git remote add upstream https://github.com/timvw/wt.git
git fetch upstream
```

#### Keep `origin/upstream-main` in sync with `upstream/main`

Create the mirror branch once (on your fork):

```bash
git fetch upstream
git push origin upstream/main:upstream-main
```

Update the mirror branch later (fast-forward it again):

```bash
git fetch upstream
git push origin upstream/main:upstream-main
```

#### F1: personal-only feature (stays on your fork)

- Create a feature branch from `origin/main`.
- Open a PR to `origin/main`.

```bash
git checkout -b feat/f1 origin/main
git push -u origin feat/f1
# PR: origin/feat/f1 -> origin/main
```

#### F2: upstream-only feature

- Create a feature branch from `upstream/main` (or your fork mirror branch `origin/upstream-main`).
- Open a PR to `upstream/main` only.

```bash
git checkout -b feat/f2 origin/upstream-main   # or upstream/main
git push -u origin feat/f2
# PR: origin/feat/f2 -> upstream/main
```

#### F3: shared feature (both `origin/main` and upstream)

1) Create the branch from `upstream/main` (or `origin/upstream-main`) and open a PR to upstream.
2) If you want it in `origin/main`, bring it in via a PR-only staging branch (no direct cherry-picks into `origin/main`).

```bash
# Upstream PR
git checkout -b feat/f3 origin/upstream-main   # or upstream/main
git push -u origin feat/f3
# PR: origin/feat/f3 -> upstream/main

# Origin PR (after or before upstream merges)
git checkout -b chore/f3-to-origin origin/main
git cherry-pick <f3-commit-sha>   # or merge if you prefer
git push -u origin chore/f3-to-origin
# PR: origin/chore/f3-to-origin -> origin/main
```

### Branch Naming Convention

- `feat/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation changes
- `refactor/description` - Code refactoring
- `chore/description` - Maintenance tasks

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add interactive selection for checkout`
- `fix: filter out invalid branch names`
- `docs: update installation instructions`
- `refactor: simplify branch filtering logic`
- `chore: update dependencies`
- `security: update vulnerable dependency`

### Running Tests Locally

Before pushing:

```bash
# Run tests
go test ./...

# Run linter
golangci-lint run

# Build
go build -o bin/wt .
```

### Branch Protection

The `main` branch is protected and requires:
- ✅ All CI checks must pass (Test, Build, Lint, Cross Compile)
- ✅ Branch must be up to date with main
- ❌ No direct pushes to main

## CI/CD

### Continuous Integration

Every push triggers:
- Tests on Go 1.25
- Linting with golangci-lint
- Build verification
- Cross-compilation checks

### Release Process

1. All changes merged to `main` via PRs
2. When ready to release:
   ```bash
   git tag v0.1.x
   git push origin v0.1.x
   ```
3. Automated workflow:
   - Builds binaries for all platforms
   - Creates Homebrew bottles
   - Publishes GitHub release
   - Updates Homebrew formula automatically

## Getting Help

- Check existing issues: https://github.com/timvw/wt/issues
- Read the README: https://github.com/timvw/wt#readme
- Ask questions in discussions
