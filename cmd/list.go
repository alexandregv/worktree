package cmd

import (
	"fmt"
	"os"

	"github.com/alexandregv/worktree/pkg/git"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "list worktrees",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		worktrees, err := git.GitWorktreeList()
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Could not get Git worktrees: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println(git.BuildWorktreeList(worktrees, false, "~"))
	},
}

// init registers the list command
func init() {
	rootCmd.AddCommand(listCmd)
}
