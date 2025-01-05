package main

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	fzf "github.com/junegunn/fzf/src"
	fzfProtector "github.com/junegunn/fzf/src/protector"
)

func buildList(worktrees []*Worktree) (list []string) {
	// Capture the tabbed output in a buffer
	var buf bytes.Buffer
	writer := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Loop through worktrees and write their formatted output
	for i, wt := range worktrees {
		var str string
		if wt.Bare {
			str = fmt.Sprintf("%d: %s (bare)", i, wt.Path)
		} else {
			str = fmt.Sprintf("%d: %s\t[%s]\t%s", i, wt.Path, wt.Head[:7], wt.Branch)
		}
		if wt.Locked {
			str += "\tlocked"
		}
		fmt.Fprintln(writer, str)
	}

	// Flush the writer to write the data into the buffer
	writer.Flush()

	// Capture the tabbed content as an array of strings
	lines := bytes.Split(buf.Bytes(), []byte("\n"))
	for _, line := range lines {
		if string(line) == "" {
			continue
		}
		list = append(list, string(line))
	}
	return list
}

func main() {
	worktrees, err := GitWorktreeList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not get Git worktrees: %s\n", err.Error())
		os.Exit(1)
	}

	// Init fzf options with defaults + our custom values
	fzfOptions, err := initFzfOptions(
		buildList(worktrees),
		[]string{
			"--height=40%",
			"--prompt=worktree: ",
			"--with-nth=2,3,4,5",
			"--bind=space:jump",
			"--jump-labels=0123456789",
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "worktree: Could not initalize fzf options: %s\n", err.Error())
		os.Exit(fzf.ExitError)
	}

	go func() {
		out := <-fzfOptions.Output
		i, err := strconv.Atoi(strings.Split(out, ":")[0])
		if err != nil {
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
