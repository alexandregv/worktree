package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	fzf "github.com/junegunn/fzf/src"
	fzfProtector "github.com/junegunn/fzf/src/protector"
)

func main() {
	worktrees, err := GitWorktreeList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not get Git worktrees: %s\n", err.Error())
		os.Exit(1)
	}

	// Init fzf options with defaults + our custom values
	fzfOptions, err := InitFzfOptions(BuildWorktreeList(worktrees))
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not initalize fzf options: %s\n", err.Error())
		os.Exit(fzf.ExitError)
	}

	// Capture fzf output, get corresponding worktree by index
	go func() {
		out := <-fzfOptions.Output
		i, err := strconv.Atoi(strings.Split(out, ":")[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "worktree: Error parsing fzf output: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Println(worktrees[i].Path)
	}()

	// Run fzf (with BSD protector)
	fzfProtector.Protect()
	exitCode, err := fzf.Run(fzfOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not run fzf: %s\n", err.Error())
		os.Exit(exitCode)
	}
}
