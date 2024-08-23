package routes

import (
	"github.com/go-chi/chi/v5"
	"main.go/app"
	"main.go/handlers"
)

func RegisterUserRoutes(r chi.Router, app *app.App) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/create", handlers.CreateUserHandler(app))
        r.Post("/login", handlers.LoginUserHandler(app))
	})
}
