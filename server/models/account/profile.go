package account

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

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

	snap, err := DB.Restore("users/" + utils.SHA256(userID))
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

		snap, err := DB.Restore("courses/" + storedMajor.Hash())

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

// SearchCurriculum queries the user's given major subjects and returns which ones they have completed and if so, their record information (grade, status and frequency)
func SearchCurriculum(ctx *gin.Context, DB db.Env, userID string, controller *controllers.CurriculumQuery) {
	courseSubjectIDs, err := DB.Client.Collection("subjects").
		Where("course", "==", controller.Course).
		Where("specialization", "==", controller.Specialization).
		Where("optional", "==", controller.Optional).
		Where("semester", "==", controller.Semester).
		Documents(ctx).
		GetAll()

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("error running curriculum query: %s", err.Error()))
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error running curriculum query: %s", err.Error()))
		return
	}

	userHash := utils.SHA256(userID)
	results := make([]*views.CurriculumResult, 0, len(courseSubjectIDs))

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

		result := &views.CurriculumResult{
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

	account.SearchCurriculum(ctx, results)
}

// GetTranscriptYears retrieves the last few years a user's has been in USP
func GetTranscriptYears(ctx *gin.Context, DB db.Env, userID string) {
	userHash := utils.SHA256(userID)

	// fetch all final scores from users, we cannot use restore collection here because final scores are missing documents
	finalScores, err := DB.RestoreCollectionRefs(
		fmt.Sprintf(
			"users/%s/final_scores",
			userHash,
		),
	)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user final scores: %s", err.Error()))
		return
	}

	// fetch years the user has been in USP
	years := make(map[int][]int)

	for _, fs := range finalScores {
		curYear := time.Now().Year()
		for year := curYear - 10; year <= curYear; year++ {
			for _, semester := range []int{1, 2} {
				recordHash := models.Record{Year: year, Semester: semester}.Hash()

				subHash := fs.ID

				// get final score with given record hash (year + semester)
				_, err := DB.Restore(
					fmt.Sprintf(
						"users/%s/final_scores/%s/records/%s",
						userHash,
						subHash,
						recordHash,
					),
				)

				if err != nil {
					if status.Code(err) == codes.NotFound { // no document with given record hash
						continue
					}

					ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get record: %s", err.Error()))
					return
				}

				// record was found with this year + semester

				if _, ok := years[year]; !ok { // create array if needed
					years[year] = make([]int, 0)
				}

				years[year] = append(years[year], semester)
			}
		}
	}

	flattenedYears := make([]*views.TranscriptYear, 0, len(years))
	for year, semesters := range years {
		semesters = utils.UniqueInts(semesters)
		sort.Ints(semesters)

		flattenedYears = append(flattenedYears, &views.TranscriptYear{
			Year:      year,
			Semesters: semesters,
		})
	}

	sort.Slice(flattenedYears, func(i, j int) bool {
		return flattenedYears[i].Year < flattenedYears[j].Year
	}) // sort years

	account.GetTranscriptYears(ctx, flattenedYears)
}

// SearchTranscript takes a transcript query and retrieves its records with subject data attached to them
func SearchTranscript(ctx *gin.Context, DB db.Env, userID string, controller *controllers.TranscriptQuery) {
	userHash := utils.SHA256(userID)

	// fetch all final scores from users, we cannot use restore collection here because final scores are missing documents
	finalScores, err := DB.RestoreCollectionRefs(
		fmt.Sprintf(
			"users/%s/final_scores",
			userHash,
		),
	)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get user final scores: %s", err.Error()))
		return
	}

	// get all records with same hash as transcript query data
	results := make([]*views.Record, 0, len(finalScores))
	recordHash := models.Record{Year: controller.Year, Semester: controller.Semester}.Hash()
	for _, fs := range finalScores {
		subHash := fs.ID

		// get final score with given record hash (year + semester)
		snap, err := DB.Restore(
			fmt.Sprintf(
				"users/%s/final_scores/%s/records/%s",
				userHash,
				subHash,
				recordHash,
			),
		)

		if err != nil {
			if status.Code(err) == codes.NotFound { // no document with given record hash
				continue
			}

			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get record: %s", err.Error()))
			return
		}

		var model models.Record
		if err := snap.DataTo(&model); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind record: %s", err.Error()))
			return
		}

		record := views.NewRecordFromModel(&model)

		// get subject data to inject into view object
		snap, err = DB.Restore("subjects/" + subHash)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				log.Printf("failed to get subject with hash %s in records query, this should not happen, maybe subject does not exist anymore?\n", subHash)
				continue
			}

			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to get subject with record hash: %s", err.Error()))
			return
		}

		var subject models.Subject
		if err := snap.DataTo(&subject); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind subject: %s", err.Error()))
			return
		}

		// inject subject header
		record.Code = subject.Code
		record.Course = subject.CourseCode
		record.Specialization = subject.Specialization
		record.Name = subject.Name

		results = append(results, record)
	}

	account.SearchTranscript(ctx, results)
}