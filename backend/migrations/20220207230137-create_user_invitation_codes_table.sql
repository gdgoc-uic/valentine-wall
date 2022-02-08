
-- +migrate Up
CREATE TABLE user_invitation_codes (
  id TEXT PRIMARY KEY NOT NULL,
  user_id TEXT NOT NULL,
  max_users INTEGER DEFAULT 1,
  user_count INTEGER DEFAULT 0,
  created_at TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE user_invitation_codes;