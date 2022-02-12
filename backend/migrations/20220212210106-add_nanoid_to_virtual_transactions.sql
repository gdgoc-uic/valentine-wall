
-- +migrate Up
ALTER TABLE virtual_transactions ALTER COLUMN id SET DEFAULT nanoid();

-- +migrate Down
ALTER TABLE virtual_transactions ALTER COLUMN id DROP DEFAULT;
