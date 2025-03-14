-- name: CreateUser :one
INSERT INTO users (username, password, encryption_key, email)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: CreatePasswordEntry :one
INSERT INTO passwords (user_id, name, login, password)
VALUES ($1, $2, $3, $4)
    ON CONFLICT (user_id, name) DO NOTHING
RETURNING *;

-- name: GetPasswordEntriesByUserID :many
SELECT * FROM passwords
WHERE user_id = $1;

-- name: GetPasswordEntryByID :one
SELECT * FROM passwords
WHERE id = $1;

-- name: UpdatePasswordEntry :one
UPDATE passwords
SET name = $1, login = $2, password = $3
WHERE id = $4
    RETURNING *;

-- name: DeletePasswordEntry :exec
DELETE FROM passwords
WHERE id = $1;