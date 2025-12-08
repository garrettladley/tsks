CREATE TABLE IF NOT EXISTS tasks (
    id TEXT NOT NULL,
    version DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    archived_at DATETIME DEFAULT NULL,
    PRIMARY KEY (id, version)
);

CREATE INDEX IF NOT EXISTS idx_tasks_id_archived ON tasks(id, archived_at);

CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_tasks_id_version ON tasks(id, version DESC);
