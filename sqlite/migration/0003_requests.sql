CREATE TABLE IF NOT EXISTS requests (
    id TEXT PRIMARY KEY,
    method TEXT,
    url TEXT,
    host TEXT,
    path TEXT,
    timestamp TEXT,
    body TEXT,
    headers TEXT,
    namespace_id TEXT,
    endpoint_id TEXT,
    FOREIGN KEY (namespace_id) REFERENCES namespaces(id),
    FOREIGN KEY (endpoint_id) REFERENCES endpoints(id)
);