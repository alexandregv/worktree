package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a multi-worktrees repository",
	Long:  `Clone a repository (as bare) and its branches as local worktrees`,
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		flagAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Errors reading flags: %s\n", err.Error())
			os.Exit(1)
		}

		repoURL := args[0]

		var path string
		if len(args) >= 2 {
			path = args[1]
		} else {
			splits := strings.Split(strings.TrimRight(args[0], "/"), "/")
			if len(splits) == 0 {
				fmt.Fprintf(os.Stderr, "worktree: Invalid repository URL: %s\n", repoURL)
				os.Exit(1)
			}
			path = splits[len(splits)-1]
		}

		shellCmd := exec.Command("git", "clone", "--no-checkout", repoURL, path)
		shellCmd.Stdout = os.Stdout
		shellCmd.Stderr = os.Stderr
		err = shellCmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Error running git command `%s`: %s\n", shellCmd.String(), err.Error())
			os.Exit(1)
		}

		err = os.Chdir(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Error changing working directory: %s\n", err.Error())
			os.Exit(1)
		}

		shellCmd = exec.Command("git", "config", "core.bare", "true")
		shellCmd.Stdout = os.Stdout
		shellCmd.Stderr = os.Stderr
		err = shellCmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Error running git command `%s`: %s\n", shellCmd.String(), err.Error())
			os.Exit(1)
		}

		if flagAll {
			var sb strings.Builder
			shellCmd = exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/remotes/origin")
			shellCmd.Stdout = &sb
			shellCmd.Stderr = os.Stderr
			err = shellCmd.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "worktree: Error running git command `%s`: %s\n", shellCmd.String(), err.Error())
				os.Exit(1)
			}

			scanner := bufio.NewScanner(strings.NewReader(sb.String()))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				splits := strings.SplitN(scanner.Text(), "origin/", 2)
				if len(splits) <= 1 {
					continue
				}
				branch := splits[1]

				shellCmd = exec.Command("git", "worktree", "add", branch, "-B", branch, "origin/"+branch)
				shellCmd.Stdout = os.Stdout
				shellCmd.Stderr = os.Stderr
				err = shellCmd.Run()
				if err != nil {
					fmt.Fprintf(os.Stderr, "worktree: Error running git command `%s`: %s\n", shellCmd.String(), err.Error())
					os.Exit(1)
				}
			}
		} else {
			shellCmd = exec.Command("git", "remote", "show", "origin")
			out, err := shellCmd.Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "worktree: Error running git command `%s`: %s\n", shellCmd.String(), err.Error())
				os.Exit(1)
			}

			splits := strings.SplitN(string(out), "HEAD branch: ", 2)
			if len(splits) <= 1 {
				fmt.Fprint(os.Stderr, "worktree: Could not determine default branch\n")
				os.Exit(1)
			}
			defaultBranch := strings.Split(splits[1], "\n")[0]

			shellCmd = exec.Command("git", "worktree", "add", defaultBranch, "-B", defaultBranch, "origin/"+defaultBranch)
			shellCmd.Stdout = os.Stdout
			shellCmd.Stderr = os.Stderr
			err = shellCmd.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "worktree: Error running git command `%s`: %s\n", shellCmd.String(), err.Error())
				os.Exit(1)
			}
		}
	},
}

// init registers the list command
func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().BoolP("all", "a", false, "Create local worktrees for all branches")
}
