package tables

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	if err := db.AutoMigrate(
		&ProjectTable{},
		&WorktreeTable{},
		&TaskTable{},
		&TaskCommentTable{},
	); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	return db
}
