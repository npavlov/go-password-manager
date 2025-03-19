package utils

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/npavlov/go-password-manager/internal/server/storage"
	"github.com/pkg/errors"
)

const (
	TokenExpiration = time.Minute * 60
)

func GenerateJWT(userID string, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
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

func GetUserKey(ctx context.Context, storage *storage.DBStorage, userUUID pgtype.UUID, masterKey string) (string, error) {
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

func GetIdFromString(id string) pgtype.UUID {
	var uuid pgtype.UUID

	_ = uuid.Scan(id)

	return uuid
}
