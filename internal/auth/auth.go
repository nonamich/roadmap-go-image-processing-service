package Auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	UserPackage "nonamich/image-processing-service/internal/user"
)

func getSecret() string {
	secret := os.Getenv("JWT_SECRET")

	return secret
}

func JwtDecode(bearerToken string) (*jwt.Token, error) {
	secret := getSecret()
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func JwtEncode(user UserPackage.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Username": user.Username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour).Unix(),
	})

	secret := getSecret()
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
