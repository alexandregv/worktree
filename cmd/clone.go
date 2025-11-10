package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alexandregv/worktree/git"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a multi-worktrees repository",
	Long: `
	Clone a repository (as bare) and its branches as local worktrees.
	Native git args can be passed by adding '-- [--key=val]': wt clone <url> -- --depth=10

	This is the equivalent of running:
	  git clone --no-checkout <url> <path>
	  cd <path>
	  git config core.bare true
	  git worktree add $(basename <branch>) -B <branch> origin/<branch>  # for each branch if --all
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flagAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Errors reading flags: %s\n", err.Error())
			os.Exit(1)
		}

		repoURL := args[0]

		var path string
		var gitArgs []string
		switch len(args) {
		case 1: // wt clone <url>
			splits := strings.Split(strings.TrimRight(args[0], "/"), "/")
			if len(splits) == 0 {
				fmt.Fprintf(os.Stderr, "worktree: Invalid repository URL: %s\n", repoURL)
				os.Exit(1)
			}
			path = splits[len(splits)-1]
		case 2: // wt clone <url> [path]
			path = args[1]
		default: // wt clone <url> [path] [--] [git-args...]
			path = args[1]
			gitArgs = args[2:]
		}

		path, _ = strings.CutSuffix(path, "/")
		path, _ = strings.CutSuffix(path, ".git")

		err = git.Clone(repoURL, append([]string{path, "--no-checkout", "--no-single-branch"}, gitArgs...)...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Error cloning repository: %s\n", err.Error())
			os.Exit(1)
		}

		err = os.Chdir(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Error changing working directory: %s\n", err.Error())
			os.Exit(1)
		}

		err = git.SetBare(true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Error setting repository as bare: %s\n", err.Error())
			os.Exit(1)
		}

		if flagAll {
			refs, err := git.Refs("remotes/origin")
			if err != nil {
				fmt.Fprintf(os.Stderr, "worktree: Error listing refs/remotes/origin: %s\n", err.Error())
				os.Exit(1)
			}

			scanner := bufio.NewScanner(strings.NewReader(refs))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				splits := strings.SplitN(scanner.Text(), "origin/", 2)
				if len(splits) <= 1 {
					continue
				}
				branch := splits[1]

				err := git.NewWorktree(branch)
				if err != nil {
					fmt.Fprintf(os.Stderr, "worktree: Error adding worktree: %s\n", err.Error())
					os.Exit(1)
				}
			}
		} else {
			refs, err := git.Refs("heads")
			if err != nil {
				fmt.Fprintf(os.Stderr, "worktree: Error listings refs/heads: %s\n", err.Error())
				os.Exit(1)
			}

			scanner := bufio.NewScanner(strings.NewReader(refs))
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				defaultBranch := scanner.Text()

				err := git.NewWorktree(defaultBranch)
				if err != nil {
					fmt.Fprintf(os.Stderr, "worktree: Error adding worktree: %s\n", err.Error())
					os.Exit(1)
				}
			}
		}
	},
}

// init registers the clone command
func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().BoolP("all", "a", false, "Create local worktrees for all branches")
}
