package utils_test

import (
	"testing"
	"time"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/utils"
)

func BenchmarkHashPasswords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash, _ := utils.Bcrypt("SenhaU3l34178!Fodida18723@#!")
		b.Log(hash)
	}
}

func TestGenerateToken(t *testing.T) {
	jwt, err := utils.GenerateJWT(map[string]interface{}{
		"user":      "login",
		"timestamp": time.Now().Unix(),
	}, config.Env.JWTSecret)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(jwt)
}

func TestValidateToken(t *testing.T) {
	jwt, err := utils.GenerateJWT(map[string]interface{}{
		"user":      "login",
		"timestamp": time.Now().Unix(),
	}, config.Env.JWTSecret)

	if err != nil {
		t.Fatal(err)
	}

	_, err = utils.ValidateJWT(jwt, config.Env.JWTSecret)

	if err != nil {
		t.Fatal(err)
	}
}
