package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	cleanupDryRun bool
	cleanupForce  bool
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove worktrees for merged branches",
	Long: `Remove worktrees for branches that have been merged into the base branch.

This command finds all worktrees whose branches have been merged into main/master,
and removes them. Use --dry-run to preview what would be removed.

Examples:
  wt cleanup              # Interactive confirmation for each worktree
  wt cleanup --dry-run    # Preview what would be removed
  wt cleanup --force      # Remove all without confirmation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		base := getDefaultBase()
		jsonMode := isJSONOutput()

		// Get merged branches
		mergedBranches, err := getMergedBranches(base)
		if err != nil {
			return err
		}

		// Get existing worktree branches
		worktreeBranches, err := getExistingWorktreeBranches()
		if err != nil {
			return fmt.Errorf("failed to get worktrees: %w", err)
		}

		// Create a set of merged branches for quick lookup
		mergedSet := make(map[string]bool)
		for _, b := range mergedBranches {
			mergedSet[b] = true
		}

		// Find worktrees that are for merged branches
		var toRemove []string
		for _, branch := range worktreeBranches {
			if mergedSet[branch] {
				toRemove = append(toRemove, branch)
			}
		}

		if len(toRemove) == 0 {
			if jsonMode {
				return emitJSONSuccess(cmd, map[string]any{"removed": 0, "skipped": 0, "base": base, "worktrees": []string{}})
			}
			fmt.Println("No worktrees found for merged branches")
			return nil
		}

		if jsonMode && !cleanupDryRun && !cleanupForce {
			return fmt.Errorf("wt cleanup with --format json requires --force or --dry-run")
		}

		// Dry run mode - just show what would be removed
		if cleanupDryRun {
			if jsonMode {
				planned := make([]map[string]string, 0, len(toRemove))
				for _, branch := range toRemove {
					if path, exists := worktreeExists(branch); exists {
						planned = append(planned, map[string]string{"branch": branch, "path": path})
					}
				}
				return emitJSONSuccess(cmd, map[string]any{"dry_run": true, "base": base, "worktrees": planned})
			}
			fmt.Printf("Would remove %d worktree(s) for merged branches:\n", len(toRemove))
			for _, branch := range toRemove {
				if path, exists := worktreeExists(branch); exists {
					fmt.Printf("  - %s (%s)\n", branch, path)
				}
			}
			return nil
		}

		// Track results
		removed := 0
		skipped := 0

		for _, branch := range toRemove {
			existingPath, exists := worktreeExists(branch)
			if !exists {
				continue
			}

			// If not force mode, ask for confirmation
			if !cleanupForce {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("Remove worktree for merged branch '%s'", branch),
					IsConfirm: true,
				}
				_, err := prompt.Run()
				if err != nil {
					fmt.Printf("  Skipped: %s\n", branch)
					skipped++
					continue
				}
			}

			// Remove the worktree
			gitCmd := exec.Command("git", "worktree", "remove", existingPath)
			if !jsonMode {
				gitCmd.Stdout = os.Stdout
				gitCmd.Stderr = os.Stderr
			}
			if err := gitCmd.Run(); err != nil {
				if jsonMode {
					skipped++
					continue
				}
				fmt.Printf("  Failed to remove %s: %v\n", branch, err)
				continue
			}

			if err := cleanupWorktreePath(existingPath); err != nil {
				if jsonMode {
					continue
				}
				fmt.Printf("  Warning: failed to cleanup path for %s: %v\n", branch, err)
			}

			if !jsonMode {
				fmt.Printf("✓ Removed worktree: %s\n", branch)
			}
			removed++
		}

		// Run prune at the end
		pruneGitCmd := exec.Command("git", "worktree", "prune")
		_ = pruneGitCmd.Run()

		if jsonMode {
			return emitJSONSuccess(cmd, map[string]any{"dry_run": false, "base": base, "removed": removed, "skipped": skipped})
		}

		fmt.Printf("\nCleanup complete: %d removed, %d skipped\n", removed, skipped)
		return nil
	},
}
