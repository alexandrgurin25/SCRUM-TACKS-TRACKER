CREATE TABLE IF NOT EXISTS projects
(
    id          VARCHAR PRIMARY KEY,
    author_id   VARCHAR,
    title       VARCHAR(256),
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);