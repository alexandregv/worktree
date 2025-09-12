package cmd

import (
	"os"

	"github.com/alexandregv/worktree/core"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                "worktree",
	Short:              "Git worktree CLI utility",
	Long:               `CLI utility to easily navigate between Git worktrees, list them, clone a multiple-worktrees-enabled repo, etc.`,
	DisableFlagParsing: true, // Needed for `wt -` to work
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 && args[0] == "-" {
			switchCmd.Run(switchCmd, args)
		} else {
			core.OpenTUI()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
