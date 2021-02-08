// package account contains functions that implement backend-db communication for every /account endpoint
package account

import (
	"github.com/tpreischadt/ProjetoJupiter/server/middleware"
	"testing"

	"github.com/joho/godotenv"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

func TestGenerateToken(t *testing.T) {
	godotenv.Load(".env")

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
	godotenv.Load(".env")

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
