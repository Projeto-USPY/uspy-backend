package account

import (
	"fmt"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
)

// Profile sets the profile data once it is successful
func Profile(ctx *gin.Context, user models.User) {
	if name, err := utils.AESDecrypt(user.NameHash, config.Env.AESKey); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error decrypting nameHash: %s", err.Error()))
	} else {
		ctx.JSON(http.StatusOK, views.Profile{User: user.ID, Name: name})
	}
}

