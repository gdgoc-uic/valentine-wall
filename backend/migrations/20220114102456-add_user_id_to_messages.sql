
-- +migrate Up
ALTER TABLE messages ADD COLUMN submitter_user_id TEXT NOT NULL;

-- +migrate Down
ALTER TABLE messages DROP COLUMN submitter_user_id;