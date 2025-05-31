CREATE TABLE IF NOT EXISTS user
(
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_email ON user(email);

CREATE TABLE IF NOT EXISTS task
(
    id UUID PRIMARY KEY,
    author UUID REFERENCES user(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    status      TEXT CHECK( status IN ('to-do','in-progress','done') ),
    deadline TIMESTAMP,
    created_at TIMESTAMP
)

