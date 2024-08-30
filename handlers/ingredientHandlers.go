package handlers

import (
	"net/http"

	"main.go/app"
)

func GetAllIngredientsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetIngredientByIdHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetIngredientByNameHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func CreateIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func UpdateIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func DeleteIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
