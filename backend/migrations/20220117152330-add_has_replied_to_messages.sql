
-- +migrate Up
ALTER TABLE messages ADD COLUMN has_replied BOOLEAN NOT NULL DEFAULT false;

-- +migrate Down
ALTER TABLE messages DROP COLUMN has_replied;