package utils

import "github.com/jackc/pgx/v5/pgtype"

func GetIdFromString(id string) pgtype.UUID {
	var uuid pgtype.UUID

	_ = uuid.Scan(id)

	return uuid
}
