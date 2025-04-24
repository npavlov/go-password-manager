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
INSERT INTO passwords (user_id, login, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPasswordEntriesByUserID :many
SELECT * FROM passwords
WHERE user_id = $1;

-- name: GetPasswordEntryByID :one
SELECT passwords.*
FROM passwords
WHERE passwords.id = $1 and passwords.user_id = $2;

-- name: UpdatePasswordEntry :one
UPDATE passwords
SET login = $1, password = $2
WHERE id = $3
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
SELECT notes.*
FROM notes
WHERE notes.id = $1 and notes.user_id = $2;

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
SELECT cards.*
FROM cards
WHERE cards.id = $1 and cards.user_id = $2;

-- name: DeleteCard :exec
DELETE FROM cards WHERE id = $1 and user_id = $2;

-- name: StoreBinaryEntry :one
INSERT INTO binary_entries (user_id, file_name, file_url, file_size)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: GetBinaryEntriesByUserID :many
SELECT * FROM binary_entries WHERE user_id = $1;

-- name: GetBinaryEntryByID :one
SELECT binary_entries.*
FROM binary_entries
WHERE binary_entries.id = $1 and binary_entries.user_id = $2;

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

-- name: AddMetaInfo :one
INSERT INTO metainfo (item_id, key, value)
VALUES ($1, $2, $3)
    ON CONFLICT (item_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP
                                      RETURNING *;

-- name: GetMetaInfoByItemID :many
SELECT key, value FROM metainfo WHERE item_id = $1;

-- name: DeleteMetaInfo :exec
DELETE FROM metainfo
WHERE item_id = $1 AND key = $2;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at
FROM refresh_tokens
WHERE token = $1;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;

-- name: DeleteUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: ExpireRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW();