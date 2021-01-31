CREATE TABLE IF NOT EXISTS history(
    id serial PRIMARY KEY,
    original VARCHAR(255),
    translation VARCHAR(255)
);