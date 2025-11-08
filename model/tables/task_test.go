package tables

import "testing"

func TestTaskTablePersistsTags(t *testing.T) {
	db := setupTestDB(t)

	project := &ProjectTable{
		Name:          "Task Project",
		Path:          "/tmp/task-project",
		DefaultBranch: "main",
	}
	if err := db.Create(project).Error; err != nil {
		t.Fatalf("create project failed: %v", err)
	}

	worktree := &WorktreeTable{
		ProjectID:  project.ID,
		BranchName: "main",
		Path:       "/tmp/task-project",
		IsMain:     true,
	}
	if err := db.Create(worktree).Error; err != nil {
		t.Fatalf("create worktree failed: %v", err)
	}

	worktreeID := worktree.ID
	task := &TaskTable{
		ProjectID:   project.ID,
		WorktreeID:  &worktreeID,
		Title:       "Implement kanban board",
		Description: "Initial task to validate table",
		Status:      "todo",
		Priority:    1,
		OrderIndex:  1.0,
		Tags:        StringArray{"backend", "high-priority"},
	}

	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create task failed: %v", err)
	}

	var loaded TaskTable
	if err := db.First(&loaded, "id = ?", task.ID).Error; err != nil {
		t.Fatalf("load task failed: %v", err)
	}

	if len(loaded.Tags) != len(task.Tags) {
		t.Fatalf("expected %d tags, got %d", len(task.Tags), len(loaded.Tags))
	}
	for i, tag := range task.Tags {
		if loaded.Tags[i] != tag {
			t.Fatalf("tag at index %d mismatch: want %q got %q", i, tag, loaded.Tags[i])
		}
	}
}

func TestStringArrayScanNil(t *testing.T) {
	var arr StringArray
	if err := arr.Scan(nil); err != nil {
		t.Fatalf("scan nil failed: %v", err)
	}
	if len(arr) != 0 {
		t.Fatalf("expected empty slice for nil scan, got %v", arr)
	}
}
