CREATE TABLE IF NOT EXISTS requests (
    id TEXT PRIMARY KEY,
    method TEXT,
    url TEXT,
    host TEXT,
    path TEXT,
    timestamp TEXT,
    body TEXT,
    headers TEXT
);