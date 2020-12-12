package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

// CreateAccount inserts a new user into firestore [TODO]
func CreateAccount(user entity.User) error {
	return nil
}

// Login authenticates the user [TODO]
func Login(user entity.User) error {
	return nil
}

// GenerateJWT generates a JWT from user struct
func GenerateJWT(user entity.User) (jwtString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":      user.Login,
		"timestamp": time.Now().Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	jwtString, err = token.SignedString([]byte(secret))

	return
}

// ValidateJWT takes a JWT token string and validates it
func ValidateJWT(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return os.Getenv("JWT_SECRET"), nil
	})

	if !token.Valid {
		return err
	}

	return nil
}
