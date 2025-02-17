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
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/mstgnz/starter-kit/web/handler"
	"github.com/mstgnz/starter-kit/web/infra/config"
	"github.com/mstgnz/starter-kit/web/infra/localization"
	"github.com/mstgnz/starter-kit/web/infra/validate"
	"github.com/mstgnz/starter-kit/web/middle"
)

var (
	PORT string

	homeHandler = handler.NewHomeHandler()
	userHandler = handler.NewUserHandler()
)

func init() {
	// Load Env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Load Env Error: %v", err)
	}
	// init conf
	_ = config.App()
	validate.CustomValidate()

	// Load Translation
	localization.LoadTranslations()
	//log.Println(localization.Translations["en"]["routes"])

	// Load Routes
	config.LoadRoutesFromJSON()
	//log.Println(config.App().Routes["home"]["tr"])

	PORT = os.Getenv("APP_PORT")
}

func main() {

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

	// without auth
	r.Group(func(r chi.Router) {
		r.Use(middle.IsAuthMiddleware)
		for _, lang := range config.App().Langs {
			r.Get(config.App().Routes["login"][lang], config.Catch(userHandler.Login))
			r.Get(config.App().Routes["register"][lang], config.Catch(userHandler.Register))
		}
		r.Post("/login", config.Catch(userHandler.Login))
		r.Post("/register", config.Catch(userHandler.Register))
	})

	// with auth
	r.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware)
		for _, lang := range config.App().Langs {
			r.Get(config.App().Routes["home"][lang], config.Catch(homeHandler.Home))
		}
	})

	// Not Found
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, config.App().Lang, http.StatusSeeOther)
		}
		http.Redirect(w, r, config.App().Routes["not-found"][config.App().Lang], http.StatusSeeOther)
	})

	// Create a context that listens for interrupt and terminate signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	// Run your HTTP server in a goroutine
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", PORT), r)
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	log.Println("sahakolay is running on", PORT)

	// Block until a signal is received
	<-ctx.Done()

	log.Println("Shutting down gracefully...")
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
