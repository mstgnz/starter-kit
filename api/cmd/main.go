package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mstgnz/starter-kit/api/handler"
	"github.com/mstgnz/starter-kit/api/infra/config"
	"github.com/mstgnz/starter-kit/api/infra/load"
	"github.com/mstgnz/starter-kit/api/infra/logger"
	"github.com/mstgnz/starter-kit/api/infra/response"
	"github.com/mstgnz/starter-kit/api/infra/validate"
	"github.com/mstgnz/starter-kit/api/middle"
	"github.com/mstgnz/starter-kit/api/router/web"
)

var (
	PORT string

	userHandler = handler.NewUserHandler()
	homeHandler = handler.NewHomeHandler()
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

	PORT = os.Getenv("APP_PORT")
}

func main() {

	defer func() {
		config.App().Redis.CloseRedis()
		config.App().Kafka.CloseKafka()
		config.App().DB.CloseDatabase()
	}()

	// Chi Define Routes
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Timestamp", "Hash", "Origin", "X-Requested-With"},
		ExposedHeaders:   []string{"Link", "Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           300, // Preflight cache time (second)
	}))

	// Hash Middleware
	r.Use(middle.HashMiddleware)

	workDir, _ := os.Getwd()
	fileServer(r, "/asset", http.Dir(filepath.Join(workDir, "asset")))

	// swagger
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/swagger.html")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middle.HeaderMiddleware)
		r.Use(middle.AuthMiddleware)
		web.WebRoutes(r)
	})

	// Not Found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		_ = response.WriteJSON(w, http.StatusUnauthorized, response.Response{Success: false, Message: "Not Found"})
		return
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	}
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
