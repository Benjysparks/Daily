-- name: GetUserByID :one
SELECT id, full_name, email FROM users WHERE id = $1;
