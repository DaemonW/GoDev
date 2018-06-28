package middleware

import (
	"github.com/gin-gonic/gin"
	"daemonw/db"
	"time"
	"strconv"
	"net/http"
	"daemonw/model"
	"daemonw/log"
)

func UserRateLimiter(limit int64) gin.HandlerFunc {
	limiter := NewLimiter(db.GetRing())
	return func(c *gin.Context) {
		uid := c.MustGet("uid").(uint64)
		key := "user_access:" + strconv.FormatUint(uid, 10)
		rate, delay, allow := limiter.Allow(key, limit, time.Second)
		if !allow {
			header := c.Writer.Header()
			header.Set("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
			header.Set("X-RateLimit-Remaining", strconv.FormatInt(limit-rate, 10))
			delaySec := int64(delay / time.Second)
			header.Set("X-RateLimit-Delay", strconv.FormatInt(delaySec, 10))
			c.JSON(http.StatusTooManyRequests, model.NewResp().SetErrMsg("access too frequently"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func UserCountLimiter(limit int64, dur time.Duration) gin.HandlerFunc {
	limiter := NewCounter(db.GetRedis(), 100, time.Now().Add(dur))
	return func(c *gin.Context) {
		uid := c.MustGet("uid").(uint64)
		username := c.MustGet("user").(string)
		key := "user_access:" + strconv.FormatUint(uid, 10)
		count, allow := limiter.Allow(key, time.Now())
		log.Info().Str("user", username).Msgf("access count = %d", count)
		if !allow {
			header := c.Writer.Header()
			header.Set("X-CountLimit-Limit", strconv.FormatInt(limit, 10))
			c.JSON(http.StatusTooManyRequests, model.NewResp().SetErrMsg("access too many times"))
			c.Abort()
			return
		}
		c.Next()
	}
}
