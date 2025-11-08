package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	goGit "github.com/go-git/go-git/v5"
)

// GitRepo is a thin wrapper around go-git's Repository with resolved path metadata.
type GitRepo struct {
	Path       string
	Repository *goGit.Repository
}

// Remote describes a configured git remote.
type Remote struct {
	Name string
	URL  string
}

var (
	errEmptyPath = errors.New("path is required")
)

// DetectRepository returns a GitRepo if the given path is a valid git repository.
func DetectRepository(path string) (*GitRepo, error) {
	p := strings.TrimSpace(path)
	if p == "" {
		return nil, errEmptyPath
	}

	absPath, err := filepath.Abs(p)
	if err != nil {
		return nil, fmt.Errorf("resolve git path: %w", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		return nil, fmt.Errorf("stat git path: %w", err)
	}

	repo, err := goGit.PlainOpen(absPath)
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	return &GitRepo{
		Path:       absPath,
		Repository: repo,
	}, nil
}

// GetRemotes lists configured remotes, returning the first URL for each remote.
func (r *GitRepo) GetRemotes() ([]Remote, error) {
	if r == nil || r.Repository == nil {
		return nil, errors.New("git repository is not initialized")
	}

	remotes, err := r.Repository.Remotes()
	if err != nil {
		return nil, err
	}

	items := make([]Remote, 0, len(remotes))
	for _, remote := range remotes {
		cfg := remote.Config()
		if len(cfg.URLs) == 0 {
			continue
		}
		items = append(items, Remote{
			Name: cfg.Name,
			URL:  cfg.URLs[0],
		})
	}

	return items, nil
}

// GetCurrentBranch returns the short name of HEAD's branch.
func (r *GitRepo) GetCurrentBranch() (string, error) {
	if r == nil || r.Repository == nil {
		return "", errors.New("git repository is not initialized")
	}

	head, err := r.Repository.Head()
	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}
