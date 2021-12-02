package public

import (
	"fmt"
	"net/http"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/public"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetOfferings(ctx *gin.Context, DB db.Env, sub *controllers.Subject) {
	model := models.NewSubjectFromController(sub)

	offerings := make([]*models.Offering, 0, 20)
	IDs := make([]string, 0, 20)

	snaps, err := DB.RestoreCollection("subjects/" + model.Hash() + "/offerings")

	if err != nil {
		if status.Code(err) == codes.NotFound {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("could not find collection offerings: %s", err.Error()))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to fetch offerings: %s", err.Error()))
		return
	} else if len(snaps) == 0 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	for _, s := range snaps {
		var off models.Offering
		if err := s.DataTo(&off); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not bind offering: %s", err.Error()))
			return
		}

		offerings = append(offerings, &off)
		IDs = append(IDs, s.Ref.ID)
	}

	public.GetOfferings(ctx, IDs, offerings)
}
