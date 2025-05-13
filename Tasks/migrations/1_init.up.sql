CREATE TABLE IF NOT EXISTS tasks(
    id TEXT PRIMARY KEY,
    title TEXT,
    description TEXT,
    author_id TEXT,
    status TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);