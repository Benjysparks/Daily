-- name: ShowAllUserPreferences :many
SELECT u.id, u.email, u.user_hours, u.user_minutes, up.preferences
FROM users u
LEFT JOIN user_preferences up ON u.id = up.user_id;