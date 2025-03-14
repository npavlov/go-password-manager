package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
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
