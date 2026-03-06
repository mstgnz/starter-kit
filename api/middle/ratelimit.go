package middle

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/mstgnz/starter-kit/api/infra/response"
)

// RateLimitConfig holds rate limiter configuration
type RateLimitConfig struct {
	Requests int           // Number of requests allowed
	Window   time.Duration // Time window for the limit
	Message  string        // Custom message for rate limit exceeded
}

// DefaultRateLimitConfig returns default rate limit settings
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Requests: 100,             // 100 requests
		Window:   1 * time.Minute, // per minute
		Message:  "Rate limit exceeded. Please try again later.",
	}
}

// StrictRateLimitConfig returns stricter rate limit for sensitive endpoints
func StrictRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		Requests: 10,              // 10 requests
		Window:   1 * time.Minute, // per minute
		Message:  "Too many requests. Please wait before trying again.",
	}
}

// rateLimitEntry holds the count and expiry for a rate limit key
type rateLimitEntry struct {
	count  int
	expiry time.Time
}

// inMemoryRateLimiter is a simple in-memory rate limiter
type inMemoryRateLimiter struct {
	entries sync.Map
	mu      sync.Mutex
}

// Global rate limiter instance
var limiter = &inMemoryRateLimiter{}

// cleanup removes expired entries periodically
func init() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			now := time.Now()
			limiter.entries.Range(func(key, value any) bool {
				entry := value.(*rateLimitEntry)
				if now.After(entry.expiry) {
					limiter.entries.Delete(key)
				}
				return true
			})
		}
	}()
}

// increment increments the counter for a key and returns the new count
// Uses mutex to prevent race conditions between Load and Store operations
func (l *inMemoryRateLimiter) increment(key string, window time.Duration) (int, time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	expiry := now.Add(window)

	// Try to load existing entry under lock
	if val, ok := l.entries.Load(key); ok {
		entry := val.(*rateLimitEntry)
		if now.Before(entry.expiry) {
			// Entry still valid, increment
			entry.count++
			return entry.count, entry.expiry
		}
		// Entry expired, will be replaced below
	}

	// Create new entry
	entry := &rateLimitEntry{
		count:  1,
		expiry: expiry,
	}
	l.entries.Store(key, entry)
	return 1, expiry
}

// RateLimitMiddleware creates a rate limiting middleware using in-memory storage
func RateLimitMiddleware(cfg RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := GetClientIP(r)

			// Create a unique key for this IP and endpoint
			key := fmt.Sprintf("ratelimit:%s:%s", ip, r.URL.Path)

			// Increment counter
			count, expiry := limiter.increment(key, cfg.Window)

			// Check if limit exceeded
			if count > cfg.Requests {
				ttl := time.Until(expiry)

				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Requests))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", expiry.Unix()))
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))

				_ = response.WriteJSON(w, http.StatusTooManyRequests, response.Response{
					Code:    http.StatusTooManyRequests,
					Success: false,
					Message: cfg.Message,
				})
				return
			}

			// Set rate limit headers
			remaining := cfg.Requests - count
			if remaining < 0 {
				remaining = 0
			}

			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Requests))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

			next.ServeHTTP(w, r)
		})
	}
}

// GlobalRateLimitMiddleware applies a global rate limit per IP (not per endpoint)
func GlobalRateLimitMiddleware(cfg RateLimitConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := GetClientIP(r)

			// Global key for this IP
			key := fmt.Sprintf("ratelimit:global:%s", ip)

			// Increment counter
			count, expiry := limiter.increment(key, cfg.Window)

			// Check if limit exceeded
			if count > cfg.Requests {
				ttl := time.Until(expiry)

				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", cfg.Requests))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(ttl.Seconds())))

				_ = response.WriteJSON(w, http.StatusTooManyRequests, response.Response{
					Code:    http.StatusTooManyRequests,
					Success: false,
					Message: cfg.Message,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
