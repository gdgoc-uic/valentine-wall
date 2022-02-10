
-- +migrate Up
ALTER TABLE associated_ids ADD COLUMN last_active_at TIMESTAMP;

-- +migrate Down
ALTER TABLE associated_ids DROP COLUMN last_active_at;