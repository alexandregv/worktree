package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
)

// Worktree represents a Git worktree with its attributes.
type Worktree struct {
	Path         string
	Branch       string
	Head         string
	Bare         bool
	Detached     bool
	Locked       bool
	LockedReason string
}

// BuildWorktreeList formats a tabulated worktree list, with emojis and optional indexes
func BuildWorktreeList(worktrees []*Worktree, withIndexes bool) (list string) {
	var sb strings.Builder
	var buf bytes.Buffer
	writer := tabwriter.NewWriter(&buf, 0, 0, 4, ' ', 0)

	for i, wt := range worktrees {
		path := strings.Replace(wt.Path, os.Getenv("HOME"), "$HOME", -1)

		if withIndexes {
			sb.WriteString(fmt.Sprintf("%d: ", i))
		}

		sb.WriteString(fmt.Sprintf("📁 %s\t", path))

		if wt.Bare {
			sb.WriteString("🗳️ (bare)\t")
		} else {
			sb.WriteString(fmt.Sprintf("🔗 %s\t🔀 %s\t", wt.Head[:7], strings.Replace(wt.Branch, "refs/heads/", "", -1)))
		}

		if wt.Locked {
			sb.WriteString("\t🔒" + wt.LockedReason + "\t")
		}

		fmt.Fprintln(writer, sb.String())
		sb.Reset()
	}

	writer.Flush()

	return buf.String()
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
				currentWorktree.LockedReason = line[len("locked"):]
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
