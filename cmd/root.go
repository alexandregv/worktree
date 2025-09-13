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
	Args:               cobra.RangeArgs(0, 1),
	ValidArgs:          []string{"-", "-h", "--help"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			switch args[0] {
			case "-":
				switchCmd.Run(switchCmd, args)
			default:
				cmd.Help()
			}
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
