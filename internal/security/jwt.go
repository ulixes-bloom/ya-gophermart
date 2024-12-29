package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

func BuildJWTString(userID int64, secretKey string, tokenLifetime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return tokenString, nil
}

func GetUserID(tokenString, secretKey string) (int64, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return -1, fmt.Errorf("security.jwt.getUserID: %w", err)
	}

	if !token.Valid {
		return -1, fmt.Errorf("security.jwt.getUserID: token not valid")
	}

	return claims.UserID, nil
}
