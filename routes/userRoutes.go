package routes

import (
	"github.com/go-chi/chi/v5"
	"main.go/app"
	"main.go/handlers"
	"main.go/middlewares"
)

func RegisterUserRoutes(r chi.Router, app *app.App) {
	r.Route("/user", func(r chi.Router) {
        r.Get("/{id}", handlers.GetUserByIdHandler(app))
		r.Post("/create", handlers.CreateUserHandler(app))
		r.Post("/login", handlers.LoginUserHandler(app))
		r.Post("/login", handlers.LoginUserHandler(app))

		r.With(middlewares.AuthMiddleware).Get("/profile", handlers.UserProfileHandler(app))
	})

}
