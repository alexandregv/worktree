package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	fzfLib "github.com/junegunn/fzf/src"
	fzfLibProtec "github.com/junegunn/fzf/src/protector"

	"github.com/alexandregv/worktree/fzf"
	"github.com/alexandregv/worktree/git"
)

func OpenTUI() {
	worktrees, err := git.GitWorktreeList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not get Git worktrees: %s\n", err.Error())
		os.Exit(1)
	}

	// Capture the tabbed content as an array of strings
	strList := git.BuildWorktreeList(worktrees, true, "$HOME")
	arrayList := []string{}
	for _, line := range strings.Split(strList, "\n") {
		if string(line) == "" {
			continue
		}
		arrayList = append(arrayList, line)
	}

	// Init fzf options with defaults + our custom values
	fzfOptions, err := fzf.InitFzfOptions(arrayList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not initalize fzf options: %s\n", err.Error())
		os.Exit(fzfLib.ExitError)
	}

	// Run fzf (with BSD protector)
	fzfLibProtec.Protect()
	exitCode, err := fzfLib.Run(fzfOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not run fzf: %s\n", err.Error())
		os.Exit(exitCode)
	}

	switch exitCode {
	case 0:
		i, err := strconv.Atoi(strings.Split(<-fzfOptions.Output, ":")[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Error parsing fzf output: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(worktrees[i].Path)
		SaveLastWorktree()
	case 130:
		fmt.Println("^C")
		os.Exit(130)
	default:
		fmt.Fprintf(os.Stderr, "worktree: Could not run fzf: %s\n", err.Error())
		os.Exit(exitCode)
	}
}

func SaveLastWorktree() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not get current working directory to save last worktree (wt switch -): %s\n", err.Error())
	}
	err = git.SetConfig("alexandregv-worktree.lastworktree", cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not save last worktree to git config (wt switch -): %s\n", err.Error())
	}
}
