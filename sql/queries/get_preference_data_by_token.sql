-- name: GetPreferencesByToken :one
SELECT up.*
FROM users u
JOIN user_preferences up ON u.id = up.user_id
WHERE u.jwt_token = $1;
