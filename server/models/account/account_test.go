// package models
package account

import (
	"testing"

	"github.com/Projeto-USPY/uspy-backend/server/middleware"
)

func TestGenerateToken(t *testing.T) {
	jwt, err := middleware.GenerateJWT("login")

	if err != nil {
		t.Fatal(err)
	}

	t.Log(jwt)
}

func TestValidateToken(t *testing.T) {
	jwt, err := middleware.GenerateJWT("login")

	if err != nil {
		t.Fatal(err)
	}

	_, err = middleware.ValidateJWT(jwt)

	if err != nil {
		t.Fatal(err)
	}
}
