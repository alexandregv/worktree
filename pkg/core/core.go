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

	// Capture the tabbed content as an array of strings
	strList := git.BuildWorktreeList(worktrees, true)
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
