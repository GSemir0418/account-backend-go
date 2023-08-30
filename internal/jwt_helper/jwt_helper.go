package jwt_helper

import (
	"crypto/rand"
	"io"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(user_id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user_id,
	})
	secret, err := generateHmacSecret()
	if err != nil {
		return "", err
	}
	return token.SignedString(secret)
}

func generateHmacSecret() ([]byte, error) {
	key := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
