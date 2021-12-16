package restricted

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Projeto-USPY/uspy-backend/config"
	"github.com/Projeto-USPY/uspy-backend/entity/models"
	"github.com/gin-gonic/gin"
)

// GetGrades takes the grades models and returns its response view object by creating buckets of frequency
//
// For example a collection of grades {7.0, 7.5, 10.0, 10.0} will become {0.0: 0, ..., 7.0: 1, ...,  7.5: 1, ..., 10.0: 2}
func GetGrades(ctx *gin.Context, grades []models.Record) {
	buckets := make(map[string]int)

	for _, g := range grades {
		buckets[fmt.Sprintf("%.1f", g.Grade)]++
	}

	avg, approval := 0.0, 0.0
	cnt := 0

	// calculate average grade and approval rate
	for k, v := range buckets {
		f, _ := strconv.ParseFloat(k, 64)
		avg += f * float64(v)

		if f >= 5.0 {
			approval += float64(v)
		}

		cnt += v
	}

	if cnt > 0 {
		avg /= float64(cnt)
		approval /= float64(cnt)

		if config.Env.IsProd() && cnt <= 10 { // do not return grades if there are too few grades
			ctx.JSON(http.StatusOK, gin.H{"grades": map[string]int{}, "average": 0.0, "approval": 0.0})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"grades": buckets, "average": avg, "approval": approval})
}
