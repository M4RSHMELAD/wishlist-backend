-- Применение миграции для настроек пользователя
-- Выполни эту команду: docker exec -i wishlist_postgres psql -U postgres -d wishlist -f /tmp/migration.sql

-- Добавляем поля настроек в таблицу users
ALTER TABLE users ADD COLUMN IF NOT EXISTS push_notifications BOOLEAN DEFAULT true;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email_notifications BOOLEAN DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS show_fulfilled_wishes BOOLEAN DEFAULT true;

-- Устанавливаем значения по умолчанию для существующих пользователей
UPDATE users SET push_notifications = true WHERE push_notifications IS NULL;
UPDATE users SET email_notifications = false WHERE email_notifications IS NULL;
UPDATE users SET show_fulfilled_wishes = true WHERE show_fulfilled_wishes IS NULL;

-- Проверка что поля добавлены
SELECT column_name, data_type, column_default 
FROM information_schema.columns 
WHERE table_name = 'users' 
AND column_name IN ('push_notifications', 'email_notifications', 'show_fulfilled_wishes');
