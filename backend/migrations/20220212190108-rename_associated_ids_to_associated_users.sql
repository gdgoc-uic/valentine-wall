
-- +migrate Up
ALTER TABLE associated_users RENAME TO associated_users;

-- +migrate Down
ALTER TABLE associated_users RENAME TO associated_users;