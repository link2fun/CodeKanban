package model

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestWorktreeServiceCreateAndRefresh(t *testing.T) {
	cleanup := initTestDB(t)
	defer cleanup()

	repoPath := createProjectTestRepo(t)
	projectService := &ProjectService{}
	project, err := projectService.CreateProject(context.Background(), CreateProjectParams{
		Name: "WT Project",
		Path: repoPath,
	})
	if err != nil {
		t.Fatalf("create project failed: %v", err)
	}

	service := NewWorktreeService()
	service.AsyncRefresh(false)
	ctx := context.Background()

	worktree, err := service.CreateWorktree(ctx, project.Id, "feature/testing", "main", true)
	if err != nil {
		t.Fatalf("CreateWorktree returned error: %v", err)
	}
	if worktree.Id == "" {
		t.Fatalf("expected worktree to have ID")
	}
	if worktree.Path == "" {
		t.Fatalf("expected worktree path to be set")
	}

	if _, err := os.Stat(worktree.Path); err != nil {
		t.Fatalf("git did not create worktree path: %v", err)
	}

	got, err := service.GetWorktree(ctx, worktree.Id)
	if err != nil {
		t.Fatalf("GetWorktree failed: %v", err)
	}
	if got.BranchName != "feature/testing" {
		t.Fatalf("expected branch name feature/testing, got %s", got.BranchName)
	}

	refreshed, err := service.RefreshWorktreeStatus(ctx, worktree.Id)
	if err != nil {
		t.Fatalf("RefreshWorktreeStatus failed: %v", err)
	}
	if refreshed.StatusUpdatedAt == nil {
		t.Fatalf("expected status updated timestamp to be set")
	}
}

func TestWorktreeServiceDeleteAndSync(t *testing.T) {
	cleanup := initTestDB(t)
	defer cleanup()

	repoPath := createProjectTestRepo(t)
	projectService := &ProjectService{}
	project, err := projectService.CreateProject(context.Background(), CreateProjectParams{
		Name: "Delete Project",
		Path: repoPath,
	})
	if err != nil {
		t.Fatalf("create project failed: %v", err)
	}

	if runtime.GOOS == "windows" {
		t.Skip("DeleteWorktree integration test is flaky on Windows due to git worktree remove permissions")
	}

	service := NewWorktreeService()
	service.AsyncRefresh(false)
	ctx := context.Background()

	worktree, err := service.CreateWorktree(ctx, project.Id, "feature/delete", "main", true)
	if err != nil {
		t.Fatalf("CreateWorktree returned error: %v", err)
	}

	if err := service.DeleteWorktree(ctx, worktree.Id, true, true); err != nil {
		t.Fatalf("DeleteWorktree returned error: %v", err)
	}

	if _, err := service.GetWorktree(ctx, worktree.Id); err == nil {
		t.Fatalf("expected worktree to be deleted")
	}

	// Manually create a worktree using git to test sync
	runGitCommand(t, repoPath, "worktree", "add", filepath.Join(repoPath, "manual"), "main")
	defer runGitCommand(t, repoPath, "worktree", "remove", filepath.Join(repoPath, "manual"))

	if err := service.SyncWorktrees(ctx, project.Id); err != nil {
		t.Fatalf("SyncWorktrees returned error: %v", err)
	}

	worktrees, err := service.ListWorktrees(ctx, project.Id)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}
	found := false
	for _, wt := range worktrees {
		if wt.Path == filepath.Join(repoPath, "manual") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected manual worktree to be synced into database")
	}
}

func TestWorktreeServiceRefreshAll(t *testing.T) {
	cleanup := initTestDB(t)
	defer cleanup()

	repoPath := createProjectTestRepo(t)
	projectService := &ProjectService{}
	project, err := projectService.CreateProject(context.Background(), CreateProjectParams{
		Name: "Refresh Project",
		Path: repoPath,
	})
	if err != nil {
		t.Fatalf("create project failed: %v", err)
	}

	service := NewWorktreeService()
	service.AsyncRefresh(false)
	ctx := context.Background()

	if _, err := service.CreateWorktree(ctx, project.Id, "feature/all", "main", true); err != nil {
		t.Fatalf("CreateWorktree returned error: %v", err)
	}

	updated, failed, err := service.RefreshAllWorktrees(ctx, project.Id)
	if err != nil {
		t.Fatalf("RefreshAllWorktrees returned error: %v", err)
	}
	if updated == 0 || failed != 0 {
		t.Fatalf("unexpected refresh counts updated=%d failed=%d", updated, failed)
	}

	// Ensure timestamps updated
	time.Sleep(10 * time.Millisecond)
	updatedWTs, err := service.ListWorktrees(ctx, project.Id)
	if err != nil {
		t.Fatalf("ListWorktrees failed: %v", err)
	}
	if len(updatedWTs) == 0 || updatedWTs[0].StatusUpdatedAt == nil {
		t.Fatalf("expected worktree status to be refreshed")
	}
}
