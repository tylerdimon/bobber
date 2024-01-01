CREATE TABLE IF NOT EXISTS requests (
    id TEXT PRIMARY KEY,
    method  TEXT NOT NULL,
    host  TEXT NOT NULL,
    path  TEXT NOT NULL,
    timestamp  TEXT NOT NULL,
    body  TEXT,
    headers  TEXT,
    namespace_id  TEXT,
    endpoint_id  TEXT,
    FOREIGN KEY (namespace_id) REFERENCES namespaces(id),
    FOREIGN KEY (endpoint_id) REFERENCES endpoints(id)
);