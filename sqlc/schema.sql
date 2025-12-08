-- Event sourced tasks table with composite key (id, version)
-- Each task mutation creates a new row with new version timestamp
-- Queries select latest version by MAX(version) per task ID
CREATE TABLE tasks (
    id TEXT NOT NULL,
    version TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    archived_at TEXT DEFAULT NULL,
    PRIMARY KEY (id, version)
);

-- Performance indexes for common query patterns
CREATE INDEX idx_tasks_id_archived ON tasks(id, archived_at);
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);
CREATE INDEX idx_tasks_id_version ON tasks(id, version DESC);
