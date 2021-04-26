// package account contains functions that implement backend-db communication for every /account endpoint
package account

import (
	"github.com/Projeto-USPY/uspy-backend/server/middleware"
	"testing"

	"github.com/Projeto-USPY/uspy-backend/entity"
)

func TestGenerateToken(t *testing.T) {
	user := entity.User{}
	user.Login = "login"
	user.Password = "pass"

	jwt, err := middleware.GenerateJWT(user)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(jwt)
}

func TestValidateToken(t *testing.T) {
	user := entity.User{}
	user.Login = "login"
	user.Password = "pass"

	jwt, err := middleware.GenerateJWT(user)

	if err != nil {
		t.Fatal(err)
	}

	_, err = middleware.ValidateJWT(jwt)

	if err != nil {
		t.Fatal(err)
	}
}
