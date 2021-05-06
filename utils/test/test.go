package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/server"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/Projeto-USPY/uspy-backend/utils/test/emulator"
	"github.com/gin-gonic/gin"
)

// setupAccessToken fetches the jwt token used for private and restricted tests
func setupAccessToken(router *gin.Engine) (*http.Cookie, error) {
	// Login data
	jsonBody := map[string]interface{}{"login": "123456789", "pwd": "r4nd0mpass123!@#", "remember": true}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(jsonBody)

	// Execute login
	w := utils.MakeRequest(router, http.MethodPost, "/account/login", payloadBuf)

	if w.Code != http.StatusOK || len(w.Result().Cookies()) < 1 {
		return nil, errors.New("could not made login")
	}

	// Fetch returned cookie
	return w.Result().Cookies()[0], nil
}

// GetEnvironment will reinitialize the testing environment.
// It requires a suite because it is meant to be run with suites.
func MustGetEnvironment() (DB db.Env, router *gin.Engine, cookie *http.Cookie) {
	DB = emulator.MustGet()
	if err := emulator.Setup(DB); err != nil {
		panic(err)
	}

	// setup router
	var err error
	router, err = server.SetupRouter(DB)
	if err != nil {
		panic(err)
	}

	// get valid AccessToken
	if cookie, err = setupAccessToken(router); err != nil {
		panic(err)
	}

	return
}
