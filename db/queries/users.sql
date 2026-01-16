
-- name: CreateUser :one
INSERT INTO users (id, email, password_hash) VALUES (
    gen_random_uuid(),
    $1,
    $2
) RETURNING id, created_at, updated_at, email, password_hash;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;  

-- name: UpdateUserPassword :one
UPDATE users SET password_hash = $2, updated_at = NOW() WHERE email = $1 RETURNING *;

-- name: SetEmailAndPassword :exec
UPDATE users
SET email = $2, password_hash = $3, updated_at = NOW()
WHERE id = $1;