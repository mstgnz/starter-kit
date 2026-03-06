package middle

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/cors"
)

// AllowedOrigins returns the list of allowed origins
func AllowedOrigins() []string {
	origins := []string{
		"https://flowize.app",
		"https://*.flowize.app",
		"https://getflowize.com",
		"https://*.getflowize.com",
	}

	// Add development origins if in development mode
	env := os.Getenv("APP_ENV")
	if env == "development" || env == "local" || env == "" {
		origins = append(origins,
			"http://localhost:3000",
			"http://localhost:3600",
			"http://localhost:4200",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3600",
			"http://127.0.0.1:4200",
		)
	}

	return origins
}

// isOriginAllowed checks if the origin matches allowed patterns
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		// Handle wildcard subdomain matching
		if after, ok := strings.CutPrefix(allowed, "https://*."); ok {
			// Extract the base domain (e.g., "flowize.app" from "https://*.flowize.app")
			baseDomain := after

			// Check if origin ends with the base domain and starts with https://
			if strings.HasPrefix(origin, "https://") {
				originDomain := strings.TrimPrefix(origin, "https://")
				// Must be a subdomain (contain a dot before the base domain) or exact match
				if originDomain == baseDomain || strings.HasSuffix(originDomain, "."+baseDomain) {
					return true
				}
			}
		} else if origin == allowed {
			return true
		}
	}
	return false
}

// CORSMiddleware returns configured CORS handler
func CORSMiddleware() func(http.Handler) http.Handler {
	allowedOrigins := AllowedOrigins()

	return cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			// Empty origin is allowed for same-origin requests
			if origin == "" {
				return true
			}
			return isOriginAllowed(origin, allowedOrigins)
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodPatch,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"Timestamp",
			"Hash",
			"Origin",
			"X-Requested-With",
			"X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"Link",
			"Content-Length",
			"Access-Control-Allow-Origin",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
			"Retry-After",
		},
		// AllowCredentials is safe here because AllowOriginFunc performs explicit
		// origin validation (exact match or *.flowize.app subdomain only).
		// This is NOT a wildcard "*" — browsers enforce CORS origin checks.
		AllowCredentials: true,
		MaxAge:           300, // Preflight cache time (seconds)
	})
}
