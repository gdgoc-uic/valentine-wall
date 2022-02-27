
-- +migrate Up
CREATE TABLE archived_stats (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  total INTEGER NOT NULL DEFAULT 0,
  archived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE archived_stats;
