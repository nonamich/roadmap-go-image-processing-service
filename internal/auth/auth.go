package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string
	Password string
}

func getSecret() string {
	secret := os.Getenv("JWT_SECRET")

	return secret
}

func JwtDecode(bearerToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (any, error) {
		return getSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func JwtEncode(user User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Username": user.Username,
	})

	tokenString, _ := token.SignedString(getSecret())

	return tokenString
}
