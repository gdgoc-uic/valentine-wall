
-- +migrate Up
CREATE TABLE associated_ids (
  user_id TEXT PRIMARY KEY,
  associated_id TEXT NOT NULL UNIQUE
);

-- +migrate Down
DROP TABLE associated_ids;