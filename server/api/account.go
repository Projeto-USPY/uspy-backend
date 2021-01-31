package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/iddigital"
	"github.com/tpreischadt/ProjetoJupiter/server/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			log.Println(fmt.Errorf("error generating jwt for user %v: %s", user, err.Error()))
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
			log.Println(errors.New("error getting pdf from iddigital: " + err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		} else if resp.Header.Get("Content-Type") != "application/pdf" {
			// wrong captcha or auth
			c.Status(http.StatusBadRequest)
			return
		}

		if pdf := entity.NewPDF(resp); pdf.Error != nil {
			// error converting PDF to text
			log.Println(errors.New("error converting pdf to text: " + pdf.Error.Error()))
			c.Status(http.StatusInternalServerError)
		} else {
			data, err := pdf.Parse(DB)

			var maxPDFAge float64
			if os.Getenv("MODE") == "dev" {
				maxPDFAge = 24 * 30 // a month
			} else {
				maxPDFAge = 1.0 // an hour
			}

			if err != nil {
				// error parsing pdf
				log.Println(errors.New("error parsing pdf: " + err.Error()))
				c.Status(http.StatusInternalServerError)
				return
			} else if time.Since(pdf.CreationDate).Hours() > maxPDFAge {
				// pdf is too old
				c.Status(http.StatusBadRequest)
				return
			}

			// user doent exist, must register
			newUser, hashErr := entity.User{
				Login:      data.Nusp,
				Password:   signupForm.Password,
				LastUpdate: time.Now(),
			}.WithHash()

			if hashErr != nil {
				log.Println(errors.New("error hashing password" + hashErr.Error()))
				c.Status(http.StatusInternalServerError)
				return
			}

			_, err = DB.Restore("users", newUser.Hash())
			if status.Code(err) == codes.NotFound {
				// user is new
				signupErr := newUser.Signup(DB, data)
				if signupErr != nil {
					log.Println(errors.New("error inserting user into db: " + signupErr.Error()))
					c.Status(http.StatusInternalServerError)
					return
				}
			} else {
				// user has already registered
				c.Status(http.StatusForbidden)
				return
			}

			c.JSON(http.StatusOK, data)
		}
	}
}

func SignupCaptcha() func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := iddigital.GetCaptcha()
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Println(errors.New("error getting captcha from iddigital: " + err.Error()))
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
