-- Add composite index for cursor-based pagination
-- Supports queries with (created_at, id) ordering and comparison
CREATE INDEX IF NOT EXISTS idx_tasks_pagination
    ON tasks(created_at DESC, id DESC)
    WHERE archived_at IS NULL;
