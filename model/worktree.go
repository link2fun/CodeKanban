package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go-template/utils"
	"go-template/utils/git"

	"go.uber.org/zap"
)

var (
	// ErrWorktreeNotFound indicates the requested worktree does not exist.
	ErrWorktreeNotFound = errors.New("worktree not found")
	// ErrWorktreeIsMain indicates the worktree is the main repository path and cannot be removed.
	ErrWorktreeIsMain = errors.New("cannot delete main worktree")
	// ErrWorktreeHasTasks indicates there are tasks referencing the worktree requiring a force delete.
	ErrWorktreeHasTasks = errors.New("worktree has active tasks")
)

// WorktreeService coordinates CRUD operations between git worktrees and the database.
type WorktreeService struct {
	asyncStatusRefresh bool
}

// NewWorktreeService builds a WorktreeService with async status refresh enabled.
func NewWorktreeService() *WorktreeService {
	return &WorktreeService{
		asyncStatusRefresh: true,
	}
}

// AsyncRefresh toggles async status refresh behaviour (useful for tests).
func (s *WorktreeService) AsyncRefresh(enabled bool) {
	if s == nil {
		return
	}
	s.asyncStatusRefresh = enabled
}

// CreateWorktree provisions a new git worktree and persists its metadata.
func (s *WorktreeService) CreateWorktree(
	ctx context.Context,
	projectID string,
	branchName string,
	baseBranch string,
	createBranch bool,
) (*Worktree, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	q, err := resolveQueries(nil)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(projectID) == "" {
		return nil, fmt.Errorf("project id is required")
	}
	if strings.TrimSpace(branchName) == "" {
		return nil, fmt.Errorf("branch name is required")
	}

	project, err := q.ProjectGetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorktreeNotFound
		}
		return nil, err
	}

	gitRepo, err := git.DetectRepository(project.Path)
	if err != nil {
		return nil, err
	}

	targetBranch := strings.TrimSpace(branchName)
	if createBranch {
		refBranch := strings.TrimSpace(baseBranch)
		if refBranch == "" {
			if project.DefaultBranch != "" {
				refBranch = project.DefaultBranch
			}
		}
		if refBranch == "" {
			refBranch = "main"
		}
		if err := gitRepo.CreateBranch(targetBranch, refBranch); err != nil {
			return nil, err
		}
	}

	worktreePath, err := s.resolveWorktreePath(project, targetBranch)
	if err != nil {
		return nil, err
	}

	if err := gitRepo.AddWorktree(worktreePath, targetBranch, false); err != nil {
		return nil, err
	}

	now := time.Now()
	worktree, err := q.WorktreeCreate(ctx, &WorktreeCreateParams{
		Id:              utils.NewID(),
		CreatedAt:       now,
		UpdatedAt:       now,
		ProjectId:       projectID,
		BranchName:      targetBranch,
		Path:            worktreePath,
		IsMain:          false,
		IsBare:          false,
		HeadCommit:      nil,
		StatusAhead:     0,
		StatusBehind:    0,
		StatusModified:  0,
		StatusStaged:    0,
		StatusUntracked: 0,
		StatusConflicts: 0,
		StatusUpdatedAt: nil,
	})
	if err != nil {
		_ = gitRepo.RemoveWorktree(worktreePath, true)
		return nil, err
	}

	if s != nil && s.asyncStatusRefresh {
		go s.RefreshWorktreeStatus(context.Background(), worktree.Id)
	} else {
		_, _ = s.RefreshWorktreeStatus(ctx, worktree.Id)
	}

	return worktree, nil
}

// ListWorktrees returns worktrees for a project ordered by main flag then creation.
func (s *WorktreeService) ListWorktrees(ctx context.Context, projectID string) ([]*Worktree, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	q, err := resolveQueries(nil)
	if err != nil {
		return nil, err
	}

	return q.WorktreeListByProject(ctx, projectID)
}

// GetWorktree fetches a worktree by identifier.
func (s *WorktreeService) GetWorktree(ctx context.Context, id string) (*Worktree, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	q, err := resolveQueries(nil)
	if err != nil {
		return nil, err
	}

	wt, err := q.WorktreeGetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWorktreeNotFound
		}
		return nil, err
	}
	return wt, nil
}

// DeleteWorktree removes a worktree from git and the database.
func (s *WorktreeService) DeleteWorktree(ctx context.Context, id string, force, deleteBranch bool) error {
	if ctx == nil {
		ctx = context.Background()
	}

	q, err := resolveQueries(nil)
	if err != nil {
		return err
	}

	worktree, err := s.GetWorktree(ctx, id)
	if err != nil {
		return err
	}
	if worktree.IsMain {
		return ErrWorktreeIsMain
	}

	worktreeID := worktree.Id
	taskCount, err := q.TaskCountByWorktree(ctx, &worktreeID)
	if err != nil {
		return err
	}
	if taskCount > 0 && !force {
		return ErrWorktreeHasTasks
	}

	project, err := q.ProjectGetByID(ctx, worktree.ProjectId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrWorktreeNotFound
		}
		return err
	}

	gitRepo, err := git.DetectRepository(project.Path)
	if err != nil {
		return err
	}

	if err := gitRepo.RemoveWorktree(worktree.Path, force); err != nil {
		return err
	}

	if deleteBranch {
		if err := gitRepo.DeleteBranch(worktree.BranchName, force); err != nil {
			utils.Logger().Warn("failed to delete branch",
				zap.Error(err),
				zap.String("branch", worktree.BranchName),
				zap.String("projectId", project.Id),
			)
		}
	}

	now := time.Now()
	_, err = q.WorktreeSoftDelete(ctx, &WorktreeSoftDeleteParams{
		DeletedAt: &now,
		UpdatedAt: now,
		Id:        id,
	})
	return err
}

