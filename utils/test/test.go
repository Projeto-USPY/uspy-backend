package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

// SetupAccessToken fetches the jwt token used for private and restricted tests
func SetupAccessToken(router *gin.Engine) (*http.Cookie, error) {
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
