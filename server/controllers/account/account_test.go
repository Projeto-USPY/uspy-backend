package account_test

import (
	"net/http"
	"testing"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/Projeto-USPY/uspy-backend/utils/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AccountSuite struct {
	suite.Suite
	DB          db.Env
	router      *gin.Engine
	accessToken *http.Cookie
}

// SetupTest runs before every test
func (s *AccountSuite) SetupTest() {
	s.DB, s.router, s.accessToken = test.MustGetEnvironment()
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountSuite))
}

func (s *AccountSuite) TestProfile() {
	// Make request without any access token
	w := utils.MakeRequest(s.router, http.MethodGet, "/account/profile", nil)
	s.Equal(http.StatusUnauthorized, w.Result().StatusCode, "status should be 401 because there is no jwt")

	// Attach cookie and re run request
	w = utils.MakeRequest(s.router, http.MethodGet, "/account/profile", nil, s.accessToken)
	s.Equal(http.StatusOK, w.Result().StatusCode, "profile should return 200 because we have the login cookie")
}

func (s *AccountSuite) TestSignupCaptcha() {
	w := utils.MakeRequest(s.router, http.MethodGet, "/account/captcha", nil)

	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal([]string{"image/png"}, w.Result().Header["Content-Type"])
}

func (s *AccountSuite) TestChangePassword() {}
func (s *AccountSuite) TestLogout()         {}
func (s *AccountSuite) TestDelete()         {}
