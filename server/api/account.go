package api

import (
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/pdfparser"
	iddigital "github.com/tpreischadt/ProjetoJupiter/pdfparser/auth"
	"github.com/tpreischadt/ProjetoJupiter/server/auth"
	"log"
	"net/http"
	"os"
	"time"
)

func Login(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user entity.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.Status(http.StatusBadRequest)
		}

		if err := auth.Login(user); err != nil {
			c.Status(http.StatusUnauthorized)
		}

		if jwt, err := auth.GenerateJWT(user); err != nil {
			c.Status(http.StatusInternalServerError)
		} else {
			domain := os.Getenv("DOMAIN")

			// expiration date = 1 monthFirestore
			secureCookie := os.Getenv("MODE") == "prod"
			c.SetCookie("access_token", jwt, 30*24*3600, "/", domain, secureCookie, true)
			c.Status(http.StatusOK)
		}
	}
}

func Signup(DB db.Env) func(g *gin.Context) {
	return func(c *gin.Context) {
		var signupForm entity.Signup
		if err := c.ShouldBindJSON(&signupForm); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		cookies := c.Request.Cookies()
		resp, err := iddigital.PostAuthCode(signupForm.AccessKey, signupForm.Captcha, cookies)
		if err != nil {
			// error getting PDF from iddigital
			c.Status(http.StatusInternalServerError)
			return
		} else if resp.Header.Get("Content-Type") != "application/pdf" {
			// wrong captcha or auth
			c.Status(http.StatusBadRequest)
			return
		}

		if pdf := pdfparser.NewPDF(resp); pdf.Error != nil {
			// error converting PDF to text
			log.Print(pdf.Error)
			c.Status(http.StatusInternalServerError)
		} else {
			data, err := pdf.ParsePDF()
			if err != nil {
				// error parsing pdf
				c.Status(http.StatusInternalServerError)
				return
			} else if time.Since(pdf.CreationDate).Hours() > 1.0 {
				// pdf is too old
				c.Status(http.StatusBadRequest)
				return
			}
			// TODO: add user to firestore
			c.JSON(http.StatusOK, data)
		}
	}
}

func SignupCaptcha() func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := iddigital.GetCaptcha()
		if err != nil || resp.StatusCode != http.StatusOK {
			c.Status(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		cookies := resp.Cookies()
		for _, ck := range cookies {
			domain := os.Getenv("DOMAIN")
			secureCookie := os.Getenv("MODE") == "prod"
			c.SetCookie(ck.Name, ck.Value, ck.MaxAge, "/account/create", domain, secureCookie, ck.HttpOnly)
		}

		c.DataFromReader(
			http.StatusOK,
			resp.ContentLength,
			resp.Header.Get("Content-Type"),
			resp.Body,
			map[string]string{},
		)

	}
}
