// package account contains the callbacks for every /account endpoint
// for backend-db communication, see /server/models/account
package account

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tpreischadt/ProjetoJupiter/utils"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/iddigital"
	"github.com/tpreischadt/ProjetoJupiter/server/middleware"
	"github.com/tpreischadt/ProjetoJupiter/server/models/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Profile returns the user's profile (in v1 it only checks for authentication, but this will be incremented later)
func Profile(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		user := entity.User{Login: userID}
		storedUser, err := account.Profile(DB, user)

		if err != nil {
			// error fetching user profile
			log.Println(fmt.Errorf("error fetching user %v profile: %s", user, err.Error()))
			c.Status(http.StatusInternalServerError)
			return
		}

		if key, ok := os.LookupEnv("AES_KEY"); ok {
			if name, err := utils.AESDecrypt(storedUser.NameHash, key); err != nil {
				log.Println("error decrypting name hash")
				c.Status(http.StatusInternalServerError)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"user": userID,
					"name": name,
				})
			}
		} else {
			log.Println(errors.New("AES_KEY 128/196/256-bit key env variable was not provided"))
			c.Status(http.StatusInternalServerError)
		}
	}
}

// ResetPassword is a closure for PUT /account/password_reset
// It differs from ChangePassword because the user does not have to be logged in.
func ResetPassword(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		// validate user data
		var signupForm entity.Signup
		if err := c.ShouldBindJSON(&signupForm); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// get user records
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

		if pdf := iddigital.NewPDF(resp); pdf.Error != nil {
			// error converting PDF to text
			log.Println(errors.New("error converting pdf to text: " + pdf.Error.Error()))
			c.Status(http.StatusInternalServerError)
		} else {
			data, err := pdf.Parse(DB)

			if err != nil {
				// error parsing pdf
				log.Println(errors.New("error parsing pdf: " + err.Error()))
				c.Status(http.StatusInternalServerError)
				return
			}

			// change user password
			user := entity.User{Login: data.Nusp}
			err = account.ChangePassword(DB, user, signupForm.Password)
			if err != nil {
				// error changing password
				log.Println(fmt.Errorf("error changing password of user %v to %v: %s", user.Login, signupForm.Password, err.Error()))
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Status(http.StatusOK)
		}
	}
}

// ChangePassword is a closure for PUT /account/password_change
// It differs from ResetPassword because the user must be logged in.
func ChangePassword(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get user info
		token := c.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		var reset entity.Reset
		// bind old and new password
		if err := c.ShouldBindJSON(&reset); err != nil {
			c.Status(http.StatusBadRequest)
		} else {
			user := entity.User{Login: userID, Password: reset.OldPassword}
			if _, loginErr := account.Login(DB, user); loginErr != nil { // old_password is incorrect
				c.Status(http.StatusForbidden)
				return
			}

			changeErr := account.ChangePassword(DB, user, reset.NewPassword)
			if changeErr != nil {
				c.Status(http.StatusInternalServerError)
				return
			}

			c.Status(http.StatusOK)
		}
	}
}

// Logout is a closure for the GET /account/logout endpoint
func Logout() func(c *gin.Context) {
	return func(c *gin.Context) {
		domain := c.MustGet("front_domain").(string)
		secureCookie := os.Getenv("LOCAL") == "FALSE"

		// delete access_token cookie
		c.SetCookie("access_token", "", -1, "/", domain, secureCookie, true)
	}
}

// Login is a closure for the POST /account/login endpoint
func Login(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		var user entity.User

		// validate user data
		if err := c.ShouldBindJSON(&user); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// check if password is correct
		if storedUser, err := account.Login(DB, user); err != nil {
			c.Status(http.StatusUnauthorized)
		} else {
			// generate access_token
			if jwtToken, err := middleware.GenerateJWT(user); err != nil {
				log.Println(fmt.Errorf("error generating jwt for user %v: %s", user, err.Error()))
				c.Status(http.StatusInternalServerError)
			} else {
				domain := c.MustGet("front_domain").(string)

				// expiration date = 1 month
				secureCookie := os.Getenv("LOCAL") == "FALSE"
				cookieAge := 0

				// remember this login?
				if user.Remember {
					cookieAge = 30 * 24 * 3600 // 30 days in seconds
				}

				if key, ok := os.LookupEnv("AES_KEY"); ok {
					if name, err := utils.AESDecrypt(storedUser.NameHash, key); err != nil {
						log.Println("error decrypting name hash")
						c.Status(http.StatusInternalServerError)
					} else {
						c.SetCookie("access_token", jwtToken, cookieAge, "/", domain, secureCookie, true)
						c.JSON(http.StatusOK, gin.H{
							"user": user.Login,
							"name": name,
						})
					}
				} else {
					log.Println(errors.New("AES_KEY 128/196/256-bit key env variable was not provided"))
					c.Status(http.StatusInternalServerError)
				}
			}
		}
	}
}

// Remove is a closure for the removes the account, forever!!
func Remove(DB db.Env) func(g *gin.Context) {
	return func(c *gin.Context) {
		domain := os.Getenv("DOMAIN")
		secureCookie := os.Getenv("MODE") == "prod"

		// delete access_token cookie
		c.SetCookie("access_token", "", -1, "/", domain, secureCookie, true)

		// get user info
		token := c.MustGet("access_token")
		claims := token.(*jwt.Token).Claims.(jwt.MapClaims)
		userID := claims["user"].(string)

		u := entity.User{Login: userID}
		if err := account.Remove(DB, u); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}

// Signup is a closure for the POST /account/create endpoint
func Signup(DB db.Env) func(g *gin.Context) {
	return func(c *gin.Context) {
		// validate user data
		var signupForm entity.Signup
		if err := c.ShouldBindJSON(&signupForm); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// get user records
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

		if pdf := iddigital.NewPDF(resp); pdf.Error != nil {
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

			newUser, userErr := entity.NewUserWithOptions(
				data.Nusp,
				signupForm.Password,
				data.Name,
				pdf.CreationDate,
				entity.WithNameHash{},
				entity.WithPasswordHash{},
			)

			if userErr != nil {
				log.Println(errors.New("error generating user" + userErr.Error()))
				c.Status(http.StatusInternalServerError)
				return
			}

			_, err = DB.Restore("users", newUser.Hash())
			if status.Code(err) == codes.NotFound {
				// user is new
				signupErr := account.Signup(DB, newUser, data)
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

			// generate JWT to auto-login user for the current session
			jwtToken, err := middleware.GenerateJWT(newUser)
			if err != nil {
				log.Println(errors.New("error generating jwt for new user: " + err.Error()))
				c.Status(http.StatusInternalServerError)
				return
			}

			domain := c.MustGet("front_domain").(string)
			secureCookie := os.Getenv("LOCAL") == "FALSE"
			c.SetCookie("access_token", jwtToken, 0, "/", domain, secureCookie, true)
			c.JSON(http.StatusOK, data)
		}
	}
}

// SignupCaptcha is a closure for the GET /account/captcha endpoint
func SignupCaptcha() func(c *gin.Context) {
	return func(c *gin.Context) {
		resp, err := iddigital.GetCaptcha()
		if err != nil || resp.StatusCode != http.StatusOK {
			log.Println(errors.New("error getting captcha from iddigital"))
			c.Status(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		cookies := resp.Cookies()
		for _, ck := range cookies {
			domain := c.MustGet("front_domain").(string)
			secureCookie := os.Getenv("LOCAL") == "FALSE"
			c.SetCookie(ck.Name, ck.Value, ck.MaxAge, "/", domain, secureCookie, true)
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
