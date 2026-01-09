
-- name: CreateUser :one
INSERT INTO users (id, email, password_hash) VALUES (
    gen_random_uuid(),
    $1,
    $2
) RETURNING id, created_at, updated_at, email, password_hash;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;