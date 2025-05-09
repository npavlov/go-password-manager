// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const AddMetaInfo = `-- name: AddMetaInfo :one
INSERT INTO metainfo (item_id, key, value)
VALUES ($1, $2, $3)
    ON CONFLICT (item_id, key) DO UPDATE SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP
                                      RETURNING id, item_id, key, value, created_at, updated_at
`

type AddMetaInfoParams struct {
	ItemID pgtype.UUID `db:"item_id"`
	Key    string      `db:"key"`
	Value  string      `db:"value"`
}

func (q *Queries) AddMetaInfo(ctx context.Context, arg AddMetaInfoParams) (Metainfo, error) {
	row := q.db.QueryRow(ctx, AddMetaInfo, arg.ItemID, arg.Key, arg.Value)
	var i Metainfo
	err := row.Scan(
		&i.ID,
		&i.ItemID,
		&i.Key,
		&i.Value,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const CreateNoteEntry = `-- name: CreateNoteEntry :one
INSERT INTO notes (user_id, encrypted_content)
VALUES ($1, $2)
    RETURNING id, user_id, encrypted_content, created_at, updated_at
`

type CreateNoteEntryParams struct {
	UserID           pgtype.UUID `db:"user_id"`
	EncryptedContent string      `db:"encrypted_content"`
}

func (q *Queries) CreateNoteEntry(ctx context.Context, arg CreateNoteEntryParams) (Note, error) {
	row := q.db.QueryRow(ctx, CreateNoteEntry, arg.UserID, arg.EncryptedContent)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.EncryptedContent,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const CreatePasswordEntry = `-- name: CreatePasswordEntry :one
INSERT INTO passwords (user_id, login, password)
VALUES ($1, $2, $3)
RETURNING id, user_id, login, password, created_at, updated_at
`

type CreatePasswordEntryParams struct {
	UserID   pgtype.UUID `db:"user_id"`
	Login    string      `db:"login"`
	Password string      `db:"password"`
}

func (q *Queries) CreatePasswordEntry(ctx context.Context, arg CreatePasswordEntryParams) (Password, error) {
	row := q.db.QueryRow(ctx, CreatePasswordEntry, arg.UserID, arg.Login, arg.Password)
	var i Password
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Login,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const CreateRefreshToken = `-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
`

type CreateRefreshTokenParams struct {
	UserID    pgtype.UUID      `db:"user_id"`
	Token     string           `db:"token"`
	ExpiresAt pgtype.Timestamp `db:"expires_at"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) error {
	_, err := q.db.Exec(ctx, CreateRefreshToken, arg.UserID, arg.Token, arg.ExpiresAt)
	return err
}

const CreateUser = `-- name: CreateUser :one
INSERT INTO users (username, password, encryption_key, email)
VALUES ($1, $2, $3, $4)
    RETURNING id, username, email, password, encryption_key
`

