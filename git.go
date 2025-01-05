package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Worktree represents a Git worktree with its attributes.
type Worktree struct {
	Path     string
	Branch   string
	Head     string
	Bare     bool
	Detached bool
	Locked   bool
}

// ParseWorktrees parses the output of `git worktree list --porcelain -z` into Worktree structs.
func ParseWorktrees(input string) ([]*Worktree, error) {
	var worktrees []*Worktree
	var currentWorktree *Worktree

	// Split each worktree group by double NUL character (--porcelain -z)
	entries := strings.Split(input, "\x00\x00")
	for _, entry := range entries {
		if entry == "" {
			continue
		}

		// Then each line (attribute) by simple NUL character
		lines := strings.Split(entry, "\x00")
		for _, line := range lines {
			switch {
			case strings.HasPrefix(line, "worktree"):
				currentWorktree = &Worktree{
					Path: lines[0][len("worktree "):],
				}
			case line == "bare":
				currentWorktree.Bare = true
			case line == "detached":
				currentWorktree.Detached = true
			case strings.HasPrefix(line, "locked"):
				currentWorktree.Locked = true
			case strings.HasPrefix(line, "HEAD "):
				currentWorktree.Head = line[len("HEAD "):]
			case strings.HasPrefix(line, "branch "):
				currentWorktree.Branch = line[len("branch "):]
			}
		}
		worktrees = append(worktrees, currentWorktree)
	}

	return worktrees, nil
}

// GitWorktreeList runs `git worktree list --porcelain -z` and parses the output.
func GitWorktreeList() ([]*Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain", "-z")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run git command: %w", err)
	}

	output := out.String()
	return ParseWorktrees(output)
}
