package middle

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/mstgnz/starter-kit/api/internal/auth"
	"github.com/mstgnz/starter-kit/api/internal/config"
	"github.com/mstgnz/starter-kit/api/internal/response"
	"github.com/mstgnz/starter-kit/api/repository"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Success: false, Message: "Invalid Token"})
			return
		}
		token = strings.Replace(token, "Bearer ", "", 1)

		userId, err := auth.GetUserIDByToken(token)
		if err != nil {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Success: false, Message: err.Error()})
			return
		}

		user_id, err := strconv.Atoi(userId)
		if err != nil && user_id == 0 {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Success: false, Message: err.Error()})
			return
		}

		userRepository := repository.NewUserRepository()
		user, err := userRepository.GetWithId(r.Context(), user_id)

		if err != nil {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Success: false, Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), config.CKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
