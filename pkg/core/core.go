package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexandregv/worktree/pkg/fzf"
	"github.com/alexandregv/worktree/pkg/git"

	fzfLib "github.com/junegunn/fzf/src"
	fzfLibProtec "github.com/junegunn/fzf/src/protector"
)

func OpenTUI() {
	worktrees, err := git.GitWorktreeList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not get Git worktrees: %s\n", err.Error())
		os.Exit(1)
	}

	// Init fzf options with defaults + our custom values
	fzfOptions, err := fzf.InitFzfOptions(git.BuildWorktreeList(worktrees))
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not initalize fzf options: %s\n", err.Error())
		os.Exit(fzfLib.ExitError)
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
	fzfLibProtec.Protect()
	exitCode, err := fzfLib.Run(fzfOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not run fzf: %s\n", err.Error())
		os.Exit(exitCode)
	}
}
