package middle

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/mstgnz/starter-kit/api/infra/config"
	"github.com/mstgnz/starter-kit/api/infra/response"
)

func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkMethod := r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH"
		if checkMethod {
			contentType := r.Header.Get("Content-Type")
			// Allow application/json (with optional charset) and multipart/form-data
			isValidContentType := strings.HasPrefix(contentType, "application/json") ||
				strings.HasPrefix(contentType, "multipart/form-data")

			if !isValidContentType {
				_ = response.WriteJSON(w, http.StatusBadRequest, response.Response{
					Success: false,
					Message: "Invalid Content-Type. Expected application/json or multipart/form-data",
				})
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// IPMiddleware sets the client IP in the request context
func IPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := GetClientIP(r)
		ctx := context.WithValue(r.Context(), config.CKey("requestIp"), clientIP)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetClientIP extracts the real client IP from various headers and sources
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (most common proxy header)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" && ip != "unknown" {
				return ip
			}
		}
	}

	// Check X-Real-IP header (common in Nginx)
	if xri := r.Header.Get("X-Real-IP"); xri != "" && xri != "unknown" {
		return xri
	}

	// Check CF-Connecting-IP header (Cloudflare)
	if cfip := r.Header.Get("CF-Connecting-IP"); cfip != "" && cfip != "unknown" {
		return cfip
	}

	// Check X-Client-IP header
	if xcip := r.Header.Get("X-Client-IP"); xcip != "" && xcip != "unknown" {
		return xcip
	}

	// Check X-Cluster-Client-IP header (used by some load balancers)
	if xccip := r.Header.Get("X-Cluster-Client-IP"); xccip != "" && xccip != "unknown" {
		return xccip
	}

	// Fallback to RemoteAddr (direct connection)
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	// If all else fails, return RemoteAddr as is
	return r.RemoteAddr
}
