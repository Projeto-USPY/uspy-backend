package validation

import (
	"time"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
)

func validateVerificationToken(f1 validator.FieldLevel) bool {
	str := f1.Field().String()
	if token, err := utils.ValidateJWT(str, config.Env.JWTSecret); err != nil {
		return false
	} else {
		claims := token.Claims.(jwt.MapClaims)

		// assert correct operation
		if operation, ok := claims["type"].(string); !ok {
			return false
		} else if operation != "email_verification" {
			return false
		}

		// assert timestamp is less than an hour old
		if creationDate, ok := claims["timestamp"].(string); !ok {
			return false
		} else if t, err := time.Parse(time.RFC3339Nano, creationDate); err != nil {
			return false
		} else {
			return time.Since(t) < time.Hour
		}
	}
}

func validateRecoveryToken(f1 validator.FieldLevel) bool {
	str := f1.Field().String()
	if token, err := utils.ValidateJWT(str, config.Env.JWTSecret); err != nil {
		return false
	} else {
		claims := token.Claims.(jwt.MapClaims)

		// assert correct operation
		if operation, ok := claims["type"].(string); !ok {
			return false
		} else if operation != "password_reset" {
			return false
		}

		// assert user hash is sent and is sha256
		if hash, ok := claims["user"].(string); !ok {
			return false
		} else if len(hash) != 64 || !utils.IsHex(hash) {
			return false
		}

		// assert timestamp is less than an hour old
		if creationDate, ok := claims["timestamp"].(string); !ok {
			return false
		} else if t, err := time.Parse(time.RFC3339Nano, creationDate); err != nil {
			return false
		} else {
			return time.Since(t) < time.Hour
		}
	}
}
