/* package iddigital contains all logic necessary to interact with uspdigital */
package iddigital

import (
	"net/http/httputil"
	"testing"
)

func TestGetCaptcha(t *testing.T) {
	response, err := GetCaptcha()

	if err != nil {
		t.Fatalf("failed with error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		t.Fatal("could not get captcha")
	}

	dump, err := httputil.DumpResponse(response, true)

	if err != nil {
		t.Fatal("could not dump response")
	}

	t.Log(string(dump))
}
