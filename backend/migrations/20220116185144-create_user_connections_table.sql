
-- +migrate Up
CREATE TABLE user_connections (
  user_id TEXT PRIMARY KEY,
  provider TEXT NOT NULL,
  -- refresh_token TEXT NOT NULL,
  -- access_token TEXT NOT NULL,
  -- expires_at DATETIME NOT NULL
  token TEXT NOT NULL,
  token_secret TEXT NOT NULL
);

-- +migrate Down
DROP TABLE user_connections;