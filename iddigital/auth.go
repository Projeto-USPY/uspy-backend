/*Package iddigital contains all logic necessary to interact with uspdigital */
package iddigital

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Projeto-USPY/uspy-backend/config"
	log "github.com/sirupsen/logrus"
)

// GetPDF sends a GET request to the auth API with the auth code
func GetPDF(auth string) (*http.Response, error) {
	code := strings.ReplaceAll(auth, "-", "")
	URL := fmt.Sprintf("%s/%s", config.Env.AuthEndpoint, code)

	// retry 3 times
	for i := 0; i < 3; i++ {
		resp, err := http.Get(URL)

		if err != nil {
			log.Errorf("error getting pdf from uspdigital: %s", err.Error())
			return nil, err
		}

		if resp.StatusCode == http.StatusOK {
			return resp, nil
		} else if resp.StatusCode == http.StatusBadRequest {
			return resp, nil
		}

		log.Info("retrying to get pdf from uspdigital")
	}

	return nil, fmt.Errorf("error getting pdf from uspdigital: %s", "max retries exceeded")
}
