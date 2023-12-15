CREATE TABLE IF NOT EXISTS endpoints (
    id TEXT PRIMARY KEY,
    method TEXT NOT NULL,
    path TEXT NOT NULL,
    response TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT,
    namespace_id TEXT,
    FOREIGN KEY (namespace_id) REFERENCES namespaces(id)
);