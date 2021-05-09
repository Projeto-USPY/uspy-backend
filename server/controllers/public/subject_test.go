package public_test

import (
	"fmt"
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
	DB          db.Env
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

func (s *SubjectSuite) TestGetAll() {
	w := utils.MakeRequest(s.router, http.MethodGet, "/api/subject/all", nil)

	expectedResponse := `
		[
			{
				"name": "Bacharelado em Ciência de Dados",
				"code": "55090",
				"specialization": "0",
				"subjects": {
					"SCC0230": "Inteligência Artificial"
				}
			},
			{
				"name": "Bacharelado em Ciências de Computação",
				"code": "55041",
				"specialization": "0",
				"subjects": {
					"SCC0222": "Laboratório de Introdução à Ciência de Computação I",
					"SCC0217": "Linguages de Programação e Compiladores"
				}
			}
		]
	`

	s.Equal(http.StatusOK, w.Result().StatusCode)

	bytes, err := io.ReadAll(w.Result().Body)
	s.NoError(err)
	s.JSONEq(expectedResponse, string(bytes))
}
func (s *SubjectSuite) TestGetByCode() {
	endpointMask := "/api/subject?code=%s&course=%s&specialization=%s"

	// invalid request
	w := utils.MakeRequest(
		s.router,
		http.MethodGet,
		fmt.Sprintf(
			endpointMask,
			"'--",
			"55041",
			"0",
		),
		nil,
	)

	s.Equal(http.StatusBadRequest, w.Result().StatusCode, "subject is invalid")

	// subject does not exist
	w = utils.MakeRequest(
		s.router,
		http.MethodGet,
		fmt.Sprintf(
			endpointMask,
			"SCC0999",
			"55041",
			"0",
		),
		nil,
	)

	s.Equal(http.StatusNotFound, w.Result().StatusCode, "subject does not exist")

	// OK
	w = utils.MakeRequest(
		s.router,
		http.MethodGet,
		fmt.Sprintf(
			endpointMask,
			"SCC0222",
			"55041",
			"0",
		),
		nil,
	)

	s.Equal(http.StatusOK, w.Result().StatusCode, "subject should be returned")

	expectedResponse := `
		{
			"code": "SCC0222",
			"course": "55041",
			"specialization": "0",
			"semester": 2,
			"name": "Laboratório de Introdução à Ciência de Computação I",
			"description": "Implementar em laboratório as técnicas de programação apresentadas em Introdução à Ciência da Computação I, utilizando uma linguagem de programação estruturada.",
			"class": 2,
			"assign": 2,
			"hours": "90 h",
			"optional": true,
			"stats": {
				"total": 0,
				"worth_it": 0
			},
			"requirements": []
		}
	`

	bytes, err := io.ReadAll(w.Result().Body)
	s.NoError(err)
	s.JSONEq(expectedResponse, string(bytes))
}
func (s *SubjectSuite) TestGetGraph() {}
