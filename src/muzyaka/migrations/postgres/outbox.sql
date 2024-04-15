CREATE TABLE IF NOT EXISTS outbox(
    id SERIAL PRIMARY KEY,
    event_id TEXT NOT NULL,
    track_id INTEGER NOT NULL,
    source VARCHAR(254),
    name VARCHAR(100),
    genre INT,
    type VARCHAR(100) NOT NULL,
    sent BOOLEAN DEFAULT FALSE NOT NULL
);
