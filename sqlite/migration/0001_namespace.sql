CREATE TABLE IF NOT EXISTS namespaces (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT DEFAULT ''
);