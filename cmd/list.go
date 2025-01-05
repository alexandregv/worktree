package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

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

		// Capture the tabbed output in a buffer
		var buf bytes.Buffer
		writer := tabwriter.NewWriter(&buf, 0, 0, 4, ' ', 0)

		// Loop through worktrees and write their formatted output
		for _, wt := range worktrees {
			var str string
			if wt.Bare {
				str = fmt.Sprintf("📁 %s\t🗳️ (bare)", strings.Replace(wt.Path, os.Getenv("HOME"), "~", -1))
			} else {
				str = fmt.Sprintf("📁 %s\t🔗 %s\t🔀 %s", strings.Replace(wt.Path, os.Getenv("HOME"), "~", -1), wt.Head[:7], wt.Branch)
			}
			if wt.Locked {
				str += "\t🔒" + wt.LockedReason
			}
			fmt.Fprintln(writer, str)
		}

		// Flush the writer to write the data into the buffer
		writer.Flush()

		fmt.Println(buf.String())
	},
}

// init registers the list command
func init() {
	rootCmd.AddCommand(listCmd)
}
