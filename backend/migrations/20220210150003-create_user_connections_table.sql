
-- +migrate Up
CREATE TABLE user_connections_new (
  user_id TEXT,
  provider TEXT NOT NULL,
  token TEXT NOT NULL,
  token_secret TEXT NOT NULL
);

INSERT INTO user_connections_new SELECT * FROM user_connections;
DROP TABLE user_connections;

-- +migrate Down
CREATE TABLE user_connections (
  user_id TEXT PRIMARY KEY,
  provider TEXT NOT NULL,
  -- refresh_token TEXT NOT NULL,
  -- access_token TEXT NOT NULL,
  -- expires_at DATETIME NOT NULL
  token TEXT NOT NULL,
  token_secret TEXT NOT NULL
);

INSERT INTO user_connections SELECT * FROM user_connections_new;
DROP TABLE user_connections_new;