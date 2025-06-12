-- +goose Up
ALTER TABLE users 
ADD COLUMN user_hours INTEGER NOT NULL DEFAULT 06,
ADD COLUMN user_minutes INTEGER NOT NULL DEFAULT 00;

-- +goose Down
ALTER TABLE users 
DROP COLUMN user_minutes,
DROP COLUMN user_hours;
