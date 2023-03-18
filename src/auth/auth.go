package auth

import (
	"errors"
	"reflect"
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

	if claims["authorized"] != nil && reflect.TypeOf(claims["authorized"]).Name() != "bool" {
		return claims, errors.New("invalid JWT authorization")
	}

	if claims["authorized"].(bool) != true {
		return claims, errors.New("not authorized")
	}

	if claims["username"] != nil && reflect.TypeOf(claims["username"]).Name() != "string" {
		return claims, errors.New("invalid JWT username")
	}

	if claims["playerId"] != nil && reflect.TypeOf(claims["playerId"]).Name() != "string" {
		return claims, errors.New("invalid JWT playerId")
	}

	return claims, nil
}
