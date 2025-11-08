-- name: ProjectCreate :one
INSERT INTO projects (
  id,
  created_at,
  updated_at,
  name,
  path,
  description,
  default_branch,
  worktree_base_path,
  remote_url,
  last_sync_at
) VALUES (
  @id,
  @created_at,
  @updated_at,
  @name,
  @path,
  @description,
  @default_branch,
  @worktree_base_path,
  @remote_url,
  @last_sync_at
) RETURNING
  id,
  created_at,
  updated_at,
  deleted_at,
  name,
  path,
  description,
  default_branch,
  worktree_base_path,
  remote_url,
  last_sync_at;

-- name: ProjectGetByID :one
SELECT
  id,
  created_at,
  updated_at,
  deleted_at,
  name,
  path,
  description,
  default_branch,
  worktree_base_path,
  remote_url,
  last_sync_at
FROM projects
WHERE id = @id
  AND deleted_at IS NULL
LIMIT 1;

-- name: ProjectList :many
SELECT
  id,
  created_at,
  updated_at,
  deleted_at,
  name,
  path,
  description,
  default_branch,
  worktree_base_path,
  remote_url,
  last_sync_at
FROM projects
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ProjectSoftDelete :execrows
UPDATE projects
SET
  deleted_at = @deleted_at,
  updated_at = @updated_at
WHERE id = @id
  AND deleted_at IS NULL;
