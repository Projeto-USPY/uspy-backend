package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

// GenerateJWT generates a JWT from map
func GenerateJWT(data map[string]interface{}, secret string) (jwtString string, err error) {
	claims := make(jwt.MapClaims)
	for k, v := range data {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err = token.SignedString([]byte(secret))
	return
}

// ValidateJWT takes a JWT token string and validates it
func ValidateJWT(tokenString, secret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if token == nil || !token.Valid {
		return nil, err
	}

	return token, nil
}
