
-- +migrate Up
ALTER TABLE associated_ids ADD COLUMN terms_agreed int;

-- +migrate Down
ALTER TABLE associated_ids DROP COLUMN terms_agreed;