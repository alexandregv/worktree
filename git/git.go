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
		} else if wt.Detached {
			sb.WriteString("🔎 (detached)\t")
		} else {
			sb.WriteString("🔗 ")
			sb.WriteString(wt.Head[:7])
			sb.WriteString("\t🔀 ")
			sb.WriteString(strings.TrimPrefix(wt.Branch, "refs/heads/"))
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
	out, err := CommandOutput("worktree", "list", "--porcelain", "-z")
	if err != nil {
		return nil, err
	}

	return ParseWorktrees(out.String())
}

// NewWorktree adds a new worktree.
func NewWorktree(branch string) error {
	return Command("worktree", "add", branch, "-B", branch, "origin/"+branch)
}

// SetBare changes the core.bare config value.
func SetBare(bare bool) error {
	return Command("config", "core.bare", "true")
}

// Clone clones a repository
func Clone(url string, args ...string) error {
	fullArgs := []string{"clone", url}
	for _, arg := range args {
		fullArgs = append(fullArgs, arg)
	}

	return Command(fullArgs...)
}

// Refs returns specified refs, using `git for-each-ref --format='%(refname:short)'`.
func Refs(refs string) (heads string, err error) {
	out, err := CommandOutput("for-each-ref", "--format=%(refname:short)", "refs/"+refs)
	return out.String(), err
}

// Command executes the specified git command.
// Stdout and Stderr are redirected to os.Stdout and os.Stderr.
func Command(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run git command `%s`: %w", cmd.String(), err.Error())
	}
	return nil
}

// CommandOutput executes the specified git command and returns its output.
// Stderr is redirected to os.Stderr.
func CommandOutput(args ...string) (out strings.Builder, err error) {
	cmd := exec.Command("git", args...)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return out, fmt.Errorf("failed to run git command `%s`: %w", cmd.String(), err.Error())
	}
	return out, nil
}

// SetConfig sets the git config key/value pair
func SetConfig(key string, value string, args ...string) (err error) {
	fullArgs := []string{"config", "set", key, value}
	for _, arg := range args {
		fullArgs = append(fullArgs, arg)
	}
	return Command(fullArgs...)
}

// GetConfig gets the git config key/value pair
func GetConfig(key string, args ...string) (out strings.Builder, err error) {
	fullArgs := []string{"config", "get", key}
	for _, arg := range args {
		fullArgs = append(fullArgs, arg)
	}
	return CommandOutput(fullArgs...)
}
