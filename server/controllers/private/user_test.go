package private_test

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

type UserSuite struct {
	suite.Suite
	DB          db.Env
	router      *gin.Engine
	accessToken *http.Cookie
}

// SetupSuite runs before suite (fetches the emulator)
func (s *UserSuite) SetupSuite() {
	s.DB = emulator.MustGet()

	// setup router
	var err error
	s.router, err = server.SetupRouter(s.DB)
	s.Assertions.Nil(err)

	// get valid AccessToken
	s.SetupAccessToken()
}

// SetupAccessToken fetches the jwt token used for private and restricted tests
func (s *UserSuite) SetupAccessToken() {
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

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) TestUser() {
	s.Run("getGrade", s.getGrade)
	s.Run("getSubjectReview", s.getSubjectReview)
	s.Run("updateSubjectReview", s.updateSubjectReview)
}

func (s *UserSuite) getGrade()            {}
func (s *UserSuite) getSubjectReview()    {}
func (s *UserSuite) updateSubjectReview() {}
