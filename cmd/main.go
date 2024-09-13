package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mstgnz/starter-kit/handler"
	"github.com/mstgnz/starter-kit/internal/auth"
	"github.com/mstgnz/starter-kit/internal/config"
	"github.com/mstgnz/starter-kit/internal/load"
	"github.com/mstgnz/starter-kit/internal/localization"
	"github.com/mstgnz/starter-kit/internal/logger"
	"github.com/mstgnz/starter-kit/internal/response"
	"github.com/mstgnz/starter-kit/internal/validate"
	"github.com/mstgnz/starter-kit/model"
)

var (
	PORT string

	userHandler = handler.UserHandler{}
	homeHandler = handler.HomeHandler{}
)

func init() {
	// Load Env
	if err := godotenv.Load(".env"); err != nil {
		logger.Warn(fmt.Sprintf("Load Env Error: %v", err))
		log.Fatalf("Load Env Error: %v", err)
	}
	// init conf
	_ = config.App()
	validate.CustomValidate()

	// Load Sql
	config.App().QUERY = make(map[string]string)
	if query, err := load.LoadSQLQueries(); err != nil {
		logger.Warn(fmt.Sprintf("Load Sql Error: %v", err))
		log.Fatalf("Load Sql Error: %v", err)
	} else {
		config.App().QUERY = query
	}

	// Load Translation
	localization.LoadTranslations()
	//log.Println(localization.Translations["en"]["routes"])

	// Load Routes
	config.LoadRoutesFromJSON()
	//log.Println(config.App().Routes["home"]["tr"])

	PORT = os.Getenv("APP_PORT")
}

type HttpHandler func(w http.ResponseWriter, r *http.Request) error

func Catch(h HttpHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			logger.Info("HTTP Handler Error", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

func main() {

	defer func() {
		config.App().Redis.CloseRedis()
		config.App().Kafka.CloseKafka()
		config.App().DB.CloseDatabase()
	}()

	// Chi Define Routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	workDir, _ := os.Getwd()
	fileServer(r, "/asset", http.Dir(filepath.Join(workDir, "asset")))

	// swagger
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/swagger.html")
	})

	// web without auth
	r.Group(func(r chi.Router) {
		r.Use(isAuthMiddleware)
		for _, lang := range config.App().Langs {
			r.Get(config.App().Routes["login"][lang], Catch(userHandler.LoginHandler))
			r.Get(config.App().Routes["register"][lang], Catch(userHandler.RegisterHandler))
		}
		r.Post("/login", Catch(userHandler.LoginHandler))
		r.Post("/register", Catch(userHandler.RegisterHandler))
	})

	// web with auth
	r.Group(func(r chi.Router) {
		r.Use(webAuthMiddleware)
		for _, lang := range config.App().Langs {
			r.Get(config.App().Routes["home"][lang], Catch(homeHandler.HomeHandler))
		}
	})

	// api without auth
	r.With(headerMiddleware).Post("/api/login", Catch(userHandler.LoginHandler))
	r.With(headerMiddleware).Post("/api/register", Catch(userHandler.RegisterHandler))

	r.Route("/api", func(r chi.Router) {
		r.Use(headerMiddleware)
		r.Use(apiAuthMiddleware)

	})

	// Not Found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "api") {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: "Not Found"})
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}
}

func isAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")

		if err == nil {
			token := strings.Replace(cookie.Value, "Bearer ", "", 1)
			_, err = auth.GetUserIDByToken(token)
			if err == nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func webAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Authorization")

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		token := strings.Replace(cookie.Value, "Bearer ", "", 1)

		userId, err := auth.GetUserIDByToken(token)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user_id, err := strconv.Atoi(userId)
		if err != nil && user_id == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user := &model.User{}
		err = user.GetWithId(user_id)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), config.CKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func apiAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: "Invalid Token"})
			return
		}
		token = strings.Replace(token, "Bearer ", "", 1)

		userId, err := auth.GetUserIDByToken(token)
		if err != nil {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: err.Error()})
			return
		}

		user_id, err := strconv.Atoi(userId)
		if err != nil && user_id == 0 {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: err.Error()})
			return
		}

		user := &model.User{}
		err = user.GetWithId(user_id)

		if err != nil {
			_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Status: false, Message: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), config.CKey("user"), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		checkMethod := r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH"
		if checkMethod && r.Header.Get("Content-Type") != "application/json" {
			_ = response.WriteJSON(w, http.StatusBadRequest, response.Response{Status: false, Message: "Invalid Content-Type"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
