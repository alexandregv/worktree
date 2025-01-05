package cmd

import (
	"os"

	"github.com/alexandregv/worktree/pkg/core"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Git worktree CLI utility",
	Long:  `CLI utility to easily navigate between Git worktrees, list them, clone a multiple-worktrees-enabled repo, etc.`,
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 1)),
	Run: func(cmd *cobra.Command, args []string) {
		core.OpenTUI()
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
