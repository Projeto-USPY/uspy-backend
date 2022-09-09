package private_test

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

type UserSuite struct {
	suite.Suite
	DB          db.Database
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

func (s *UserSuite) TestGetGrade() {
	baseURL := "/private/subject/grade"
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

	// assert response JSON
	expectedResponse := `{
		"grade": 9.0,
		"status": "A",
		"frequency": 100,
		"reviewed": false
	}`

	bytes, err := io.ReadAll(w.Result().Body)
	s.NoError(err)
	s.JSONEq(expectedResponse, string(bytes))
}

/** This has not been implemented yet because the Firestore emulator does not support transactional gets **/
//func (s *UserSuite) TestGetSubjectReview() {}
//func (s *UserSuite) TestUpdateSubjectReview() {}
