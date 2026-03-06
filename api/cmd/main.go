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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/mstgnz/starter-kit/api/handler"
	"github.com/mstgnz/starter-kit/api/infra/config"
	"github.com/mstgnz/starter-kit/api/infra/load"
	"github.com/mstgnz/starter-kit/api/infra/logger"
	"github.com/mstgnz/starter-kit/api/infra/response"
	"github.com/mstgnz/starter-kit/api/infra/validate"
	"github.com/mstgnz/starter-kit/api/middle"
	"github.com/mstgnz/starter-kit/api/router/web"
	"github.com/mstgnz/starter-kit/api/schedule"
)

var (
	PORT string

	userHandler = handler.NewUserHandler()
	homeHandler = handler.NewHomeHandler()
)

func init() {
	// Load Env
	if os.Getenv("APP_ENV") == "local" {
		if err := godotenv.Load(".env"); err != nil {
			logger.Warn(fmt.Sprintf("Load Env Error: %v", err))
		}
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
	if PORT == "" {
		PORT = "8080"
	}
}

func main() {
	args := os.Args[1:]

	if len(args) > 0 {
		HandleCommand(args)
		return
	}

	// Create a context that listens for interrupt and terminate signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	// Scheduler Call in goroutine
	go func() {
		schedule.CallSchedule(ctx, config.App().Cron)
		config.App().Cron.Start()

		// Context dinle ve scheduler'ı kapat
		<-ctx.Done()
		log.Println("Shutting down scheduler...")
		config.App().Cron.Stop()
	}()

	startWebServer(ctx)
}

func startWebServer(ctx context.Context) {

	defer func() {
		config.App().Redis.CloseRedis()
		config.App().Kafka.CloseKafka()
		config.App().DB.CloseDatabase()
	}()

	// init conf
	_ = config.App()

	validate.CustomValidate()

	// Chi Define Routes
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS - Only allow specific origins
	r.Use(middle.CORSMiddleware())

	// Global Rate Limit - 500 requests per minute per IP
	r.Use(middle.GlobalRateLimitMiddleware(middle.RateLimitConfig{
		Requests: 500,
		Window:   time.Minute,
		Message:  "Too many requests. Please slow down.",
	}))

	// IP Middleware - Extract client IP and set in context
	r.Use(middle.IPMiddleware)

	// Hash Middleware
	//r.Use(middle.HashMiddleware)

	workDir, _ := os.Getwd()
	fileServer(r, "/asset", http.Dir(filepath.Join(workDir, "asset")))

	// swagger
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/scalar.html")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middle.HeaderMiddleware)
		r.Use(middle.AuthMiddleware)
		web.WebRoutes(r)
	})

	// Not Found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		_ = response.WriteJSON(w, http.StatusNotFound, response.Response{Success: false, Message: "Not Found"})
	})

	// Run your HTTP server in a goroutine
	go func() {
		server := &http.Server{
			Addr:              fmt.Sprintf(":%s", PORT),
			Handler:           r,
			ReadTimeout:       60 * time.Second,
			WriteTimeout:      60 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 60 * time.Second,
		}
		err := server.ListenAndServe()
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

	logger.Info("Shutting down gracefully...")
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
