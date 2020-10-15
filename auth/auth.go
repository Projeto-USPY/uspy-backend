package auth

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

// GetCaptcha returns a new auth captcha for uspiddigital
func GetCaptcha() (*http.Response, error) {
	captchaURL := "https://uspdigital.usp.br/iddigital/CriarImagemTuring"

	resp, err := http.Get(captchaURL)

	if err != nil {
		return nil, fmt.Errorf("unable to get captcha: %v", err)
	}

	return resp, nil
}

// PostAuthCode submits the auth code along with the captcha to uspiddigital
// The response object will contain the Grades PDF
func PostAuthCode(auth string, captcha string, cookies []*http.Cookie) (*http.Response, error) {
	fields := strings.Split(auth, "-")
	postURL := "https://uspdigital.usp.br/iddigital/mostradocweb"

	parsedURL, err := url.Parse(postURL)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	jar.SetCookies(parsedURL, cookies)
	if err != nil {
		return nil, fmt.Errorf("unable to create cookie jar: %v", err)
	}

	client := &http.Client{Jar: jar}
	data := url.Values{}
	data.Set("chars", strings.TrimSpace(captcha))

	for i, v := range fields {
		data.Set("codctl"+strconv.Itoa(i+1), v)
	}

	resp, err := client.PostForm(postURL, data)
	if err != nil {
		return nil, fmt.Errorf("unable to submit captcha and/or auth: %v", err)
	}

	return resp, nil
}
