
-- +migrate Up
ALTER TABLE associated_ids ADD COLUMN terms_agreed BOOLEAN;

-- +migrate Down
ALTER TABLE associated_ids DROP COLUMN terms_agreed;