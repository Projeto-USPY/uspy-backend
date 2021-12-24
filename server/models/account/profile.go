package account

import (
	"fmt"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
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

// SearchTranscript queries the user's given major subjects and returns which ones they have completed and if so, their record information (grade, status and frequency)
func SearchTranscript(ctx *gin.Context, DB db.Env, userID string, controller *controllers.TranscriptQuery) {
	courseSubjectIDs, err := DB.Client.Collection("subjects").
		Where("course", "==", controller.Course).
		Where("specialization", "==", controller.Specialization).
		Where("optional", "==", controller.Optional).
		Where("semester", "==", controller.Semester).
		Documents(ctx).
		GetAll()

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("error running transcript query: %s", err.Error()))
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error running transcript query: %s", err.Error()))
		return
	}

	userHash := utils.SHA256(userID)
	results := make([]*views.TranscriptResult, 0, len(courseSubjectIDs))

	for _, subDoc := range courseSubjectIDs {
		// query if user has done this subject
		snaps, err := DB.Client.Collection(fmt.Sprintf(
			"users/%s/final_scores/%s/records", // users/#user/final_scores/#subject/records
			userHash,
			subDoc.Ref.ID,
		)).Documents(ctx).GetAll()

		completed := len(snaps) > 0

		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting user record: %s", err.Error()))
			return
		}

		// bind subject data
		var subject models.Subject
		if err := subDoc.DataTo(&subject); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error binding subject: %s", err.Error()))
			return
		}

		result := &views.TranscriptResult{
			Name:      subject.Name,
			Code:      subject.Code,
			Completed: completed,
		}

		if completed {
			for _, recordDoc := range snaps {
				// bind record
				var record models.Record
				if err := recordDoc.DataTo(&record); err != nil {
					ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error binding record: %s", err.Error()))
					return
				}

				result.Frequency = record.Frequency
				result.Grade = record.Grade
				result.Status = record.Status
				results = append(results, result) // insert all times the user has done subject (usually this for runs only once)
			}
		} else { // insert oly once if not completed
			results = append(results, result)
		}
	}

	account.SearchTranscript(ctx, results)
}
