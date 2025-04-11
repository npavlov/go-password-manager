package utils

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"

	"github.com/npavlov/go-password-manager/internal/server/db"
)

func GenerateJWT(userID, jwtSecret string, expiration int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expiration,
	})
	result, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.Wrap(err, "Error generating JWT")
	}

	return result, nil
}

func ValidateJWT(tokenString string, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.Wrap(err, "Error validating JWT")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.Wrap(err, "Error validating JWT")
	}

	return userID, nil
}

func GetUserId(ctx context.Context) (pgtype.UUID, error) {
	var uuid pgtype.UUID

	userId, ok := ctx.Value("user_id").(string)
	if !ok || userId == "" {
		return uuid, errors.New("Error getting user id")
	}

	if err := uuid.Scan(userId); err != nil {
		return uuid, errors.Wrap(err, "Error getting user id")
	}

	return uuid, nil
}

type UserGetter interface {
	GetUserById(ctx context.Context, id pgtype.UUID) (*db.User, error)
}

func GetUserKey(ctx context.Context, storage UserGetter, userUUID pgtype.UUID, masterKey string) (string, error) {
	user, err := storage.GetUserById(ctx, userUUID)
	if err != nil {
		return "", errors.Wrap(err, "Error getting user id")
	}

	decryptedUserKey, err := Decrypt(user.EncryptionKey, masterKey)
	if err != nil {
		return "", errors.Wrap(err, "Error decrypting user id")
	}

	return decryptedUserKey, nil
}
