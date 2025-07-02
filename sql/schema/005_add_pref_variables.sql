-- +goose Up
ALTER TABLE user_preferences
ADD COLUMN preference_variables JSONB NOT NULL;


-- +goose Down
ALTER TABLE user_preferences
DROP COLUMN preference_variables;
