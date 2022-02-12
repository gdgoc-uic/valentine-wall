
-- +migrate Up
ALTER TABLE virtual_wallets ADD COLUMN id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY;

-- +migrate Down
ALTER TABLE virtual_wallets DROP CONSTRAINT virtual_wallets_pkey;
ALTER TABLE virtual_wallets DROP COLUMN id;