// package middleware contains useful middleware handlers
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// GenerateJWT generates a JWT from user struct
func GenerateJWT(user *models.User) (jwtString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":      user.ID,
		"timestamp": time.Now().Unix(),
	})

	jwtString, err = token.SignedString([]byte(config.Env.JWTSecret))
	return
}

// ValidateJWT takes a JWT token string and validates it
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Env.JWTSecret), nil
	})

	if token == nil || !token.Valid {
		return nil, err
	}

	return token, nil
}

// JWT is used to ensure authorization with the JWT Access Cookie.
func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := ValidateJWT(cookie)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("access_token", token)
	}
}
