
-- +migrate Up
CREATE TABLE message_gifts (
  message_id TEXT REFERENCES messages(id),
  gift_id INT NOT NULL
);

-- +migrate Down
DROP TABLE message_gifts;