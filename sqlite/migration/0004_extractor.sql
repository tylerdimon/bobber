-- endpoints have extractors
-- when request matches endpoints
-- get extractors and use the path to
-- traverse JSON for each one
-- and get extracted values (extracted table)
CREATE TABLE IF NOT EXISTS extractors (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT,
    endpoint_id TEXT,
    FOREIGN KEY (endpoint_id) REFERENCES endpoints(id)
);