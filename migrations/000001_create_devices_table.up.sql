CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    brand TEXT,
    state TEXT,
    creation_time TIMESTAMP DEFAULT NOW()
);
