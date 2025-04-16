package utils

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
)

func GetDecryptionKey(ctx context.Context, storage UserGetter, masterKey string) (pgtype.UUID, string, error) {
	userUUID, err := GetUserID(ctx)
	if err != nil {
		return pgtype.UUID{}, "", errors.Wrap(err, "error getting user id")
	}

	decryptedUserKey, err := GetUserKey(ctx, storage, userUUID, masterKey)
	if err != nil {
		return pgtype.UUID{}, "", errors.Wrap(err, "error getting decryption key")
	}

	return userUUID, decryptedUserKey, nil
}
