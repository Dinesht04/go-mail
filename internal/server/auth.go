package server

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	userId string
	jwt.RegisteredClaims
}

func CreateJWT(id string) (string, error) {

	claims := MyCustomClaims{
		userId: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("my_secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte("my_secret"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return false, err
	}

	if token.Valid {
		return true, nil
	} else {
		return false, fmt.Errorf("Token is invalid")
	}

}
