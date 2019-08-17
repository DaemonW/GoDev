package middleware

import (
	"daemonw/entity"
	"daemonw/xerr"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func RateLimiter(limiter *entity.Limiter, n int64) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.ClientIP()
		rate, delay, allowed := limiter.Allow(key, n, time.Second)
		if !allowed {
			c.Header("X-RateLimit-Limit", strconv.FormatInt(n, 10))
			c.Header("X-RateLimit-Remaining", strconv.FormatInt(n-rate, 10))
			delaySec := int64(delay / time.Second)
			c.Header("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))
			c.JSON(http.StatusNotAcceptable, entity.NewRespErr(xerr.CodeRateLimit, xerr.MsgAccessFrequency))
			c.Abort()
			return
		}
	}
}
