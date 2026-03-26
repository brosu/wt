package main

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all worktrees",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isJSONOutput() {
			entries, err := getWorktreeListPorcelain()
			if err != nil {
				return err
			}
			return emitJSONSuccess(cmd, map[string]any{"worktrees": entries})
		}

		gitCmd := exec.Command("git", "worktree", "list")
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr
		if err := gitCmd.Run(); err != nil {
			return err
		}
		return nil
	},
}
