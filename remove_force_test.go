package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRemoveForceFlagRemovesDirtyWorktree(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping force remove test in short mode")
	}

	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "test-repo")
	worktreeRoot := filepath.Join(tmpDir, "worktrees")

	setupTestRepo(t, repoDir)
	wtBinary := buildWtBinary(t, tmpDir)

	// Prepare a branch and checkout a worktree for it
	runGitCommand(t, repoDir, "branch", "force-remove-branch")

	checkoutCmd := exec.Command(wtBinary, "checkout", "force-remove-branch")
	checkoutCmd.Dir = repoDir
	checkoutCmd.Env = append(os.Environ(), "WORKTREE_ROOT="+worktreeRoot)
	if output, err := checkoutCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to create worktree: %v\nOutput: %s", err, output)
	}

	worktreePath := filepath.Join(worktreeRoot, "test-repo", "force-remove-branch")

	// Make the worktree dirty so a normal remove should fail
	if err := os.WriteFile(filepath.Join(worktreePath, "dirty.txt"), []byte("dirty"), 0o644); err != nil {
		t.Fatalf("Failed to create dirty file in worktree: %v", err)
	}

	removeCmd := exec.Command(wtBinary, "remove", "force-remove-branch")
	removeCmd.Dir = repoDir
	removeCmd.Env = append(os.Environ(), "WORKTREE_ROOT="+worktreeRoot)
	if output, err := removeCmd.CombinedOutput(); err == nil {
		t.Fatal("Expected remove without --force to fail for dirty worktree")
	} else if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Fatalf("Worktree unexpectedly removed without --force: %v\nOutput: %s", err, output)
	}

	forceCmd := exec.Command(wtBinary, "remove", "--force", "force-remove-branch")
	forceCmd.Dir = repoDir
	forceCmd.Env = append(os.Environ(), "WORKTREE_ROOT="+worktreeRoot)
	if output, err := forceCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to remove worktree with --force: %v\nOutput: %s", err, output)
	}

	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		t.Fatalf("Expected worktree path to be removed with --force, got err: %v", err)
	}
}
