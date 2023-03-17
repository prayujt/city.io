package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var SecretKey []byte

func GenerateJWT(username string, playerId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour)
	claims["authorized"] = true
	claims["username"] = username
	claims["playerId"] = playerId
	tokenString, err := token.SignedString(SecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(_token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(_token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(_token), nil
	})

	if (err != nil && err.Error() != "signature is invalid") || len(claims) == 0 {
		return claims, errors.New("internal error")
	}
	return claims, nil
}
