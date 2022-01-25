
-- +migrate Up
INSERT INTO message_gifts SELECT id message_id, gift_id FROM messages WHERE gift_id IS NOT NULL;
ALTER TABLE messages DROP COLUMN gift_id;

-- +migrate Down
ALTER TABLE messages ADD COLUMN gift_id INT;
UPDATE messages SET gift_id = (SELECT gift_id FROM message_gifts WHERE message_gifts.message_id = messages.id);
DELETE FROM message_gifts;