// RefreshWorktreeStatus updates cached status fields for a worktree and returns the refreshed record.
func (s *WorktreeService) RefreshWorktreeStatus(ctx context.Context, id string) (*Worktree, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	q, err := resolveQueries(nil)
	if err != nil {
		return nil, err
	}

	worktree, err := s.GetWorktree(ctx, id)
	if err != nil {
		return nil, err
	}

	status, err := git.GetWorktreeStatus(worktree.Path)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var headPtr *string
	if status.LastCommit != nil {
		head := status.LastCommit.SHA
		headPtr = &head
	}

	updated, err := q.WorktreeUpdateStatus(ctx, &WorktreeUpdateStatusParams{
		UpdatedAt:       now,
		StatusAhead:     int64(status.Ahead),
		StatusBehind:    int64(status.Behind),
		StatusModified:  int64(status.Modified),
		StatusStaged:    int64(status.Staged),
		StatusUntracked: int64(status.Untracked),
		StatusConflicts: int64(status.Conflicted),
		StatusUpdatedAt: &now,
		HeadCommit:      headPtr,
		Id:              worktree.Id,
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// RefreshAllWorktrees refreshes status for every worktree belonging to a project.
func (s *WorktreeService) RefreshAllWorktrees(ctx context.Context, projectID string) (updated, failed int, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	worktrees, err := s.ListWorktrees(ctx, projectID)
	if err != nil {
		return 0, 0, err
	}

	for _, wt := range worktrees {
		if _, err := s.RefreshWorktreeStatus(ctx, wt.Id); err != nil {
			failed++
			utils.Logger().Warn("failed to refresh worktree status",
				zap.Error(err),
				zap.String("worktreeId", wt.Id),
				zap.String("projectId", projectID),
			)
		} else {
			updated++
		}
	}
	return updated, failed, nil
}

// SyncWorktrees ensures git worktrees and the database remain aligned.
func (s *WorktreeService) SyncWorktrees(ctx context.Context, projectID string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	q, err := resolveQueries(nil)
	if err != nil {
		return err
	}

	project, err := q.ProjectGetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrWorktreeNotFound
		}
		return err
	}

	gitRepo, err := git.DetectRepository(project.Path)
	if err != nil {
		return err
	}

	gitWorktrees, err := gitRepo.ListWorktrees()
	if err != nil {
		return err
	}

	dbWorktrees, err := s.ListWorktrees(ctx, projectID)
	if err != nil {
		return err
	}

	gitByPath := make(map[string]git.WorktreeInfo, len(gitWorktrees))
	for _, wt := range gitWorktrees {
		gitByPath[normalizePathCase(wt.Path)] = wt
	}

	dbByPath := make(map[string]*Worktree, len(dbWorktrees))
	for _, wt := range dbWorktrees {
		dbByPath[normalizePathCase(wt.Path)] = wt
	}

	now := time.Now()
	for normPath, gitWT := range gitByPath {
		if existing, ok := dbByPath[normPath]; ok {
			var headPtr *string
			if gitWT.HeadCommit != "" {
				commit := gitWT.HeadCommit
				headPtr = &commit
			}
			if err := q.WorktreeUpdateMetadata(ctx, &WorktreeUpdateMetadataParams{
				UpdatedAt:  now,
				BranchName: gitWT.Branch,
				HeadCommit: headPtr,
				IsMain:     gitWT.IsMain,
				IsBare:     gitWT.IsBare,
				Id:         existing.Id,
			}); err != nil {
				return err
			}
			continue
		}

		var headPtr *string
		if gitWT.HeadCommit != "" {
			commit := gitWT.HeadCommit
			headPtr = &commit
		}
		if _, err := q.WorktreeCreate(ctx, &WorktreeCreateParams{
			Id:              utils.NewID(),
			CreatedAt:       now,
			UpdatedAt:       now,
			ProjectId:       projectID,
			BranchName:      gitWT.Branch,
			Path:            filepath.Clean(gitWT.Path),
			IsMain:          gitWT.IsMain,
			IsBare:          gitWT.IsBare,
			HeadCommit:      headPtr,
			StatusAhead:     0,
			StatusBehind:    0,
			StatusModified:  0,
			StatusStaged:    0,
			StatusUntracked: 0,
			StatusConflicts: 0,
			StatusUpdatedAt: nil,
		}); err != nil {
			return err
		}
	}

	for normPath, dbWT := range dbByPath {
		if _, ok := gitByPath[normPath]; ok {
			continue
		}
		if _, err := q.WorktreeSoftDelete(ctx, &WorktreeSoftDeleteParams{
			DeletedAt: &now,
			UpdatedAt: now,
			Id:        dbWT.Id,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *WorktreeService) resolveWorktreePath(project *Project, branchName string) (string, error) {
	basePath := ""
	if project.WorktreeBasePath != nil && strings.TrimSpace(*project.WorktreeBasePath) != "" {
		basePath = *project.WorktreeBasePath
	} else {
		basePath = filepath.Join(project.Path, "worktrees")
	}
	if !filepath.IsAbs(basePath) {
		basePath = filepath.Join(project.Path, basePath)
	}
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return "", err
	}

	dirName := sanitizeBranchName(branchName)
	return filepath.Join(basePath, dirName), nil
}

func sanitizeBranchName(branch string) string {
	replacer := strings.NewReplacer(
		"/", "__",
		"\\", "__",
		":", "_",
		"*", "_",
		"?", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(strings.TrimSpace(branch))
}

func normalizePathCase(path string) string {
	clean := filepath.Clean(path)
	if runtime.GOOS == "windows" {
		return strings.ToLower(clean)
	}
	return clean
}
