CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    pass_hash VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    UNIQUE(username)
);