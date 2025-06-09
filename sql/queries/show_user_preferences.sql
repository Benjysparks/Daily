-- name: ShowUserPreferencesByEmail :many
SELECT u.id, u.email, up.preferences
FROM users u
LEFT JOIN user_preferences up ON u.id = up.user_id;