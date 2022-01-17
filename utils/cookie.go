package utils

import (
	"github.com/gin-gonic/gin"
)

// RemoveAccessToken deletes the user access token
func RemoveAccessToken(ctx *gin.Context, secureCookie bool) {
	domain := ctx.MustGet("front_domain").(string)

	// delete access_token cookie
	ctx.SetCookie("access_token", "", -1, "/", domain, secureCookie, true)
}
