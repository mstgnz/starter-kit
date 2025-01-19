package middle

import (
	"net/http"

	"github.com/cemilsahin/arabamtaksit/internal/response"
)

func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkMethod := r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH"
		if checkMethod && r.Header.Get("Content-Type") != "application/json" {
			_ = response.WriteJSON(w, http.StatusBadRequest, response.Response{Success: false, Message: "Invalid Content-Type"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
