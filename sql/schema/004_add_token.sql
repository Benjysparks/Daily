-- +goose Up
ALTER TABLE users 
ADD COLUMN jwt_token TEXT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE users 
DROP COLUMN jwt_token;
