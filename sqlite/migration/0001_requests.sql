CREATE TABLE IF NOT EXISTS requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    method TEXT,
    url TEXT,
    host TEXT,
    path TEXT,
    timestamp TEXT,
    body TEXT,
    headers TEXT
);