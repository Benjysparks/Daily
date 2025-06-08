-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, pword, full_name)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
RETURNING *;