package restricted

import (
	"github.com/gin-gonic/gin"
	"github.com/tpreischadt/ProjetoJupiter/db"
	"github.com/tpreischadt/ProjetoJupiter/entity"
	"github.com/tpreischadt/ProjetoJupiter/server/models/restricted"
	"net/http"
	"strconv"
)

func GetSubjectGrades(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		sub := c.MustGet("Subject").(entity.Subject)

		buckets, err := restricted.GetGrades(DB, sub)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		avg, approval := 0.0, 0.0
		cnt := 0

		for k, v := range buckets {
			f, _ := strconv.ParseFloat(k, 64)
			avg += f * float64(v)

			if f >= 5.0 {
				approval += float64(v)
			}

			cnt += v
		}

		if len(buckets) == 0 {
			c.Status(http.StatusNotFound)
			return
		}

		avg /= float64(cnt)
		approval /= float64(cnt)

		c.JSON(http.StatusOK, gin.H{"grades": buckets, "average": avg, "approval": approval})
	}
}
