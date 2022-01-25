
-- +migrate Up
ALTER TABLE associated_ids ADD COLUMN department TEXT NOT NULL DEFAULT "unknown";

-- +migrate Down
ALTER TABLE associated_ids DROP COLUMN department;