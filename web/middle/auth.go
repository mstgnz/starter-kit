package middle

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mstgnz/starter-kit/web/config"
	"github.com/mstgnz/starter-kit/web/model"
	"github.com/mstgnz/starter-kit/web/service/api"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, token := isAuth(r)
		if auth {
			ctx := context.WithValue(r.Context(), config.CKey("user"), config.App().IsAuth)
			ctx = context.WithValue(ctx, config.CKey("token"), token)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		http.Redirect(w, r, config.App().Routes["login"][config.App().Lang], http.StatusSeeOther)
	})
}

func IsAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, token := isAuth(r)
		if auth {
			ctx := context.WithValue(r.Context(), config.CKey("user"), config.App().IsAuth)
			ctx = context.WithValue(ctx, config.CKey("token"), token)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isAuth(r *http.Request) (bool, string) {
	if cookie, err := r.Cookie("Authorization"); err == nil {
		token := strings.Replace(cookie.Value, "Bearer ", "", 1)
		config.App().Token = token
		response, err := api.New().WithToken(token).Get("/verify", nil)

		if err == nil && response.Success {
			if userData, err := json.Marshal(response.Data["user"]); err == nil {
				var user model.User

				if err = json.Unmarshal(userData, &user); err == nil {
					return true, token
				}
			}
		}
	}
	return false, ""
}
