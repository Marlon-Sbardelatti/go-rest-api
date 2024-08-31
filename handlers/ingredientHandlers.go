package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"main.go/app"
	"main.go/models"
)

func GetAllIngredientsHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ingredients []models.Ingredient

		result := app.DB.Find(&ingredients)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying ingredients: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		if len(ingredients) == 0 {
			http.Error(w, "No ingredients found", http.StatusNotFound)
			return
		}

		ingredientsJson, err := json.Marshal(ingredients)

		if err != nil {
			http.Error(w, "Error encoding ingredients to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(ingredientsJson)
	}
}

func GetIngredientByIdHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var ingredient models.Ingredient

		result := app.DB.Where("id = ?", id).First(&ingredient)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Ingredient not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying ingredients: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		ingredientJson, err := json.Marshal(ingredient)

		if err != nil {
			http.Error(w, "Error encoding ingredient to JSON", http.StatusInternalServerError)
			return
		}

		log.Println(ingredientJson)
		w.Header().Set("Content-Type", "application/json")
		w.Write(ingredientJson)

	}
}

func GetIngredientByNameHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		name = "%" + name + "%"

		var ingredient []models.Ingredient

		result := app.DB.Where("LOWER(name) LIKE ?", strings.ToLower(name)).Find(&ingredient)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Ingredient not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying ingredients: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		ingredientJson, err := json.Marshal(ingredient)

		if err != nil {
			http.Error(w, "Error encoding ingredient to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(ingredientJson)

	}
}

func CreateIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ingredient models.Ingredient

		err := json.NewDecoder(r.Body).Decode(&ingredient)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

        // validate := Validator.new()

        if ingredient.Name == "" {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
        }

        result := app.DB.Create(&ingredient)

		if result.Error != nil {
			http.Error(w, "Ingredient already exists or data is incorrect", http.StatusBadRequest)
			return
		}

		w.Write([]byte("Ingredient created!"))
	}
}

func UpdateIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func DeleteIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
