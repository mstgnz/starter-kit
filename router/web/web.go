package web

import (
	"github.com/cemilsahin/arabamtaksit/handler"
	"github.com/cemilsahin/arabamtaksit/internal/config"
	"github.com/cemilsahin/arabamtaksit/middle"
	"github.com/go-chi/chi/v5"
)

var (
	userHandler = handler.NewUserHandler()
)

func WebRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware)
		r.Get("/verify", config.Catch(userHandler.Verify))
	})
	r.Post("/login", config.Catch(userHandler.Login))
	r.Post("/register", config.Catch(userHandler.Register))
}
