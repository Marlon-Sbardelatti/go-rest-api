package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"main.go/app"
	"main.go/models"
)

func GetAllRecipesHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipes []models.Recipe

		// Retorna as receitas e ingredientes associados a elas da tabela ingredients_recipes
		result := app.DB.Preload("IngredientsRecipes.Ingredient").Find(&recipes)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying recipes: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Transforma Structs das receitas para JSON
		recipesJson, err := json.Marshal(recipes)
		if err != nil {
			http.Error(w, "Error encoding recipes to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(recipesJson)
	}
}

func GetRecipeByIdHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var recipe models.Recipe

		result := app.DB.Preload("IngredientsRecipes.Ingredient").Where("id = ?", id).First(&recipe)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				http.Error(w, "Recipe not found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying recipe: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		// Transforma Structs da receita para JSON
		recipeJson, err := json.Marshal(recipe)
		if err != nil {
			http.Error(w, "Error encoding recipe to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(recipeJson)
	}
}

func GetRecipeByNameHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func CreateRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recipe models.Recipe

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&recipe)
		if err != nil || recipe.UserID == 0 || recipe.Name == "" || recipe.Instructions == "" {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		result := app.DB.Create(&recipe)
		if result.Error != nil {
			http.Error(w, "Recipe already exists or data is incorrect", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Recipe created!"))
	}
}

func UpdateRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var reqRecipe models.Recipe

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&reqRecipe)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var recipe models.Recipe

		result := app.DB.Where("id = ?", id).First(&recipe)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("Recipe not found")
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying recipe: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		recipe.Name = reqRecipe.Name
		recipe.Instructions = reqRecipe.Instructions
		app.DB.Save(&recipe)

		w.Write([]byte("Recipe updated!"))

	}
}

func DeleteRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var recipe models.Recipe

		result := app.DB.Where("id = ?", id).Delete(&recipe)

		if result.Error != nil {
			fmt.Printf("Error querying recipe: %v\n", result.Error)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			fmt.Println("Recipe not found")
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Write([]byte("Recipe deleted"))
	}
}

func AddIngredientRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Error converting id param to uint", http.StatusInternalServerError)
			return
		}

		var ingredientRecipeReq models.IngredientsRecipes

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&ingredientRecipeReq)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		newRecipe := models.IngredientsRecipes{
			RecipeID:     uint(id),
			IngredientID: ingredientRecipeReq.IngredientID,
			Quantity:     ingredientRecipeReq.Quantity,
		}

		result := app.DB.Create(&newRecipe)
		if result.Error != nil {
			http.Error(w, "Error adding ingredient to recipe", http.StatusBadRequest)
			return
		}

		w.Write([]byte("Ingredient added!"))
		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteIngredientRecipeHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ingredient_id := chi.URLParam(r, "ingredient_id")

		var ingredientRecipe models.IngredientsRecipes
		result := app.DB.Where("recipe_id = ? AND ingredient_id = ?", id, ingredient_id).Delete(&ingredientRecipe)

		if result.Error != nil {
			fmt.Printf("Error querying recipe: %v\n", result.Error)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			fmt.Println("Recipe not found")
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Write([]byte("Ingredient removed from recipe!"))
	}
}
