package account_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Projeto-USPY/uspy-backend/db"
	emulator "github.com/Projeto-USPY/uspy-backend/emulator"
	"github.com/Projeto-USPY/uspy-backend/server"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AccountSuite struct {
	suite.Suite
	DB          db.Env
	router      *gin.Engine
	accessToken *http.Cookie
}

// SetupSuite runs before suite (fetches the emulator)
func (s *AccountSuite) SetupSuite() {
	s.DB = emulator.MustGet()

	// setup router
	var err error
	s.router, err = server.SetupRouter(s.DB)
	s.Assertions.Nil(err)

	// get valid AccessToken
	s.SetupAccessToken()
}

// SetupAccessToken fetches the jwt token used for private and restricted tests
func (s *AccountSuite) SetupAccessToken() {
	// Login data
	jsonBody := map[string]interface{}{"login": "123456789", "pwd": "r4nd0mpass123!@#", "remember": true}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(jsonBody)

	// Execute login
	w := utils.MakeRequest(s.router, http.MethodPost, "/account/login", payloadBuf)
	s.Equal(http.StatusOK, w.Code, "status should be 200, login is correct")
	s.GreaterOrEqual(len(w.Result().Cookies()), 1)

	// Fetch returned cookie
	s.accessToken = w.Result().Cookies()[0]
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountSuite))
}

func (s *AccountSuite) TestAccount() {
	s.Run("profile", s.profile)
	s.Run("signupCaptcha", s.signupCaptcha)
	s.Run("changePassword", s.changePassword)
	s.Run("logout", s.logout)
	s.Run("delete", s.delete)
}

func (s *AccountSuite) profile() {
	// Make request without any access token
	w := utils.MakeRequest(s.router, http.MethodGet, "/account/profile", nil)
	s.Equal(http.StatusUnauthorized, w.Result().StatusCode, "status should be 401 because there is no jwt")

	// Attach cookie and re run request
	w = utils.MakeRequest(s.router, http.MethodGet, "/account/profile", nil, s.accessToken)
	s.Equal(http.StatusOK, w.Result().StatusCode, "profile should return 200 because we have the login cookie")
}

func (s *AccountSuite) signupCaptcha() {
	w := utils.MakeRequest(s.router, http.MethodGet, "/account/captcha", nil)

	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal([]string{"image/png"}, w.Result().Header["Content-Type"])
}

func (s *AccountSuite) changePassword() {}

func (s *AccountSuite) logout() {}

func (s *AccountSuite) delete() {}
