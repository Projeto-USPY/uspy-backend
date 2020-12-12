package auth

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/tpreischadt/ProjetoJupiter/entity"
)

func TestGenerateToken(t *testing.T) {
	godotenv.Load(".env")

	user := entity.User{}
	user.Login = "login"
	user.Password = "pass"

	jwt, err := GenerateJWT(user)

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

	jwt, err := GenerateJWT(user)

	if err != nil {
		t.Fatal(err)
	}

	err = ValidateJWT(jwt)

	if err != nil {
		t.Fatal(err)
	}
}
