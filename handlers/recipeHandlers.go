package handlers

import (
	"net/http"

	"main.go/app"
)

func GetAllRecipesHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetRecipeByIdHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetRecipeByNameHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func CreateRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func UpdateRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func DeleteRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
