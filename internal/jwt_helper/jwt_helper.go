package jwt_helper

import (
	"crypto/rand"
	"io"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func GenerateJWT(user_id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user_id,
	})
	secret, err := getHmacSecret()
	if err != nil {
		return "", err
	}
	return token.SignedString(secret)
}

func getHmacSecret() ([]byte, error) {
	keyPath := viper.GetString("jwt.hmac.keyPath")
	return os.ReadFile(keyPath)
}

func GenerateHmacSecret() ([]byte, error) {
	key := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func ParseJWT(jwtString string) (*jwt.Token, error) {
	key, err := getHmacSecret()
	if err != nil {
		return nil, err
	}
	return jwt.Parse(jwtString, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
}
