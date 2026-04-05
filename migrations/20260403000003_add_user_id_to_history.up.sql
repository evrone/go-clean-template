DROP TABLE IF EXISTS history;

CREATE TABLE IF NOT EXISTS history (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source VARCHAR(255),
    destination VARCHAR(255),
    original VARCHAR(255),
    translation VARCHAR(255)
);

CREATE INDEX idx_history_user_id ON history(user_id);
