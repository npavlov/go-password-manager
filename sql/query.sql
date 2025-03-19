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
WHERE id = $1 and user_id = $2;

-- name: CreateNoteEntry :one
INSERT INTO notes (user_id, encrypted_content)
VALUES ($1, $2)
    RETURNING *;

-- name: GetNotesByUserID :many
SELECT * FROM notes
WHERE user_id = $1;

-- name: GetNoteByID :one
SELECT * FROM notes WHERE id = $1;

-- name: DeleteNoteEntry :exec
DELETE FROM notes WHERE id = $1 and user_id = $2;

-- name: StoreCard :one
INSERT INTO cards (user_id, hashed_card_number, encrypted_card_number, encrypted_expiry_date, encrypted_cvv, cardholder_name)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: UpdateCard :one
UPDATE cards
SET encrypted_card_number = $1, encrypted_expiry_date = $2, encrypted_cvv = $3, cardholder_name = $4, hashed_card_number = $5
WHERE id = $6
    RETURNING *;

-- name: GetCardsByUserID :many
SELECT * FROM cards WHERE user_id = $1;

-- name: GetCardByID :one
SELECT * FROM cards WHERE id = $1;

-- name: DeleteCard :exec
DELETE FROM cards WHERE id = $1 and user_id = $2;

-- name: StoreBinaryEntry :one
INSERT INTO binary_entries (user_id, file_name, file_url, file_size)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: GetBinaryEntriesByUserID :many
SELECT * FROM binary_entries WHERE user_id = $1;

-- name: GetBinaryEntryByID :one
SELECT * FROM binary_entries WHERE id = $1;

-- name: DeleteBinaryEntry :exec
DELETE FROM binary_entries WHERE id = $1 and user_id = $2;

-- name: GetItemsByUserID :many
SELECT
    i.id,
    i.type,
    i.id_resource,
    i.created_at,
    COALESCE(p.updated_at, n.updated_at, c.updated_at, b.updated_at) AS updated_at
FROM items i
         LEFT JOIN passwords p ON i.type = 'password' AND i.id_resource = p.id
         LEFT JOIN notes n ON i.type = 'text' AND i.id_resource = n.id
         LEFT JOIN cards c ON i.type = 'card' AND i.id_resource = c.id
         LEFT JOIN binary_entries b ON i.type = 'binary' AND i.id_resource = b.id
WHERE i.user_id = $1
ORDER BY i.created_at DESC
    LIMIT $2 OFFSET $3;

-- name: GetTotalItemCountByUserID :one
SELECT COUNT(*) FROM items WHERE user_id = $1;