package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/tpreischadt/ProjetoJupiter/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

var cases = []struct {
	Subject        string
	Course         string
	StatusExpected int
}{
	{
		"SSC0952",
		"55051",
		http.StatusOK,
	},
	{
		"SME0221",
		"55030",
		http.StatusOK,
	},
	{
		"SME0221",
		"",
		http.StatusBadRequest,
	},
}

func TestGetSubjectByCode(t *testing.T) {
	r := server.SetupRouter(
		server.SetupDB(
			"C:\\Users\\srtp-\\GolandProjects\\ProjetoJupiter\\.env",
		),
	)

	for _, c := range cases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET",
			fmt.Sprintf("/api/subject?code=%s&course=%s", c.Subject, c.Course),
			nil,
		)
		r.ServeHTTP(w, req)
		assert.Equal(t, w.Code, c.StatusExpected)
	}
}
