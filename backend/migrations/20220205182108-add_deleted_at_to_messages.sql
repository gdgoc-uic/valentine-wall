
-- +migrate Up
ALTER TABLE messages ADD COLUMN deleted_at created_at TIMESTAMP;

-- +migrate Down
ALTER TABLE messages DROP COLUMN deleted_at;