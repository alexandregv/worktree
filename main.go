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
	fzfOptions, err := initFzfOptions(
		BuildWorktreeList(worktrees),
		[]string{
			"--height=40%",
			"--prompt=worktree: ",
			"--with-nth=2..",

			// Jump to a line with Space + <n>
			"--bind=space:jump",
			"--jump-labels=0123456789;:,<.>/?'\"!@#$%^&*",

			// Preview window with git log on the right
			// sh -c '' is used to expand $HOME
			"--preview=sh -c \"git -C {3} log --color=always --oneline --graph --decorate --all -n20\"",
			"--preview-window=right,45%",

			// Theme: Catppuccin Macchiato (https://github.com/catppuccin/fzf)
			"--color=bg+:#363a4f,bg:#24273a,spinner:#f4dbd6,hl:#ed8796",
			"--color=fg:#cad3f5,header:#ed8796,info:#c6a0f6,pointer:#f4dbd6",
			"--color=marker:#b7bdf8,fg+:#cad3f5,prompt:#c6a0f6,hl+:#ed8796",
			"--color=selected-bg:#494d64",
		},
	)
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
