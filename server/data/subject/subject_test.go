package subject

import (
	"net/http"
	"testing"
)

var cases = []struct {
	Input          string
	StatusExpected int
}{
	{
		"SSC0104", // EHCA
		http.StatusOK,
	},
	{
		"SCC0221", // ICC1
		http.StatusOK,
	},
	{
		"'--", // Bad
		http.StatusBadRequest,
	},
}

func TestGetSubjectByCode(t *testing.T) {
	for _, c := range cases {
		resp, err := http.Get("http://127.0.0.1:8080/api/subject?code=" + c.Input)
		if err != nil {
			t.Fatal(err)
		} else if resp.StatusCode != c.StatusExpected {
			t.Fail()
		} else {
			t.Log(resp)
		}
	}
}
