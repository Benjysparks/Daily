-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, pword, full_name, user_hours, user_minutes, jwt_token)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;