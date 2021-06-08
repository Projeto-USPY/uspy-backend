package restricted

import (
	"fmt"
	"net/http"
	"sync"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/controllers"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/restricted"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetOfferings is a closure for the GET /api/restricted/offerings endpoint
func GetOfferingsWithStats(ctx *gin.Context, DB db.Env, sub *controllers.Subject) {
	subModel := models.NewSubjectFromController(sub)

	offerings := make([]*models.Offering, 0, 20)
	IDs := make([]string, 0, 20)
	stats := make([]*models.OfferingStats, 0, 20)

	offeringsMask := "subjects/%s/offerings"
	offeringsPath := fmt.Sprintf(offeringsMask, subModel.Hash())

	snaps, err := DB.RestoreCollection(offeringsPath)

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

	var wg sync.WaitGroup
	offsChannel := make(chan *models.Offering, len(snaps))
	IDchan := make(chan string, len(snaps))
	statsChan := make(chan *models.OfferingStats, len(snaps))

	for _, s := range snaps {
		wg.Add(1)
		go func(snap *firestore.DocumentSnapshot, wg *sync.WaitGroup) {
			defer wg.Done()
			posQt := make(map[string]int)
			negQt := make(map[string]int)
			neutQt := make(map[string]int)

			var off models.Offering
			if err := snap.DataTo(&off); err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("could not bind offering: %s", err.Error()))
				return
			}

			commentsPath := fmt.Sprintf("%s/%s/comments", offeringsPath, snap.Ref.ID)
			commentsCol := DB.Client.Collection(commentsPath)

			queryComments := func(op string) (int, int, error) {
				query := commentsCol.Where("rating", op, "3")
				res, err := query.Documents(DB.Ctx).GetAll()
				if err != nil {
					if status.Code(err) == codes.NotFound {
						return -1, http.StatusNotFound, fmt.Errorf("could not find collection comments: %s", err.Error())
					}
					return -1, http.StatusInternalServerError, fmt.Errorf("failed to fetch comments: %s", err.Error())
				}

				return len(res), -1, nil
			}

			for _, op := range []string{"<", "==", ">"} {
				qt, status, err := queryComments(op)
				if err != nil {
					ctx.AbortWithError(status, err)
					return
				}

				switch op {
				case ">":
					posQt[snap.Ref.ID] += qt
				case "==":
					neutQt[snap.Ref.ID] += qt
				case "<":
					negQt[snap.Ref.ID] += qt
				}
			}

			offsChannel <- &off
			IDchan <- snap.Ref.ID
			statsChan <- &models.OfferingStats{
				Approval:    posQt[snap.Ref.ID],
				Disapproval: negQt[snap.Ref.ID],
				Neutral:     neutQt[snap.Ref.ID],
			}
		}(s, &wg)
	}

	wg.Wait()

	close(offsChannel)
	close(IDchan)
	close(statsChan)

	for i := 0; i < len(snaps); i++ {
		offerings = append(offerings, <-offsChannel)
		IDs = append(IDs, <-IDchan)
		stats = append(stats, <-statsChan)
	}

	limit := len(IDs)
	if sub.Limit > 0 {
		limit = sub.Limit
	}

	restricted.GetOfferingsWithStats(ctx, IDs, offerings, stats, limit)
}
