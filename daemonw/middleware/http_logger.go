package middleware

import (
	"daemonw/xlog"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger() gin.HandlerFunc {
	return LoggerWithWriter()
}

// LoggerWithWriter instance a Logger middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(ignore ...string) gin.HandlerFunc {

	var skip map[string]struct{}
	if length := len(ignore); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range ignore {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		method := c.Request.Method
		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)
			clientIP := c.ClientIP()
			statusCode := c.Writer.Status()
			if raw != "" {
				path = path + "?" + raw
			}
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()
			xlog.Debug().Str("scope", "http_router").
				Msgf("%s    %d    %10v    %s    %s    %s    %s",
					end.Format("2006/01/02 - 15:04:05"),
					statusCode,
					latency,
					clientIP,
					method,
					path,
					comment,
				)
		}
	}
}
