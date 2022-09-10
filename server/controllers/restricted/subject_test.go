package restricted_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/Projeto-USPY/uspy-backend/utils/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type SubjectSuite struct {
	suite.Suite
	DB          db.Database
	router      *gin.Engine
	accessToken *http.Cookie
}

// SetupTest runs before every test
func (s *SubjectSuite) SetupTest() {
	s.DB, s.router, s.accessToken = test.MustGetEnvironment(s.Suite)
}

func TestSubjectSuite(t *testing.T) {
	suite.Run(t, new(SubjectSuite))
}

func (s *SubjectSuite) TestGetGrades() {
	baseURL := "/api/restricted/subject/grades"
	correctQueryParams := "?code=SCC0217&course=55041&specialization=0"
	incorrectQueryParams := "?code=SCC0999&course=55041&specialization=0"

	// no access token
	w := utils.MakeRequest(s.router, http.MethodGet, baseURL+correctQueryParams, nil)
	s.Equal(http.StatusUnauthorized, w.Result().StatusCode, "access token is not set so it should not return grade")

	// with access token, but invalid parameters
	w = utils.MakeRequest(s.router, http.MethodGet, baseURL+incorrectQueryParams, nil, s.accessToken)
	s.Equal(http.StatusNotFound, w.Result().StatusCode, "access token is set, but subject is not correct")

	// with access token
	w = utils.MakeRequest(s.router, http.MethodGet, baseURL+correctQueryParams, nil, s.accessToken)
	s.Equal(http.StatusOK, w.Result().StatusCode, "access token is set and subject is correct")

	// assert response
	expectedResponse := `{
		"grades": {
			"4.0": 1,
			"9.0": 2
		},
		"average": 7.333333333333333,
		"approval":  0.6666666666666666
	}`

	bytes, err := io.ReadAll(w.Result().Body)
	s.NoError(err)
	s.JSONEq(expectedResponse, string(bytes))
}
