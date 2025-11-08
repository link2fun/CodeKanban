package git

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	goGit "github.com/go-git/go-git/v5"
)

// WorktreeStatus aggregates repository state insights for a worktree.
type WorktreeStatus struct {
	Branch     string
	Ahead      int
	Behind     int
	Modified   int
	Staged     int
	Untracked  int
	Conflicted int
	LastCommit *CommitInfo
}

// CommitInfo describes a git commit summary.
type CommitInfo struct {
	SHA     string
	Message string
	Author  string
	Date    time.Time
}

// GetWorktreeStatus gathers branch, diff, and status metrics for a worktree path.
func GetWorktreeStatus(path string) (*WorktreeStatus, error) {
	repo, err := goGit.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	status := &WorktreeStatus{}

	if head, err := repo.Head(); err == nil {
		status.Branch = head.Name().Short()
		if status.Branch == "" || status.Branch == "HEAD" {
			status.Branch = describeBranch(path)
		}
		if commit, err := repo.CommitObject(head.Hash()); err == nil {
			status.LastCommit = &CommitInfo{
				SHA:     shortCommit(commit.Hash.String()),
				Message: firstLine(commit.Message),
				Author:  commit.Author.Name,
				Date:    commit.Author.When,
			}
		}
	} else {
		status.Branch = describeBranch(path)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	snap, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	for _, fs := range snap {
		if fs.Staging == goGit.Untracked || fs.Worktree == goGit.Untracked {
			status.Untracked++
			continue
		}
		if fs.Staging == goGit.UpdatedButUnmerged || fs.Worktree == goGit.UpdatedButUnmerged {
			status.Conflicted++
			continue
		}
		switch fs.Worktree {
		case goGit.Modified, goGit.Added, goGit.Deleted, goGit.Renamed:
			status.Modified++
		}
		if fs.Staging != goGit.Unmodified && fs.Staging != goGit.Untracked {
			status.Staged++
		}
	}

	status.Ahead, status.Behind = getAheadBehind(path)
	return status, nil
}

// GetWorktreeStatus returns the status for the provided worktree path. When
// path is empty the receiver's repository path is used.
func (r *GitRepo) GetWorktreeStatus(path string) (*WorktreeStatus, error) {
	if r == nil {
		return nil, errors.New("git repository is not initialized")
	}
	target := strings.TrimSpace(path)
	if target == "" {
		target = r.Path
	}
	return GetWorktreeStatus(target)
}

func getAheadBehind(path string) (ahead, behind int) {
	cmd := exec.Command("git", "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return 0, 0
	}

	parts := strings.Fields(string(output))
	if len(parts) >= 2 {
		ahead, _ = strconv.Atoi(parts[0])
		behind, _ = strconv.Atoi(parts[1])
	}
	return ahead, behind
}

func shortCommit(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

func firstLine(msg string) string {
	if idx := strings.Index(msg, "\n"); idx >= 0 {
		return strings.TrimSpace(msg[:idx])
	}
	return strings.TrimSpace(msg)
}

func describeBranch(path string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = path
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	name := strings.TrimSpace(string(output))
	if name == "HEAD" {
		return ""
	}
	return name
}
