package public

import (
	"fmt"
	"net/http"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/public"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetProfessor(ctx *gin.Context, DB db.Database, controller *controllers.Professor) {
	instituteHash := models.Institute{Code: controller.Institute}.Hash()
	snap, err := DB.Restore(
		fmt.Sprintf(
			"institutes/%s/professors/%s",
			instituteHash,
			controller.Hash,
		),
	)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find professor %v: %s", controller, err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch professor: %s", err))
		return
	}

	var model models.Professor
	err = snap.DataTo(&model)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind professor to object: %s", err))
		return
	}

	// get stats for professor
	snaps, err := DB.RestoreCollection(
		fmt.Sprintf(
			"institutes/%s/professors/%s/professor_reviews",
			instituteHash,
			controller.Hash,
		),
	)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find professor %v reviews: %s", controller, err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch professor reviews: %s", err))
		return
	}

	for _, s := range snaps {
		var review models.ProfessorReview
		err := s.DataTo(&review)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind professor review to object: %s", err))
			return
		}

		model.Reviews = append(model.Reviews, review)
	}

	public.GetProfessor(ctx, &model)
}

// GetProfessorOfferings retrieves the offerings associated to a given professor
func GetProfessorOfferings(ctx *gin.Context, DB db.Database, prof *controllers.Professor) {
	// fetch professor name
	instituteHash := models.Institute{Code: prof.Institute}.Hash()
	snap, err := DB.Restore(
		fmt.Sprintf(
			"institutes/%s/professors/%s",
			instituteHash,
			prof.Hash,
		),
	)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find professor %v: %s", prof, err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch professor: %s", err))
		return
	}

	name, err := snap.DataAt("name")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind professor to object: %s", err))
		return
	}

	// query professor offerings using professor name
	snaps, err := DB.Client.CollectionGroup("offerings").Where("professor", "==", name).Documents(ctx).GetAll()
	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find user offerings on given professor: %s", err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch user offerings on given professor: %s", err.Error()))
		return
	} else if len(snaps) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	offerings := make([]*models.Offering, len(snaps))

	var wg sync.WaitGroup
	wg.Add(len(snaps))

	for idx, s := range snaps {
		go func(i int, s *firestore.DocumentSnapshot) {
			defer wg.Done()
			var off models.Offering

			if err := s.DataTo(&off); err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not bind data to comment: %s", err.Error()))
				return
			}

			subjectSnap, err := s.Ref.Parent.Parent.Get(ctx) // subject ref
			if err != nil {
				ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not get subject for given offering: %s", err.Error()))
				return
			}

			err = subjectSnap.DataTo(&off.Subject)
			if err != nil {
				ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not bind subject for given offering: %s", err.Error()))
				return
			}

			offerings[i] = &off
		}(idx, s)
	}

	wg.Wait()
	public.GetProfessorOfferings(ctx, prof.Hash, offerings)
}
