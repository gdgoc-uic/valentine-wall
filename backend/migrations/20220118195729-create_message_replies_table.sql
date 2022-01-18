
-- +migrate Up
CREATE TABLE message_replies (
  message_id TEXT NOT NULL REFERENCES messages(id),
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE message_replies;