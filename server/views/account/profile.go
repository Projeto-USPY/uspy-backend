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

// GetMajors sets the major list as the response
func GetMajors(ctx *gin.Context, majors []*views.Major) {
	ctx.JSON(http.StatusOK, majors)
}

// SearchCurriculum sets the curriculum results as the response
func SearchCurriculum(ctx *gin.Context, results []*views.CurriculumResult) {
	ctx.JSON(http.StatusOK, results)
}

// SearchTranscript sets the transcript results as the response
func SearchTranscript(ctx *gin.Context, results []*views.Record) {
	ctx.JSON(http.StatusOK, results)
}
