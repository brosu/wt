package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isJSONOutput() {
			return emitJSONSuccess(cmd, map[string]string{"version": version})
		}
		fmt.Printf("wt version %s\n", version)
		return nil
	},
}
