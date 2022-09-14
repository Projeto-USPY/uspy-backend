package public

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/Projeto-USPY/uspy-backend/server/views/public"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetStats fetches some statistics from the database
func GetStats(ctx *gin.Context, DB db.Database) {
	// query cached stats collection
	snaps, err := DB.RestoreCollection("stats")
	if len(snaps) == 1 { // found stats document
		var stats models.Stats
		documentSnap := snaps[0]

		if err := documentSnap.DataTo(&stats); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to bind stats to model: %s", err.Error()))
			return
		}

		public.GetStats(
			ctx,
			&stats,
		)

		// update users and grades count if it's too old
		// get number of users and number of grades
		var users, grades int

		if userSnaps, err := DB.Client.Collection("users").Select().Documents(ctx).GetAll(); err != nil {
			log.Errorf("error getting user count: %s", err.Error())
			return
		} else {
			users = len(userSnaps)
		}

		if gradeSnaps, err := DB.Client.CollectionGroup("grades").Select().Documents(ctx).GetAll(); err != nil {
			log.Errorf("error getting grades count: %s", err.Error())
			return
		} else {
			grades = len(gradeSnaps)
		}

		now := time.Now().In(stats.Users.LastUpdate.Location())
		if now.Sub(stats.Users.LastUpdate) > 24*time.Hour || now.Sub(stats.Grades.LastUpdate) > 24*time.Hour { // more than a day old
			stats.Grades.LastUpdate = time.Now()
			stats.Users.LastUpdate = time.Now()

			if err := DB.Update(models.Stats{
				Users: models.StatsEntry{
					Name:       "users",
					Count:      users,
					LastUpdate: now,
				},
				Grades: models.StatsEntry{
					Name:       "grades",
					Count:      grades,
					LastUpdate: now,
				},
			}, "stats"); err != nil {
				log.Errorf("error updating user and grades stats counters: %s", err.Error())
			}
		}

		return

	} else if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error retrieving stats document: %s", err.Error()))
		return
	} else {
		ctx.Status(http.StatusInternalServerError)
	}
}
