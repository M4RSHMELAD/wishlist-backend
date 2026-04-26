-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS push_notifications BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_notifications BOOLEAN DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS show_fulfilled_wishes BOOLEAN DEFAULT true;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS push_notifications;
ALTER TABLE users DROP COLUMN IF EXISTS email_notifications;
ALTER TABLE users DROP COLUMN IF EXISTS show_fulfilled_wishes;
