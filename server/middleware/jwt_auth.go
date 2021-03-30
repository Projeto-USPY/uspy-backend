// package middleware contains useful middleware handlers
package middleware

import (
	"fmt"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

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
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if token == nil || !token.Valid {
		return nil, err
	}

	return token, nil
}

// JWT is used to ensure authorization with the JWT Access Cookie.
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := ValidateJWT(cookie)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("access_token", token)
	}
}
