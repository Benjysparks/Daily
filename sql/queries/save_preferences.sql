-- name: SavePreferences :exec
INSERT INTO user_preferences (user_id, preferences, preference_variables)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE
SET preferences = EXCLUDED.preferences,
    preference_variables = EXCLUDED.preference_variables;
