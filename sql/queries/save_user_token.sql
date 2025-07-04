-- name: SaveUserToken :exec
UPDATE users
SET jwt_token = $2
WHERE id = $1;
