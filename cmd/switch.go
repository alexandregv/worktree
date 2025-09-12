package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alexandregv/worktree/core"
	"github.com/alexandregv/worktree/git"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:     "switch",
	Short:   "Move to a worktree",
	Aliases: []string{"s", "cd"},
	Args:    cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) (validArgs []string, directive cobra.ShellCompDirective) {
		validArgs = append(validArgs, "-")

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
		if args[0] == "-" {
			lastWorktree, err := git.GetConfig("alexandregv-worktree.lastworktree")
			if err != nil {
				fmt.Fprintf(os.Stderr, "worktree: Could not get last worktree: %s\n", err.Error())
				os.Exit(1)
			}

			core.SaveLastWorktree()
			fmt.Println(lastWorktree.String())
			os.Exit(0)
		}

		worktrees, err := git.GitWorktreeList()
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Could not get Git worktrees: %s\n", err.Error())
			os.Exit(1)
		}

		for _, wt := range worktrees {
			if strings.HasSuffix(wt.Path, "/"+args[0]) {
				core.SaveLastWorktree()
				fmt.Println(wt.Path)
				os.Exit(0)
			}
		}
		fmt.Fprintf(os.Stderr, "worktree: Could not find Git worktree: %s\n", args[0])
		os.Exit(1)
	},
}

// init registers the switch command
func init() {
	rootCmd.AddCommand(switchCmd)
}
