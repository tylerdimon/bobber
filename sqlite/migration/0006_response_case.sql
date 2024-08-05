-- once extracted values have all been calculated
-- go through response cases in order of priority
-- and if one match send back that specific response
-- TODO after response has been decided populate any response template {{json.path}} variables
CREATE TABLE IF NOT EXISTS response_case (
    id TEXT PRIMARY KEY,
    match_value TEXT,
    response TEXT,
    priority INT,
    created_at TEXT NOT NULL,
    updated_at TEXT,
    extractor_id TEXT,
    FOREIGN KEY (extractor_id) REFERENCES extractors(id)
);