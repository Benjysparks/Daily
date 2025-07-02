-- +goose Up
CREATE TABLE user_preferences (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  preferences JSONB NOT NULL
);

-- +goose Down
DROP TABLE user_preferences;
