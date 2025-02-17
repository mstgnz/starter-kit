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
	"github.com/mstgnz/starter-kit/api/internal/config"
	"github.com/mstgnz/starter-kit/api/internal/load"
	"github.com/mstgnz/starter-kit/api/internal/logger"
	"github.com/mstgnz/starter-kit/api/internal/response"
	"github.com/mstgnz/starter-kit/api/internal/validate"
	"github.com/mstgnz/starter-kit/api/middle"
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

/* Handler http isteklerinden soyutla
type Request any // struct yapısında `json:"id", param:"id", query:"id", header:"id"` gibi tagleri kullanılabilir
type Response map[string]any

type HandlerInterface[Req Request, Res Response] interface {
	Handle(ctx context.Context, request Req) (Res, error)
}

func handler[Req Request, Res Response](handler HandlerInterface[Req, Res]) HttpHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var request Req
		// body parser
		if err := json.NewDecoder(r.Body).Decode(&request["body"]); err != nil {
			return err
		}
		// param parser
		request["params"] = chi.Vars(r)
		// query parser
		request["query"] = r.URL.Query()

		// header parser
		request["headers"] = r.Header

		res, err := handler.Handle(r.Context(), request)
		if err != nil {
			return err
		}

		_ = response.WriteJSON(w, http.StatusOK, response.Response{
			Success: true,
			Message: "Success",
			Data:    res,
		})
		return nil
	}
} */

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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	workDir, _ := os.Getwd()
	fileServer(r, "/asset", http.Dir(filepath.Join(workDir, "asset")))

	// swagger
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./view/swagger.html")
	})

	// web without auth
	r.Group(func(r chi.Router) {
		for _, lang := range config.App().Langs {
			r.Get(config.App().Routes["login"][lang], config.Catch(userHandler.Login))
			r.Get(config.App().Routes["register"][lang], config.Catch(userHandler.Register))
		}
		r.Post("/login", config.Catch(userHandler.Login))
		r.Post("/register", config.Catch(userHandler.Register))
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(middle.HeaderMiddleware)
		r.Use(middle.AuthMiddleware)

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
