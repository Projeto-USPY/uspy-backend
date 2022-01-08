package account_test

import (
	"io"
	"net/http"
	"strings"
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
	s.DB, s.router, s.accessToken = test.MustGetEnvironment(s.Suite)
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

func (s *AccountSuite) TestGetMajors() {
	// Make request without any access token
	w := utils.MakeRequest(s.router, http.MethodGet, "/account/profile/majors", nil)
	s.Equal(http.StatusUnauthorized, w.Result().StatusCode, "status should be 401 because there is no jwt")

	expectedResponse := `
		[
			{
				"name": "Bacharelado em Ciências de Computação",
				"code": "55041",
				"specialization": "0"
			}
		]
	`

	// Attach cookie and re run request
	w = utils.MakeRequest(s.router, http.MethodGet, "/account/profile/majors", nil, s.accessToken)
	s.Equal(http.StatusOK, w.Result().StatusCode, "profile should return 200 because we have the login cookie")

	// expect correct response
	bytes, err := io.ReadAll(w.Result().Body)
	s.NoError(err)
	s.JSONEq(expectedResponse, string(bytes))
}

func (s *AccountSuite) TestSignupCaptcha() {
	w := utils.MakeRequest(s.router, http.MethodGet, "/account/captcha", nil)

	s.Equal(http.StatusOK, w.Result().StatusCode)
	s.Equal([]string{"image/png"}, w.Result().Header["Content-Type"])
}

func (s *AccountSuite) TestChangePassword() {
	loginBody := `
		{
			"login": "123456789",
			"pwd": "r4nd0mpass123!@#"
		}
	`

	newLoginBody := `
	{
		"login": "123456789",
		"pwd": "p4ssw0rdr4nd0m123!@#"
	}
	`

	incorrectBody := `
		{
			"old_password": "wr0ngpass123!@#",
			"new_password": "p4ssw0rdr4nd0m123!@#"
		}
	`

	invalidBody := `
		{
			"old_password": "r4nd0mpass123!@#",
			"new_password": "shortpass"
		}
	`

	changePwdBody := `
		{
			"old_password": "r4nd0mpass123!@#",
			"new_password": "p4ssw0rdr4nd0m123!@#"
		}
	`

	w := utils.MakeRequest(s.router, http.MethodPost, "/account/login", strings.NewReader(loginBody))
	s.Equal(http.StatusOK, w.Result().StatusCode, "failed to login with original credentials")

	w = utils.MakeRequest(s.router, http.MethodPut, "/account/password_change", strings.NewReader(incorrectBody), s.accessToken)
	s.Equal(http.StatusForbidden, w.Result().StatusCode, "changed password with wrong old pass")

	w = utils.MakeRequest(s.router, http.MethodPut, "/account/password_change", strings.NewReader(invalidBody), s.accessToken)
	s.Equal(http.StatusBadRequest, w.Result().StatusCode, "password should be invalid")

	w = utils.MakeRequest(s.router, http.MethodPut, "/account/password_change", strings.NewReader(changePwdBody), s.accessToken)
	s.Equal(http.StatusOK, w.Result().StatusCode, "failed to change password")

	w = utils.MakeRequest(s.router, http.MethodPost, "/account/login", strings.NewReader(loginBody))
	s.Equal(http.StatusUnauthorized, w.Result().StatusCode, "managed to login with old credentials")

	w = utils.MakeRequest(s.router, http.MethodPost, "/account/login", strings.NewReader(newLoginBody))
	s.Equal(http.StatusOK, w.Result().StatusCode, "failed to login with new credentials")
}
func (s *AccountSuite) TestLogout() {
	w := utils.MakeRequest(s.router, http.MethodGet, "/account/logout", nil)
	s.Equal(http.StatusUnauthorized, w.Result().StatusCode, "managed to logout without authorization")

	w = utils.MakeRequest(s.router, http.MethodGet, "/account/logout", nil, s.accessToken)
	s.Equal(http.StatusOK, w.Result().StatusCode, "did not manage to log out")

	// no cookies for you
	cookies := w.Result().Cookies()
	if len(cookies) > 0 {
		for _, c := range cookies {
			if c.Name == "access_token" {
				s.Equal(cookies[0].Value, "")
			}
		}
	}
}

/** This has not been implemented yet because the Firestore emulator does not support transactional gets **/
// func (s *AccountSuite) TestDelete() {
// 	w := utils.MakeRequest(s.router, http.MethodDelete, "/account", nil)
// 	s.Equal(http.StatusUnauthorized, w.Result().StatusCode, "managed to delete account without authorization")

// 	w = utils.MakeRequest(s.router, http.MethodDelete, "/account", nil, s.accessToken)
// 	s.Equal(http.StatusOK, w.Result().StatusCode, "could not delete account even with cookie")

// 	user, _ := entity.NewUserWithOptions(
// 		"123456789",
// 		"r4nd0mpass123!@#",
// 		"Usuário teste",
// 		time.Now(),
// 		entity.WithPasswordHash{},
// 		entity.WithNameHash{},
// 	)

// 	// assert user does not exist anymore
// 	_, err := s.DB.Restore("users", user.Hash())
// 	s.Error(err)

// 	// get all subjects
// 	subs, err := s.DB.RestoreCollection("subjects")
// 	s.NoError(err)

// 	for _, sub := range subs {
// 		colRef := sub.Ref.Collection("grades")
// 		snaps, err := colRef.Documents(s.DB.Ctx).GetAll()

// 		s.NoError(err)
// 		s.Empty(snaps)
// 	}
// }
