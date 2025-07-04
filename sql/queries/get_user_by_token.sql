-- name: GetUserByToken :one
SELECT * FROM users
WHERE jwt_token = $1;