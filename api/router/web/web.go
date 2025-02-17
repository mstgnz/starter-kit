package web

import (
	"github.com/go-chi/chi/v5"
	"github.com/mstgnz/starter-kit/api/handler"
	"github.com/mstgnz/starter-kit/api/infra/basehandler"
	"github.com/mstgnz/starter-kit/api/infra/config"
)

var (
	userHandler = handler.NewUserHandler()
)

func WebRoutes(r chi.Router) {
	/* r.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware)
		r.Get("/verify", config.Catch(userHandler.Verify))
	})
	r.Post("/login", config.Catch(userHandler.Login))
	r.Post("/register", config.Catch(userHandler.Register)) */

	r.Get("/test", config.Catch(basehandler.Handle(userHandler.Login)))
}
