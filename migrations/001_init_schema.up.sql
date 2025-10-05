CREATE TABLE IF NOT EXISTS calendar (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    date DATE NOT NULL,
    event TEXT
);