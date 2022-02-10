
-- +migrate Up
CREATE TABLE cheques (
  id TEXT PRIMARY KEY NOT NULL,
  user_id TEXT NOT NULL,
  amount REAL NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE cheques;