
-- +migrate Up
ALTER TABLE associated_users ALTER COLUMN user_id TYPE VARCHAR(255);

-- +migrate Down
ALTER TABLE associated_users ALTER COLUMN user_id TYPE TEXT;