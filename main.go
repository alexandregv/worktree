package main

import (
	"fmt"
	"os"
	"reflect"

	fzf "github.com/junegunn/fzf/src"
	fzfProtector "github.com/junegunn/fzf/src/protector"

	git "github.com/go-git/go-git/v5"
)

func getGitRepo(path string) (repo *git.Repository, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repo, err = git.PlainOpenWithOptions(cwd, &git.PlainOpenOptions{DetectDotGit: true, EnableDotGitCommonDir: true})
	if err != nil {
		return nil, err
	}
	// godump.Dump(repo)

	return repo, err
}

func getRootRepo(repo *git.Repository) (commonDirRepo *git.Repository, err error) {
	commonDotGitPath := fmt.Sprintf("%+v\n",
		reflect.ValueOf(repo.Storer).
			Elem().FieldByName("fs").
			Elem().Elem().FieldByName("commonDotGitFs").
			Elem().Elem().FieldByName("base"),
	)
	// fmt.Println(commonDotGitPath)
	// fmt.Println(commonDotGitPath[:len(commonDotGitPath)-7])

	commonDirRepo, err = getGitRepo(commonDotGitPath[:len(commonDotGitPath)-7])
	return commonDirRepo, err
}

func getGitWorktrees(repo *git.Repository) (paths []string, err error) {
	worktrees, err := repo.Worktrees()

	// Keep only paths as strings
	paths = []string{}
	for _, wt := range worktrees {
		paths = append(paths, fmt.Sprintf("%s", reflect.ValueOf(wt.Filesystem).Elem().FieldByName("base")))
		fmt.Println(fmt.Sprintf("%s", reflect.ValueOf(wt.Filesystem).Elem().FieldByName("base")))
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
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not get current working directory: %s\n", err.Error())
		os.Exit(1)
	}

	repo, err := getGitRepo(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not get Git repository: %s\n", err.Error())
		os.Exit(1)
	}

	rrepo, err := getRootRepo(repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not get root Git repository: %s\n", err.Error())
		os.Exit(1)
	}

	worktrees, err := getGitWorktrees(rrepo)
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
