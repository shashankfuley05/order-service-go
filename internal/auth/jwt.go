package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId string) (string, error) {

	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(20 * time.Minute).Unix(),
	}

	secret := []byte("AbraKadabra")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ParseToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {

		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte("AbraKadabra"), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Token Invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("Invalid Claim")
	}

	user, ok := claims["userId"].(string)

	if !ok {
		return "", errors.New("User not found")
	}
	return user, nil
}
