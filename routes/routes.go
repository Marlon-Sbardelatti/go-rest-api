package routes

import (
	"github.com/go-chi/chi/v5"
	"main.go/app"
	"main.go/handlers"
	"main.go/middlewares"
)

func RegisterRoutes(r chi.Router, app *app.App) {
	// Usuário
	r.Route("/user", func(r chi.Router) {
		r.Post("/create", handlers.CreateUserHandler(app))
		r.Post("/login", handlers.LoginUserHandler(app))

		// Sub-rotas com autenticação
		r.With(middlewares.AuthMiddleware).Put("/{id}", handlers.UpdateUserHandler(app))
		r.With(middlewares.AuthMiddleware).Delete("/{id}", handlers.DeleteUserHandler(app))
		r.With(middlewares.AuthMiddleware).Get("/{id}", handlers.GetUserByIdHandler(app))
		r.With(middlewares.AuthMiddleware).Get("/{id}/recipes", handlers.GetUserRecipesHandler(app))
		r.With(middlewares.AuthMiddleware).Get("/", handlers.GetAllUsersHandler(app))
	})

	// Ingrediente
	r.Route("/ingredient", func(r chi.Router) {
		r.Get("/", handlers.GetAllIngredientsHandler(app))
		r.Get("/{id}", handlers.GetIngredientByIdHandler(app))
		r.Get("/name/{name}", handlers.GetIngredientByNameHandler(app))

		// // Sub-rotas com autenticação
		r.With(middlewares.AuthMiddleware).Post("/create", handlers.CreateIngredientHandler(app))
		r.With(middlewares.AuthMiddleware).Put("/{id}", handlers.UpdateIngredientHandler(app))
		r.With(middlewares.AuthMiddleware).Delete("/{id}", handlers.DeleteIngredientHandler(app))
	})

	// Receita
	r.Route("/recipe", func(r chi.Router) {
		r.Get("/", handlers.GetAllRecipesHandler(app))
		r.Get("/{id}", handlers.GetRecipeByIdHandler(app))
		r.Get("/name/{name}", handlers.GetRecipeByNameHandler(app))

		// Sub-rotas com autenticação
		r.With(middlewares.AuthMiddleware).Post("/create", handlers.CreateRecipeHandler(app))
		r.With(middlewares.AuthMiddleware).Put("/{id}", handlers.UpdateRecipeHandler(app))
		r.With(middlewares.AuthMiddleware).Delete("/{id}", handlers.DeleteRecipeHandler(app))

		// Adição e remoção de ingredientes associados à receita
		r.With(middlewares.AuthMiddleware).Post("/ingredients/{id}", handlers.AddIngredientRecipeHandler(app))
		r.With(middlewares.AuthMiddleware).Delete("/ingredients/{id}/{ingredient_id}", handlers.DeleteIngredientRecipeHandler(app))
	})

}
