package tables

import "testing"

func TestTaskCommentTable(t *testing.T) {
	db := setupTestDB(t)

	project := &ProjectTable{
		Name:          "Comment Project",
		Path:          "/tmp/comment-project",
		DefaultBranch: "main",
	}
	if err := db.Create(project).Error; err != nil {
		t.Fatalf("create project failed: %v", err)
	}

	task := &TaskTable{
		ProjectID:  project.ID,
		Title:      "Seed task",
		Status:     "todo",
		OrderIndex: 1,
	}
	if err := db.Create(task).Error; err != nil {
		t.Fatalf("create task failed: %v", err)
	}

	comment := &TaskCommentTable{
		TaskID:  task.ID,
		Content: "Looks good to me",
	}

	if err := db.Create(comment).Error; err != nil {
		t.Fatalf("create comment failed: %v", err)
	}

	var loaded TaskCommentTable
	if err := db.Preload("Task").First(&loaded, "id = ?", comment.ID).Error; err != nil {
		t.Fatalf("load comment failed: %v", err)
	}

	if loaded.Task == nil || loaded.Task.ID != task.ID {
		t.Fatalf("expected associated task with id %s, got %#v", task.ID, loaded.Task)
	}
}
