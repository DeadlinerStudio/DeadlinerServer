package middleware

import (
	"context"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type tokenBucket struct {
	tokens   float64
	lastSeen time.Time
}

type clientRateLimiter struct {
	mu               sync.Mutex
	buckets          map[string]*tokenBucket
	tokensPerSecond  float64
	burst            float64
	staleAfter       time.Duration
	cleanupThreshold int
}

func LimitByClientIP(scope string, requestsPerMinute int, burst int) app.HandlerFunc {
	if requestsPerMinute <= 0 || burst <= 0 {
		return func(ctx context.Context, c *app.RequestContext) {
			c.Next(ctx)
		}
	}

	limiter := &clientRateLimiter{
		buckets:          make(map[string]*tokenBucket),
		tokensPerSecond:  float64(requestsPerMinute) / 60.0,
		burst:            float64(burst),
		staleAfter:       15 * time.Minute,
		cleanupThreshold: 1024,
	}

	return func(ctx context.Context, c *app.RequestContext) {
		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = "unknown"
		}
		if !limiter.allow(scope + ":" + clientIP) {
			c.Header("Retry-After", "60")
			writeJSONError(c, consts.StatusTooManyRequests, "too many requests")
			return
		}
		c.Next(ctx)
	}
}

func (l *clientRateLimiter) allow(key string) bool {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.buckets) >= l.cleanupThreshold {
		l.cleanup(now)
	}

	bucket, ok := l.buckets[key]
	if !ok {
		l.buckets[key] = &tokenBucket{
			tokens:   l.burst - 1,
			lastSeen: now,
		}
		return true
	}

	elapsed := now.Sub(bucket.lastSeen).Seconds()
	bucket.tokens += elapsed * l.tokensPerSecond
	if bucket.tokens > l.burst {
		bucket.tokens = l.burst
	}
	bucket.lastSeen = now

	if bucket.tokens < 1 {
		return false
	}

	bucket.tokens--
	return true
}

func (l *clientRateLimiter) cleanup(now time.Time) {
	for key, bucket := range l.buckets {
		if now.Sub(bucket.lastSeen) > l.staleAfter {
			delete(l.buckets, key)
		}
	}
}
