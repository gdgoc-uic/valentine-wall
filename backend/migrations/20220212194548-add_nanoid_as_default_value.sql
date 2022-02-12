
-- +migrate Up
ALTER TABLE messages ALTER COLUMN id SET DEFAULT nanoid();
ALTER TABLE user_invitation_codes ALTER COLUMN id SET DEFAULT nanoid();
ALTER TABLE cheques ALTER COLUMN id SET DEFAULT nanoid();

-- +migrate Down
ALTER TABLE messages ALTER COLUMN id DROP DEFAULT;
ALTER TABLE user_invitation_codes ALTER COLUMN id DROP DEFAULT;
ALTER TABLE cheques ALTER COLUMN id DROP DEFAULT;