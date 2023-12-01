CREATE TABLE IF NOT EXISTS endpoints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT,
    response TEXT,
    created_at TEXT,
    updated_at TEXT,
    namespace_id INTEGER NOT NULL,
    FOREIGN KEY (namespace_id) REFERENCES namespaces(id)
);