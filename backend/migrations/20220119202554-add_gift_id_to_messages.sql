
-- +migrate Up
ALTER TABLE messages ADD COLUMN gift_id int;

-- +migrate Down
-- ALTER TABLE messages DROP COLUMN gift_id;