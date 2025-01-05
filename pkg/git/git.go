package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
func BuildWorktreeList(worktrees []*Worktree, withIndexes bool, replaceHome string) (list string) {
	var sb strings.Builder
	var buf bytes.Buffer
	writer := tabwriter.NewWriter(&buf, 0, 0, 4, ' ', 0)

	for i, wt := range worktrees {
		if withIndexes {
			sb.WriteString(strconv.Itoa(i))
			sb.WriteString(": ")
		}

		sb.WriteString("📁 ")
		if replaceHome != "" {
			sb.WriteString(strings.Replace(wt.Path, os.Getenv("HOME"), replaceHome, -1))
		} else {
			sb.WriteString(wt.Path)
		}
		sb.WriteString("\t")

		if wt.Bare {
			sb.WriteString("🗳️ (bare)\t")
		} else {
			sb.WriteString("🔗 ")
			sb.WriteString(wt.Head[:7])
			sb.WriteString("\t🔀 ")
			sb.WriteString(wt.Branch[len("refs/heads/"):]) // Slice to remove the "refs/heads/" prefix
			sb.WriteString("\t")

		}

		if wt.Locked {
			sb.WriteString("\t🔒 ")
			sb.WriteString(wt.LockedReason)
			sb.WriteString("\t")
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
