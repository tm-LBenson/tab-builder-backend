CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  uid         TEXT PRIMARY KEY,
  display_name TEXT
);

CREATE TABLE songs (
  id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  owner_uid   TEXT REFERENCES users(uid),
  title       TEXT NOT NULL,
  is_public   BOOLEAN DEFAULT FALSE,
  payload     JSONB NOT NULL,
  created_at  TIMESTAMPTZ DEFAULT NOW(),
  updated_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX songs_owner_idx ON songs(owner_uid);
