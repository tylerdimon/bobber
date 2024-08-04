-- requests will come in and generate extracted values
-- which are used to send back specific responses
-- and dynamically update values in the response
CREATE TABLE IF NOT EXISTS extracted (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL,
    value TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT,
    extractor_id TEXT,
    request_id TEXT,
    FOREIGN KEY (extractor_id) REFERENCES extractors(id),
    FOREIGN KEY (request_id) REFERENCES requests(id)
);