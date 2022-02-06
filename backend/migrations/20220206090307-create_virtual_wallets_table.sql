
-- +migrate Up
CREATE TABLE virtual_wallets (
  user_id TEXT NOT NULL UNIQUE,
  balance REAL NOT NULL DEFAULT 0.0
);

INSERT INTO virtual_wallets (user_id)
SELECT user_id FROM associated_ids;

-- +migrate Down
DROP TABLE virtual_wallets;