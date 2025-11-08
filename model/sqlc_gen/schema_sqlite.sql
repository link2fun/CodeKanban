-- Database schema for sqlite
-- Updated: 2025-11-08 04:30:00

CREATE TABLE "projects" (
  "id" text PRIMARY KEY,
  "created_at" datetime NOT NULL,
  "updated_at" datetime NOT NULL,
  "deleted_at" datetime,
  "name" text NOT NULL,
  "path" text NOT NULL,
  "description" text,
  "default_branch" text NOT NULL DEFAULT 'main',
  "worktree_base_path" text,
  "remote_url" text,
  "last_sync_at" datetime
);
CREATE UNIQUE INDEX "idx_projects_path" ON "projects"("path") WHERE deleted_at IS NULL;
CREATE INDEX "idx_projects_name" ON "projects"("name");
CREATE INDEX "idx_projects_deleted_at" ON "projects"("deleted_at");

CREATE TABLE "worktrees" (
  "id" text PRIMARY KEY,
  "created_at" datetime NOT NULL,
  "updated_at" datetime NOT NULL,
  "deleted_at" datetime,
  "project_id" text NOT NULL,
  "branch_name" text NOT NULL,
  "path" text NOT NULL,
  "is_main" numeric NOT NULL DEFAULT 0,
  "is_bare" numeric NOT NULL DEFAULT 0,
  "head_commit" text,
  "status_ahead" integer NOT NULL DEFAULT 0,
  "status_behind" integer NOT NULL DEFAULT 0,
  "status_modified" integer NOT NULL DEFAULT 0,
  "status_staged" integer NOT NULL DEFAULT 0,
  "status_untracked" integer NOT NULL DEFAULT 0,
  "status_conflicts" integer NOT NULL DEFAULT 0,
  "status_updated_at" datetime,
  FOREIGN KEY("project_id") REFERENCES "projects"("id") ON DELETE CASCADE
);
CREATE UNIQUE INDEX "idx_worktrees_path" ON "worktrees"("path") WHERE deleted_at IS NULL;
CREATE INDEX "idx_worktrees_project" ON "worktrees"("project_id");
CREATE INDEX "idx_worktrees_branch" ON "worktrees"("branch_name");
CREATE INDEX "idx_worktrees_deleted_at" ON "worktrees"("deleted_at");

CREATE TABLE "tasks" (
  "id" text PRIMARY KEY,
  "created_at" datetime NOT NULL,
  "updated_at" datetime NOT NULL,
  "deleted_at" datetime,
  "project_id" text NOT NULL,
  "worktree_id" text,
  "title" text NOT NULL,
  "description" text,
  "status" text NOT NULL,
  "priority" integer NOT NULL DEFAULT 0,
  "order_index" real NOT NULL,
  "tags" text,
  "due_date" datetime,
  "completed_at" datetime,
  FOREIGN KEY("project_id") REFERENCES "projects"("id") ON DELETE CASCADE,
  FOREIGN KEY("worktree_id") REFERENCES "worktrees"("id") ON DELETE SET NULL
);
CREATE INDEX "idx_tasks_project" ON "tasks"("project_id");
CREATE INDEX "idx_tasks_worktree" ON "tasks"("worktree_id");
CREATE INDEX "idx_tasks_status" ON "tasks"("status");
CREATE INDEX "idx_tasks_priority" ON "tasks"("priority");
CREATE INDEX "idx_tasks_deleted_at" ON "tasks"("deleted_at");

CREATE TABLE "task_comments" (
  "id" text PRIMARY KEY,
  "created_at" datetime NOT NULL,
  "updated_at" datetime NOT NULL,
  "deleted_at" datetime,
  "task_id" text NOT NULL,
  "content" text NOT NULL,
  FOREIGN KEY("task_id") REFERENCES "tasks"("id") ON DELETE CASCADE
);
CREATE INDEX "idx_task_comments_task" ON "task_comments"("task_id");
CREATE INDEX "idx_task_comments_deleted_at" ON "task_comments"("deleted_at");

CREATE TABLE "users" (
  "id" text,
  "created_at" datetime,
  "updated_at" datetime,
  "deleted_at" datetime,
  "nickname" text,
  "avatar" text,
  "brief" text,
  "username" text NOT NULL,
  "password" text NOT NULL,
  "salt" text NOT NULL,
  "disabled" numeric NOT NULL DEFAULT false,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "idx_users_username" ON "users"("username");
CREATE INDEX "idx_users_deleted_at" ON "users"("deleted_at");

CREATE TABLE "user_access_tokens" (
  "id" text,
  "created_at" datetime,
  "updated_at" datetime,
  "deleted_at" datetime,
  "user_id" text NOT NULL,
  "expired_at" datetime NOT NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "idx_access_tokens_user_id" ON "user_access_tokens"("user_id");
CREATE INDEX "idx_user_access_tokens_deleted_at" ON "user_access_tokens"("deleted_at");
