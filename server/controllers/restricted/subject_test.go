package restricted_test

import (
	"net/http"
	"testing"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type SubjectSuite struct {
	suite.Suite
	DB          db.Env
	router      *gin.Engine
	accessToken *http.Cookie
}

// SetupTest runs before every test
func (s *SubjectSuite) SetupTest() {
	s.DB, s.router, s.accessToken = test.MustGetEnvironment()
}

func TestSubjectSuite(t *testing.T) {
	suite.Run(t, new(SubjectSuite))
}

func (s *SubjectSuite) TestGetGrades() {}
