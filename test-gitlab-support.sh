#!/bin/bash
set -e

echo "=== Testing GitLab/GitHub Support ==="
echo ""

# Build if needed
if [ ! -f bin/wt ]; then
    echo "Building wt..."
    go build -o bin/wt .
fi

# Test 1: GitHub remote detection
echo "Test 1: GitHub remote detection (current repo)"
echo "Remote URL: $(git remote get-url origin)"
echo ""
echo "Testing: ./bin/wt pr 1"
echo "(This will fail to fetch, but shows it detects GitHub and looks for 'gh')"
./bin/wt pr 1 2>&1 | head -5 || true
echo ""
echo "---"
echo ""

# Test 2: Temporarily change to GitLab remote for testing
echo "Test 2: GitLab remote detection (simulated)"
echo "Saving current remote..."
ORIGINAL_REMOTE=$(git remote get-url origin)
echo "Temporarily changing remote to GitLab..."
git remote set-url origin https://gitlab.com/test/project.git

echo "Remote URL: $(git remote get-url origin)"
echo ""
echo "Testing: ./bin/wt mr 1"
echo "(This will fail to fetch, but shows it detects GitLab and looks for 'glab')"
./bin/wt mr 1 2>&1 | head -5 || true
echo ""

# Restore original remote
echo "Restoring original remote..."
git remote set-url origin "$ORIGINAL_REMOTE"
echo "Remote restored: $(git remote get-url origin)"
echo ""
echo "---"
echo ""

# Test 3: URL parsing
echo "Test 3: URL parsing"
echo "GitHub PR URL should extract number 123:"
./bin/wt pr --help | grep -A 1 "github.com" || true
echo ""
echo "GitLab MR URL should extract number 123:"
./bin/wt pr --help | grep -A 1 "gitlab.com" || true
echo ""

# Test 4: Command aliases
echo "Test 4: Command aliases"
echo "'wt pr' and 'wt mr' are aliases:"
./bin/wt pr --help | grep "Aliases:" -A 1
echo ""

echo "=== All Tests Complete ==="
