package handler

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type redisRateLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

var rateLimitScript = redis.NewScript(`
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local member = ARGV[4]

redis.call("ZREMRANGEBYSCORE", key, "-inf", now - window)

local count = redis.call("ZCARD", key)
if count < limit then
	redis.call("ZADD", key, now, member)
	redis.call("PEXPIRE", key, window)
	return 1
end
return 0
`)

func newRedisRateLimiter(client *redis.Client, limit int, window time.Duration) *redisRateLimiter {
	return &redisRateLimiter{client: client, limit: limit, window: window}
}

func (l *redisRateLimiter) allow(ctx context.Context, key string) (bool, error) {
	now := time.Now().UnixMilli()
	windowMs := l.window.Milliseconds()
	member := fmt.Sprintf("%d", now)

	res, err := rateLimitScript.Run(ctx, l.client, []string{"ratelimit: " + key}, now, windowMs, l.limit, member).Int()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func rateLimit(limiter *redisRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			ip = c.Request.RemoteAddr
		}
		allowed, err := limiter.allow(c.Request.Context(), ip)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, map[string]string{"message": "service unavailable"})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, map[string]string{"message": "too many requests, try later"})
			return
		}
		c.Next()
	}
}
