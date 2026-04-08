package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove worktree administrative files",
	RunE: func(cmd *cobra.Command, args []string) error {
		gitCmd := exec.Command("git", "worktree", "prune")
		if !isJSONOutput() {
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
		}
		if err := gitCmd.Run(); err != nil {
			return err
		}

		if isJSONOutput() {
			return emitJSONSuccess(cmd, map[string]any{"status": "pruned"})
		}

		fmt.Println("✓ Pruned stale worktree administrative files")
		return nil
	},
}
