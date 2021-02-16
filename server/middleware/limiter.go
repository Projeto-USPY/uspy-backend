package middleware

import (
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
)

// RateLimiter returns a middleware that sets up the rate limiting tool
func RateLimiter(format string) gin.HandlerFunc {
	rater, err := limiter.NewRateFromFormatted(format)
	if err != nil {
		return nil
	}

	store := memory.NewStore()
	instance := limiter.New(store, rater)

	return mgin.NewMiddleware(instance)
}
