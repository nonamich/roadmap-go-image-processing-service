package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"nonamich/image-processing-service/internal/user"
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

func JwtEncode(user user.User) (string, error) {
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

func CurrentUser(r *http.Request) (user.User, error) {
	authorization := r.Header.Get("Authorization")

	if authorization == "" {
		return user.User{}, errors.New("No authorization")
	}

	tokenString := strings.TrimPrefix(authorization, "Bearer ")

	if tokenString == authorization {
		return user.User{}, errors.New("Invalid authorization format")
	}

	token, err := JwtDecode(tokenString)

	if err != nil {
		return user.User{}, errors.New("Unauthorized.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return user.User{}, errors.New("Invalid token claims")
	}

	username, ok := claims["Username"].(string)

	if !ok {
		return user.User{}, fmt.Errorf("Toke must have username(%s) in claims.", username)
	}

	foundUser, err := user.FindUserByUsername(username)

	if err != nil {
		return user.User{}, err
	}

	if foundUser.ID == 0 {
		return user.User{}, errors.New("User not found.")
	}

	return foundUser, nil
}
