CREATE TABLE IF NOT EXISTS endpoints (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL,
    response TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT,
    namespace_id INTEGER NOT NULL,
    FOREIGN KEY (namespace_id) REFERENCES namespaces(id)
);