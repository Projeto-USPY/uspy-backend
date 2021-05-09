package private_test

import (
	"net/http"
	"testing"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
	DB          db.Env
	router      *gin.Engine
	accessToken *http.Cookie
}

// SetupTest runs before every test
func (s *UserSuite) SetupTest() {
	s.DB, s.router, s.accessToken = test.MustGetEnvironment(s.Suite)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) TestGetGrade()            {}
func (s *UserSuite) TestGetSubjectReview()    {}
func (s *UserSuite) TestUpdateSubjectReview() {}