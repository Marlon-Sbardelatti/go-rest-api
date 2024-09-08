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

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&ingredient)
		if err != nil || ingredient.Name == "" {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		result := app.DB.Create(&ingredient)

		if result.Error != nil {
			http.Error(w, "Ingredient already exists or data is incorrect", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Ingredient created!"))
	}
}

func UpdateIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var reqIngredient models.Ingredient

		// Transforma body da request para uma struct, sem o ID
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&reqIngredient)
		if err != nil || reqIngredient.Name == "" {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Cria struct para armazenar dados do atual ingrediente
		var ingredient models.Ingredient

		result := app.DB.Where("id = ?", id).First(&ingredient)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("Ingredient not found")
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying ingredient: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Struct com o ingrediente recebe nome da struct da request
		ingredient.Name = reqIngredient.Name
		app.DB.Save(&ingredient)

		w.Write([]byte("Ingredient updated!"))
	}
}

func DeleteIngredientHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var ingredient models.Ingredient

		result := app.DB.Where("id = ?", id).Delete(&ingredient)

		if result.Error != nil {
			fmt.Printf("Error querying user: %v\n", result.Error)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			fmt.Println("User not found")
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Write([]byte("Ingredient deleted!"))
	}
}
