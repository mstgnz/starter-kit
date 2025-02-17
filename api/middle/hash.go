package middle

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mstgnz/starter-kit/api/infra/response"
)

var skipUrls = []string{
	"/swagger",
	"/asset/swagger.yaml",
}

func HashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{
					Success: false,
					Message: "Invalid request",
				})
			}
		}()

		path := r.URL.Path

		// Skip urls
		if slices.Contains(skipUrls, path) {
			next.ServeHTTP(w, r)
			return
		}

		// Header check
		timestamp := r.Header.Get("Timestamp")
		clientHash := r.Header.Get("Hash")
		if timestamp == "" || clientHash == "" {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{
				Success: false,
				Message: "Invalid request",
			})
			return
		}

		// Time check
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil || abs(time.Now().Unix()-ts) > 60 {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{
				Success: false,
				Message: "Invalid request",
			})
			return
		}

		// Hash check
		rawData := "Saha." + timestamp + ":" + strings.TrimPrefix(path, "/api/") + ":" + os.Getenv("APP_SECRET") + ".Kolay"
		serverHash := generateHash(rawData)
		if !hashEquals(serverHash, clientHash) {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{
				Success: false,
				Message: "Invalid request",
			})
			return
		}

		// Request valid
		next.ServeHTTP(w, r)
	})
}

func generateHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func hashEquals(a, b string) bool {
	return subtleConstantTimeCompare(a, b) == 1
}

func subtleConstantTimeCompare(a, b string) int {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b))
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
