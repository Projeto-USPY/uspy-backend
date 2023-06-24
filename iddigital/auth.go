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

	log.Debugf("Getting PDF at URL", URL)
	return http.Get(URL)
}
