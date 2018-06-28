package middleware

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"golang.org/x/time/rate"
	"time"
	"sync"
	"daemonw/util"
)

func ApiCounter(c *gin.Context) {
	fmt.Printf("%v %s  %s\n", c.Request.RemoteAddr, c.Request.URL.Path, c.HandlerName())
}

type LimiterPool struct {
	pool map[string]*RateLimiter
	mu sync.RWMutex
}

func NewLimiterPool() *LimiterPool {
	lp := &LimiterPool{}
	lp.pool = make(map[string]*RateLimiter, 0)
	return lp
}

func (lp *LimiterPool) Add(key string, lmt *RateLimiter) {
	lp.mu.Lock()
	lp.pool[key] = lmt
	lp.mu.Unlock()
}

func (lp *LimiterPool) Get(key string) *RateLimiter{
	lp.mu.RLock()
	r:=lp.pool[key]
	lp.mu.RUnlock()
	return r
}

type RateLimiter struct {
	Limiter   rate.Limiter
	ip        string
	header    map[string]string
	token     string
	hitFunc   func(c *gin.Context) int
	tokenFunc func(c *gin.Context) string
}

func NewLimiter(r float64, burst int) *RateLimiter {
	lmt := &RateLimiter{}
	lmt.Limiter = *rate.NewLimiter(rate.Limit(r), burst)
	return lmt
}

func (r *RateLimiter) SetTokenFunc(f func(c *gin.Context) string) {
	r.tokenFunc = f
}

func (r *RateLimiter) LimitIp(ip string) {
	r.ip = ip
	r.hitFunc = func(c *gin.Context) int {
		clientIp := util.GetRequestIP(c.Request, false)
		if clientIp == ip {
			return 1
		}
		return 0
	}
}

func (r *RateLimiter) LimitHeader(key string, val string) {
	r.header[key] = val
	r.hitFunc = func(c *gin.Context) int {
		for k, v := range r.header {
			reqVal := c.GetHeader(k)
			if v == reqVal {
				return 1
			}
		}
		return 0
	}
}

func (r *RateLimiter) LimitToken(token string) {
	if r.tokenFunc == nil {
		panic("token function is nil!")
	}
	r.token = token
	r.hitFunc = func(c *gin.Context) int {
		clientToken := r.tokenFunc(c)
		if token == clientToken {
			return 1
		}
		return 0
	}
}

func (r *RateLimiter) LimitHeaders(headers map[string]string) {
	for k, v := range headers {
		r.header[k] = v
	}
	r.hitFunc = func(c *gin.Context) int {
		for k, v := range r.header {
			reqVal := c.GetHeader(k)
			if v == reqVal {
				return 1
			}
		}
		return 0
	}
}

func IpLimiter(lp *LimiterPool, rate float64,burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := util.GetRequestIP(c.Request, false)
		fmt.Println(ip)
		r, ok := lp.pool[ip]
		if !ok {
			r=NewLimiter(rate,burst)
			r.LimitIp(ip)
			lp.pool[ip] = r
		}
		hit := r.hitFunc(c)
		allow := r.Limiter.AllowN(time.Now(), hit)
		if allow {
			c.Next()
			return
		}
		c.JSON(http.StatusTooManyRequests, gin.H{"msg": "access too often"})
		c.Abort()
	}
}
