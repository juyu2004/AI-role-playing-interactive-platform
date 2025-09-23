package migrate

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// schemaSQL is embedded copy of database schema for automatic initialization.
const schemaSQL = `
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- PostgreSQL schema (initial draft)

CREATE TABLE IF NOT EXISTS roles (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    avatar_url TEXT,
    description TEXT,
    prompt TEXT,
    image_url TEXT
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    avatar_url TEXT,
    bio TEXT
);

CREATE TABLE IF NOT EXISTS conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id TEXT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender TEXT NOT NULL CHECK (sender IN ('user','role')),
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Friends (undirected as two directed edges)
CREATE TABLE IF NOT EXISTS friends (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY(user_id, friend_user_id)
);

-- Useful indexes
CREATE INDEX IF NOT EXISTS idx_conversations_user ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_conv ON messages(conversation_id, created_at DESC);

-- Seed default roles
INSERT INTO roles(id, name, category, avatar_url, description, prompt, image_url) VALUES
('sherlock','Sherlock Holmes','literature','', 'Brilliant detective with keen observation and deduction.', 'You are Sherlock Holmes. Respond concisely with sharp deductions.',''),
('mulan','Hua Mulan','history','', 'Heroine known for courage and loyalty.', 'You are Hua Mulan. Speak with bravery and humility.','')
ON CONFLICT (id) DO NOTHING;
`

// Apply runs the schema SQL to ensure required tables exist.
func Apply(db *sql.DB) error {
	_, err := db.Exec(schemaSQL)
	return err
}
