-- name: SavePreferences :exec
INSERT INTO user_preferences (user_id, preferences)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
SET preferences = EXCLUDED.preferences;
