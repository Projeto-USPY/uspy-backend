package account

import (
	"fmt"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/server/views/account"
	"github.com/Projeto-USPY/uspy-backend/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Profile retrieves the user profile from the database
func Profile(ctx *gin.Context, DB db.Env, userID string) {
	var storedUser models.User

	snap, err := DB.Restore("users", utils.SHA256(userID))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user with id %s: %s", userID, err.Error()))
		return
	}
	err = snap.DataTo(&storedUser)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind user %s data to model: %s", userID, err.Error()))
		return
	}

	storedUser.ID = userID
	account.Profile(ctx, storedUser)
}

// GetMajors retrieves the majors from a given user
func GetMajors(ctx *gin.Context, DB db.Env, userID string) {
	snaps, err := DB.RestoreCollection(fmt.Sprintf(
		"users/%s/majors",
		utils.SHA256(userID),
	))

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user majors: %s", err.Error()))
		return
	}

	majors := make([]*views.Major, 0, len(snaps))
	for _, s := range snaps {
		var storedMajor models.Major
		var storedCourse models.Course

		if err := s.DataTo(&storedMajor); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind user major: %s", err.Error()))
			return
		}

		snap, err := DB.Restore(
			"courses",
			storedMajor.Hash(),
		)

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get course name using major: %s", err.Error()))
			return
		}

		if err := snap.DataTo(&storedCourse); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind course: %s", err.Error()))
			return
		}

		majors = append(majors, views.NewMajorFromModels(
			&storedMajor,
			&storedCourse,
		))
	}

	account.GetMajors(ctx, majors)
}
