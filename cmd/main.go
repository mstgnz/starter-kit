package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/cemilsahin/arabamtaksit/internal/config"
	"github.com/cemilsahin/arabamtaksit/internal/logger"
	"github.com/cemilsahin/arabamtaksit/internal/response"
	"github.com/cemilsahin/arabamtaksit/internal/validate"
	"github.com/cemilsahin/arabamtaksit/middle"
	"github.com/cemilsahin/arabamtaksit/router/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

var (
	PORT string
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

	PORT = os.Getenv("APP_PORT")
}

type HttpHandler func(w http.ResponseWriter, r *http.Request) error

func main() {

	defer func() {
		config.App().Redis.CloseRedis()
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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	workDir, _ := os.Getwd()
	fileServer(r, "/asset", http.Dir(filepath.Join(workDir, "asset")))

	// handler for web
	r.Group(func(r chi.Router) {
		r.Use(middle.HeaderMiddleware)
		web.WebRoutes(r)
	})

	// Not Found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		_ = response.WriteJSON(w, http.StatusNotFound, response.Response{Success: false, Message: "Not Found"})
	})

	// Create a context that listens for interrupt and terminate signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	// Run your HTTP server in a goroutine
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
		if err != nil && err != http.ErrServerClosed {
			logger.Warn("Fatal Error", "err", err.Error())
			log.Fatal(err.Error())
		}
	}()

	logger.Info("API is running on", PORT)

	// Block until a signal is received
	<-ctx.Done()

	logger.Info("API is shutting on", PORT)

	// set Shutting
	config.App().Shutting = true

	// check Running
	for {
		if config.App().Running <= 0 {
			logger.Info("Cronjobs all done")
			break
		} else {
			logger.Info(fmt.Sprintf("Currently %d active jobs in progress. pending completion...", config.App().Running))
		}
		time.Sleep(time.Second * 5)
	}

	config.App().Cron.Stop()
	logger.Info("Shutting down gracefully...")
	config.App().DB.CloseDatabase()
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