type CreateUserParams struct {
	Username      string `db:"username"`
	Password      string `db:"password"`
	EncryptionKey string `db:"encryption_key"`
	Email         string `db:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, CreateUser,
		arg.Username,
		arg.Password,
		arg.EncryptionKey,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.EncryptionKey,
	)
	return i, err
}

const DeleteBinaryEntry = `-- name: DeleteBinaryEntry :exec
DELETE FROM binary_entries WHERE id = $1 and user_id = $2
`

type DeleteBinaryEntryParams struct {
	ID     pgtype.UUID `db:"id"`
	UserID pgtype.UUID `db:"user_id"`
}

func (q *Queries) DeleteBinaryEntry(ctx context.Context, arg DeleteBinaryEntryParams) error {
	_, err := q.db.Exec(ctx, DeleteBinaryEntry, arg.ID, arg.UserID)
	return err
}

const DeleteCard = `-- name: DeleteCard :exec
DELETE FROM cards WHERE id = $1 and user_id = $2
`

type DeleteCardParams struct {
	ID     pgtype.UUID `db:"id"`
	UserID pgtype.UUID `db:"user_id"`
}

func (q *Queries) DeleteCard(ctx context.Context, arg DeleteCardParams) error {
	_, err := q.db.Exec(ctx, DeleteCard, arg.ID, arg.UserID)
	return err
}

const DeleteMetaInfo = `-- name: DeleteMetaInfo :exec
DELETE FROM metainfo
WHERE item_id = $1 AND key = $2
`

type DeleteMetaInfoParams struct {
	ItemID pgtype.UUID `db:"item_id"`
	Key    string      `db:"key"`
}

func (q *Queries) DeleteMetaInfo(ctx context.Context, arg DeleteMetaInfoParams) error {
	_, err := q.db.Exec(ctx, DeleteMetaInfo, arg.ItemID, arg.Key)
	return err
}

const DeleteNoteEntry = `-- name: DeleteNoteEntry :exec
DELETE FROM notes WHERE id = $1 and user_id = $2
`

type DeleteNoteEntryParams struct {
	ID     pgtype.UUID `db:"id"`
	UserID pgtype.UUID `db:"user_id"`
}

func (q *Queries) DeleteNoteEntry(ctx context.Context, arg DeleteNoteEntryParams) error {
	_, err := q.db.Exec(ctx, DeleteNoteEntry, arg.ID, arg.UserID)
	return err
}

const DeletePasswordEntry = `-- name: DeletePasswordEntry :exec
DELETE FROM passwords
WHERE id = $1 and user_id = $2
`

type DeletePasswordEntryParams struct {
	ID     pgtype.UUID `db:"id"`
	UserID pgtype.UUID `db:"user_id"`
}

func (q *Queries) DeletePasswordEntry(ctx context.Context, arg DeletePasswordEntryParams) error {
	_, err := q.db.Exec(ctx, DeletePasswordEntry, arg.ID, arg.UserID)
	return err
}

const DeleteRefreshToken = `-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1
`

func (q *Queries) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, DeleteRefreshToken, token)
	return err
}

const DeleteUserRefreshTokens = `-- name: DeleteUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1
`

func (q *Queries) DeleteUserRefreshTokens(ctx context.Context, userID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, DeleteUserRefreshTokens, userID)
	return err
}

const ExpireRefreshTokens = `-- name: ExpireRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at < NOW()
`

func (q *Queries) ExpireRefreshTokens(ctx context.Context) error {
	_, err := q.db.Exec(ctx, ExpireRefreshTokens)
	return err
}

const GetBinaryEntriesByUserID = `-- name: GetBinaryEntriesByUserID :many
SELECT id, user_id, file_name, file_size, file_url, created_at, updated_at FROM binary_entries WHERE user_id = $1
`

func (q *Queries) GetBinaryEntriesByUserID(ctx context.Context, userID pgtype.UUID) ([]BinaryEntry, error) {
	rows, err := q.db.Query(ctx, GetBinaryEntriesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BinaryEntry
	for rows.Next() {
		var i BinaryEntry
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.FileName,
			&i.FileSize,
			&i.FileUrl,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetBinaryEntryByID = `-- name: GetBinaryEntryByID :one
SELECT binary_entries.id, binary_entries.user_id, binary_entries.file_name, binary_entries.file_size, binary_entries.file_url, binary_entries.created_at, binary_entries.updated_at
FROM binary_entries
WHERE binary_entries.id = $1 and binary_entries.user_id = $2
`

type GetBinaryEntryByIDParams struct {
	ID     pgtype.UUID `db:"id"`
	UserID pgtype.UUID `db:"user_id"`
}

func (q *Queries) GetBinaryEntryByID(ctx context.Context, arg GetBinaryEntryByIDParams) (BinaryEntry, error) {
	row := q.db.QueryRow(ctx, GetBinaryEntryByID, arg.ID, arg.UserID)
	var i BinaryEntry
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FileName,
		&i.FileSize,
		&i.FileUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetCardByID = `-- name: GetCardByID :one
SELECT cards.id, cards.user_id, cards.encrypted_card_number, cards.encrypted_expiry_date, cards.encrypted_cvv, cards.cardholder_name, cards.created_at, cards.updated_at, cards.hashed_card_number
FROM cards
WHERE cards.id = $1 and cards.user_id = $2
`

type GetCardByIDParams struct {
	ID     pgtype.UUID `db:"id"`
	UserID pgtype.UUID `db:"user_id"`
}

func (q *Queries) GetCardByID(ctx context.Context, arg GetCardByIDParams) (Card, error) {
	row := q.db.QueryRow(ctx, GetCardByID, arg.ID, arg.UserID)
	var i Card
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.EncryptedCardNumber,
		&i.EncryptedExpiryDate,
		&i.EncryptedCvv,
		&i.CardholderName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedCardNumber,
	)
	return i, err
}

const GetCardsByUserID = `-- name: GetCardsByUserID :many
SELECT id, user_id, encrypted_card_number, encrypted_expiry_date, encrypted_cvv, cardholder_name, created_at, updated_at, hashed_card_number FROM cards WHERE user_id = $1
`

func (q *Queries) GetCardsByUserID(ctx context.Context, userID pgtype.UUID) ([]Card, error) {
	rows, err := q.db.Query(ctx, GetCardsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Card
	for rows.Next() {
		var i Card
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.EncryptedCardNumber,
			&i.EncryptedExpiryDate,
			&i.EncryptedCvv,
			&i.CardholderName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.HashedCardNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetItemsByUserID = `-- name: GetItemsByUserID :many
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
    LIMIT $2 OFFSET $3
`

type GetItemsByUserIDParams struct {
	UserID pgtype.UUID `db:"user_id"`
	Limit  int32       `db:"limit"`
	Offset int32       `db:"offset"`
}

type GetItemsByUserIDRow struct {
	ID         pgtype.UUID      `db:"id"`
	Type       ItemType         `db:"type"`
	IDResource pgtype.UUID      `db:"id_resource"`
	CreatedAt  pgtype.Timestamp `db:"created_at"`
	UpdatedAt  pgtype.Timestamp `db:"updated_at"`
}

func (q *Queries) GetItemsByUserID(ctx context.Context, arg GetItemsByUserIDParams) ([]GetItemsByUserIDRow, error) {
	rows, err := q.db.Query(ctx, GetItemsByUserID, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetItemsByUserIDRow
	for rows.Next() {
		var i GetItemsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.Type,
			&i.IDResource,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetMetaInfoByItemID = `-- name: GetMetaInfoByItemID :many
SELECT key, value FROM metainfo WHERE item_id = $1
`

type GetMetaInfoByItemIDRow struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

func (q *Queries) GetMetaInfoByItemID(ctx context.Context, itemID pgtype.UUID) ([]GetMetaInfoByItemIDRow, error) {
	rows, err := q.db.Query(ctx, GetMetaInfoByItemID, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMetaInfoByItemIDRow
	for rows.Next() {
		var i GetMetaInfoByItemIDRow
		if err := rows.Scan(&i.Key, &i.Value); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetNoteByID = `-- name: GetNoteByID :one
SELECT notes.id, notes.user_id, notes.encrypted_content, notes.created_at, notes.updated_at
FROM notes
WHERE notes.id = $1 and notes.user_id = $2
`

type GetNoteByIDParams struct {
	ID     pgtype.UUID `db:"id"`
	UserID pgtype.UUID `db:"user_id"`
}

func (q *Queries) GetNoteByID(ctx context.Context, arg GetNoteByIDParams) (Note, error) {
	row := q.db.QueryRow(ctx, GetNoteByID, arg.ID, arg.UserID)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.EncryptedContent,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetNotesByUserID = `-- name: GetNotesByUserID :many
SELECT id, user_id, encrypted_content, created_at, updated_at FROM notes
WHERE user_id = $1
`

func (q *Queries) GetNotesByUserID(ctx context.Context, userID pgtype.UUID) ([]Note, error) {
	rows, err := q.db.Query(ctx, GetNotesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.EncryptedContent,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetPasswordEntriesByUserID = `-- name: GetPasswordEntriesByUserID :many
SELECT id, user_id, login, password, created_at, updated_at FROM passwords
WHERE user_id = $1
`

func (q *Queries) GetPasswordEntriesByUserID(ctx context.Context, userID pgtype.UUID) ([]Password, error) {
	rows, err := q.db.Query(ctx, GetPasswordEntriesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Password
	for rows.Next() {
		var i Password
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Login,
			&i.Password,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetPasswordEntryByID = `-- name: GetPasswordEntryByID :one
SELECT passwords.id, passwords.user_id, passwords.login, passwords.password, passwords.created_at, passwords.updated_at
FROM passwords
WHERE passwords.id = $1 and passwords.user_id = $2
`

type GetPasswordEntryByIDParams struct {
	ID     pgtype.UUID `db:"id"`
	UserID pgtype.UUID `db:"user_id"`
}

func (q *Queries) GetPasswordEntryByID(ctx context.Context, arg GetPasswordEntryByIDParams) (Password, error) {
	row := q.db.QueryRow(ctx, GetPasswordEntryByID, arg.ID, arg.UserID)
	var i Password
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Login,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetRefreshToken = `-- name: GetRefreshToken :one
SELECT id, user_id, token, expires_at
FROM refresh_tokens
WHERE token = $1
`

type GetRefreshTokenRow struct {
	ID        pgtype.UUID      `db:"id"`
	UserID    pgtype.UUID      `db:"user_id"`
	Token     string           `db:"token"`
	ExpiresAt pgtype.Timestamp `db:"expires_at"`
}

func (q *Queries) GetRefreshToken(ctx context.Context, token string) (GetRefreshTokenRow, error) {
	row := q.db.QueryRow(ctx, GetRefreshToken, token)
	var i GetRefreshTokenRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.ExpiresAt,
	)
	return i, err
}

const GetTotalItemCountByUserID = `-- name: GetTotalItemCountByUserID :one
SELECT COUNT(*) FROM items WHERE user_id = $1
`

func (q *Queries) GetTotalItemCountByUserID(ctx context.Context, userID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, GetTotalItemCountByUserID, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const GetUserByID = `-- name: GetUserByID :one
SELECT id, username, email, password, encryption_key FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, GetUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.EncryptionKey,
	)
	return i, err
}

const GetUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, password, encryption_key FROM users
WHERE username = $1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, GetUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.EncryptionKey,
	)
	return i, err
}

const StoreBinaryEntry = `-- name: StoreBinaryEntry :one
INSERT INTO binary_entries (user_id, file_name, file_url, file_size)
VALUES ($1, $2, $3, $4)
    RETURNING id, user_id, file_name, file_size, file_url, created_at, updated_at
`

type StoreBinaryEntryParams struct {
	UserID   pgtype.UUID `db:"user_id"`
	FileName string      `db:"file_name"`
	FileUrl  string      `db:"file_url"`
	FileSize int64       `db:"file_size"`
}

func (q *Queries) StoreBinaryEntry(ctx context.Context, arg StoreBinaryEntryParams) (BinaryEntry, error) {
	row := q.db.QueryRow(ctx, StoreBinaryEntry,
		arg.UserID,
		arg.FileName,
		arg.FileUrl,
		arg.FileSize,
	)
	var i BinaryEntry
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FileName,
		&i.FileSize,
		&i.FileUrl,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const StoreCard = `-- name: StoreCard :one
INSERT INTO cards (user_id, hashed_card_number, encrypted_card_number, encrypted_expiry_date, encrypted_cvv, cardholder_name)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, user_id, encrypted_card_number, encrypted_expiry_date, encrypted_cvv, cardholder_name, created_at, updated_at, hashed_card_number
`

type StoreCardParams struct {
	UserID              pgtype.UUID `db:"user_id"`
	HashedCardNumber    pgtype.Text `db:"hashed_card_number"`
	EncryptedCardNumber string      `db:"encrypted_card_number"`
	EncryptedExpiryDate string      `db:"encrypted_expiry_date"`
	EncryptedCvv        string      `db:"encrypted_cvv"`
	CardholderName      string      `db:"cardholder_name"`
}

func (q *Queries) StoreCard(ctx context.Context, arg StoreCardParams) (Card, error) {
	row := q.db.QueryRow(ctx, StoreCard,
		arg.UserID,
		arg.HashedCardNumber,
		arg.EncryptedCardNumber,
		arg.EncryptedExpiryDate,
		arg.EncryptedCvv,
		arg.CardholderName,
	)
	var i Card
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.EncryptedCardNumber,
		&i.EncryptedExpiryDate,
		&i.EncryptedCvv,
		&i.CardholderName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedCardNumber,
	)
	return i, err
}

