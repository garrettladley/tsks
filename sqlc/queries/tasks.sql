-- name: GetTask :one
SELECT * FROM tasks t1
WHERE t1.id = ?
  AND t1.archived_at IS NULL
  AND t1.version = (SELECT MAX(t2.version) FROM tasks t2 WHERE t2.id = t1.id)
LIMIT 1;

-- name: ListTasks :many
SELECT t1.* FROM tasks t1
INNER JOIN (
    SELECT id, MAX(version) as max_version
    FROM tasks
    GROUP BY id
) t2 ON t1.id = t2.id AND t1.version = t2.max_version
WHERE t1.archived_at IS NULL
ORDER BY t1.created_at DESC, t1.version DESC;

-- name: CreateTask :one
INSERT INTO tasks (
  id, title, description, status
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: CreateTasks :many
INSERT INTO tasks (
  id, title, description, status
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateTask :one
INSERT INTO tasks (
  id, title, description, status
) VALUES (
  ?, ?, ?, ?
)
RETURNING *;


-- name: DeleteTask :exec
UPDATE tasks
SET archived_at = CURRENT_TIMESTAMP
WHERE (id, version) IN (
  SELECT t1.id, MAX(t1.version)
  FROM tasks t1
  WHERE t1.id = ?
  GROUP BY t1.id
);

-- name: TruncateTasks :exec
DELETE FROM tasks;

-- name: ListTasksPageForward :many
-- Fetch older tasks (scrolling down). Pass NULL cursor values for initial load.
SELECT t1.* FROM tasks t1
INNER JOIN (
    SELECT id, MAX(version) as max_version
    FROM tasks
    GROUP BY id
) t2 ON t1.id = t2.id AND t1.version = t2.max_version
WHERE t1.archived_at IS NULL
  AND (
    CAST(sqlc.narg(cursor_created_at) AS TEXT) IS NULL
    OR t1.created_at < CAST(sqlc.narg(cursor_created_at) AS TEXT)
    OR (t1.created_at = CAST(sqlc.narg(cursor_created_at) AS TEXT) AND t1.id < CAST(sqlc.narg(cursor_id) AS TEXT))
  )
ORDER BY t1.created_at DESC, t1.id DESC
LIMIT sqlc.arg(page_limit);

-- name: ListTasksPageBackward :many
-- Fetch newer tasks (scrolling up).
SELECT t1.* FROM tasks t1
INNER JOIN (
    SELECT id, MAX(version) as max_version
    FROM tasks
    GROUP BY id
) t2 ON t1.id = t2.id AND t1.version = t2.max_version
WHERE t1.archived_at IS NULL
  AND (
    t1.created_at > CAST(sqlc.arg(cursor_created_at) AS TEXT)
    OR (t1.created_at = CAST(sqlc.arg(cursor_created_at) AS TEXT) AND t1.id > CAST(sqlc.arg(cursor_id) AS TEXT))
  )
ORDER BY t1.created_at ASC, t1.id ASC
LIMIT sqlc.arg(page_limit);
