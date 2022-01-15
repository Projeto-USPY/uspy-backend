package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		return fmt.Sprintf("|%s %-3s %s| %v |%s %3d %s| %13v | %15s |%#v\n",
			methodColor, param.Method, resetColor,
			param.TimeStamp.Format(time.RFC3339),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			param.Path,
		)
	})
}
