// package restricted contains the callbacks for every restricted (allowed only to users) /api/restricted endpoint
// for backend-db communication, see /server/models/restricted
package restricted

import (
	"net/http"
	"os"
	"strconv"

	"github.com/Projeto-USPY/uspy-backend/db"
	"github.com/Projeto-USPY/uspy-backend/entity"
	"github.com/Projeto-USPY/uspy-backend/server/models/restricted"
	"github.com/gin-gonic/gin"
)

// GetSubjectGrades is a closure for the GET /api/restricted/subject/grades endpoint
func GetSubjectGrades(DB db.Env) func(c *gin.Context) {
	return func(c *gin.Context) {
		sub := c.MustGet("Subject").(entity.Subject)

		// get all grades in buckets of frequency
		buckets, err := restricted.GetGrades(DB, sub)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
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

			if os.Getenv("MODE") == "prod" && cnt <= 10 { // do not return grades if there are too few grades
				c.JSON(http.StatusOK, gin.H{"grades": map[string]int{}, "average": 0.0, "approval": 0.0})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"grades": buckets, "average": avg, "approval": approval})
	}
}
