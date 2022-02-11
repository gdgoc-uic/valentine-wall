
-- +migrate Up
ALTER TABLE associated_ids ADD COLUMN sex TEXT NOT NULL DEFAULT 'unknown';

-- +migrate Down
ALTER TABLE associated_ids DROP COLUMN sex;