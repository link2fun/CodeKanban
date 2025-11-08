package tables

import "testing"

func TestWorktreeTableAssociations(t *testing.T) {
	db := setupTestDB(t)

	project := &ProjectTable{
		Name:          "WT Project",
		Path:          "/tmp/wt-project",
		DefaultBranch: "main",
	}
	if err := db.Create(project).Error; err != nil {
		t.Fatalf("create project failed: %v", err)
	}

	worktree := &WorktreeTable{
		ProjectID:  project.ID,
		BranchName: "feature/demo",
		Path:       "/tmp/wt-project-worktrees/demo",
		IsMain:     false,
		HeadCommit: "abcdef1",
	}

	if err := db.Create(worktree).Error; err != nil {
		t.Fatalf("create worktree failed: %v", err)
	}

	var loaded WorktreeTable
	if err := db.Preload("Project").First(&loaded, "id = ?", worktree.ID).Error; err != nil {
		t.Fatalf("load worktree failed: %v", err)
	}
	if loaded.Project == nil {
		t.Fatalf("expected associated project to be eagerly loaded")
	}
	if loaded.Project.ID != project.ID {
		t.Fatalf("expected project id %s got %s", project.ID, loaded.Project.ID)
	}
}
