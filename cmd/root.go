package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alexandregv/worktree/pkg/core"
	"github.com/alexandregv/worktree/pkg/git"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Git worktree CLI utility",
	Long:  `CLI utility to easily navigate between Git worktrees, list them, clone a multiple-worktrees-enabled repo, etc.`,
	Args:  cobra.MatchAll(cobra.RangeArgs(0, 1), cobra.OnlyValidArgs),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) (validArgs []string, directive cobra.ShellCompDirective) {
		worktrees, err := git.GitWorktreeList()
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Could not get Git worktrees: %s\n", err.Error())
			os.Exit(1)
		}

		for _, wt := range worktrees {
			split := strings.Split(wt.Path, "/")
			validArgs = append(validArgs, split[len(split)-1])
		}
		return validArgs, directive
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			core.OpenTUI()
			return
		}

		worktrees, err := git.GitWorktreeList()
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Could not get Git worktrees: %s\n", err.Error())
			os.Exit(1)
		}

		for _, wt := range worktrees {
			if strings.HasSuffix(wt.Path, "/"+args[0]) {
				fmt.Println(wt.Path)
				os.Exit(0)
			}
		}
		fmt.Fprintf(os.Stderr, "worktree: Could not find Git worktree: %s\n", args[0])
		os.Exit(1)
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