const UpdateCard = `-- name: UpdateCard :one
UPDATE cards
SET encrypted_card_number = $1, encrypted_expiry_date = $2, encrypted_cvv = $3, cardholder_name = $4, hashed_card_number = $5
WHERE id = $6
    RETURNING id, user_id, encrypted_card_number, encrypted_expiry_date, encrypted_cvv, cardholder_name, created_at, updated_at, hashed_card_number
`

type UpdateCardParams struct {
	EncryptedCardNumber string      `db:"encrypted_card_number"`
	EncryptedExpiryDate string      `db:"encrypted_expiry_date"`
	EncryptedCvv        string      `db:"encrypted_cvv"`
	CardholderName      string      `db:"cardholder_name"`
	HashedCardNumber    pgtype.Text `db:"hashed_card_number"`
	ID                  pgtype.UUID `db:"id"`
}

func (q *Queries) UpdateCard(ctx context.Context, arg UpdateCardParams) (Card, error) {
	row := q.db.QueryRow(ctx, UpdateCard,
		arg.EncryptedCardNumber,
		arg.EncryptedExpiryDate,
		arg.EncryptedCvv,
		arg.CardholderName,
		arg.HashedCardNumber,
		arg.ID,
	)
	var i Card
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.EncryptedCardNumber,
		&i.EncryptedExpiryDate,
		&i.EncryptedCvv,
		&i.CardholderName,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.HashedCardNumber,
	)
	return i, err
}

const UpdatePasswordEntry = `-- name: UpdatePasswordEntry :one
UPDATE passwords
SET login = $1, password = $2
WHERE id = $3
    RETURNING id, user_id, login, password, created_at, updated_at
`

type UpdatePasswordEntryParams struct {
	Login    string      `db:"login"`
	Password string      `db:"password"`
	ID       pgtype.UUID `db:"id"`
}

func (q *Queries) UpdatePasswordEntry(ctx context.Context, arg UpdatePasswordEntryParams) (Password, error) {
	row := q.db.QueryRow(ctx, UpdatePasswordEntry, arg.Login, arg.Password, arg.ID)
	var i Password
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Login,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
