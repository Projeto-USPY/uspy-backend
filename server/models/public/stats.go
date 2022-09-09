package public

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/views"
	"github.com/Projeto-USPY/uspy-backend/server/views/public"
	"github.com/gin-gonic/gin"
)

// GetStats fetches some statistics from the database
//
// This function does not scale well.
// TODO: Switch to a manual counter later on
func GetStats(ctx *gin.Context, DB db.Database) {
	statsChan := make(chan *views.StatsEntry, 5)

	performQuery := func(ctx *gin.Context, category string, action func() ([]*firestore.DocumentSnapshot, error)) {
		snaps, err := action()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error getting stats for: %s", category))
		}

		statsChan <- views.NewStatsEntry(
			category,
			len(snaps),
		)
	}

	go performQuery(ctx, "users", DB.Client.Collection("users").Select().Documents(ctx).GetAll)
	go performQuery(ctx, "subjects", DB.Client.Collection("subjects").Select().Documents(ctx).GetAll)
	go performQuery(ctx, "grades", DB.Client.CollectionGroup("grades").Select().Documents(ctx).GetAll)
	go performQuery(ctx, "comments", DB.Client.CollectionGroup("comments").Select().Documents(ctx).GetAll)
	go performQuery(ctx, "offerings", DB.Client.CollectionGroup("offerings").Select().Documents(ctx).GetAll)

	public.GetStats(
		ctx,
		statsChan,
	)
}
