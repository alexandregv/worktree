package main

import (
	"fmt"
	"os"
	"reflect"

	fzf "github.com/junegunn/fzf/src"
	fzfProtector "github.com/junegunn/fzf/src/protector"

	git "github.com/go-git/go-git/v5"
)

func getGitWorktrees() (paths []string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// cwd = "/var/home/reach/Documents/gotemplate/"

	repo, err := git.PlainOpenWithOptions(cwd, &git.PlainOpenOptions{DetectDotGit: true, EnableDotGitCommonDir: false})
	if err != nil {
		return nil, err
	}

	worktrees, err := repo.Worktrees()
	// for _, wt := range worktrees {
	// 	fileInfo, err := wt.Filesystem.Stat("")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	godump.Dump(fileInfo)
	// 	fmt.Println("\n---\n")
	// }

	// Keep only paths as strings
	paths = []string{}
	for _, wt := range worktrees {
		paths = append(paths, fmt.Sprintf("%s", reflect.ValueOf(wt.Filesystem).Elem().FieldByName("base")))
	}

	return paths, err
}

func initFzfOptions(inputs []string, customOptions []string) (options *fzf.Options, err error) {
	options, err = fzf.ParseOptions(true, customOptions)
	if err != nil {
		return nil, err
	}

	options.Input = make(chan string)

	go func() {
		defer close(options.Input)

		for _, input := range inputs {
			options.Input <- input
		}
	}()

	return options, err
}

func main() {
	// Get Git worktrees
	worktrees, err := getGitWorktrees()
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not get Git worktrees: %s\n", err.Error())
		os.Exit(1)
	}

	// Init fzf options with defaults + our custom values
	fzfOptions, err := initFzfOptions(
		worktrees,
		[]string{
			"--height=10",
			"--prompt=worktree: ",
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not initalize fzf options: %s\n", err.Error())
		os.Exit(fzf.ExitError)
	}

	// Run fzf (with BSD protector)
	fzfProtector.Protect()
	exitCode, err := fzf.Run(fzfOptions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not start fzf: %s\n", err.Error())
		os.Exit(exitCode)
	}
}
