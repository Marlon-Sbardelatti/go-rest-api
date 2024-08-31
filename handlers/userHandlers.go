package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"main.go/app"
	"main.go/models"
)

func CreateUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Encrypt senha
		hash, _ := hashPassword(user.Password)
		user.Password = hash

		result := app.DB.Create(&user)
		if result.Error != nil {
			http.Error(w, "Usuário já existe ou dados incorretos", http.StatusBadRequest)
			return
		}

		log.Println("User created successfully!")
		w.Write([]byte("User created!"))
	}

}

func GetUserByIdHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var user models.User
		result := app.DB.Where("id = ?", id).First(&user)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("User not found")
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying user: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")

		userJson, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusInternalServerError)
			return
		}
		w.Write(userJson)
	}
}

func GetUserRecipesHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		var recipes []models.Recipe
		result := app.DB.Where("user_id = ?", userID).Find(&recipes)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("User not found")
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying user: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")

		if len(recipes) == 0 {
			http.Error(w, "No recipes found for this user", http.StatusNotFound)
			return
		}

		recipesJson, err := json.Marshal(recipes)
		if err != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusInternalServerError)
			return
		}

        w.Write(recipesJson)
	}
}

func DeleteUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var user models.User

		result := app.DB.Where("id = ?", id).Delete(&user)

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

		w.Write([]byte("User removed"))
	}
}

func UpdateUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var reqUser models.User

		// Transforma de JSON para struct
		err := json.NewDecoder(r.Body).Decode(&reqUser)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		var user models.User

		result := app.DB.Where("id = ?", id).First(&user)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("User not found")
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying user: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		user.Username = reqUser.Username
		user.Email = reqUser.Email
		hash, _ := hashPassword(reqUser.Password)
		user.Password = hash
		app.DB.Save(&user)

	}
}

func LoginUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// E-mail e password
		var reqUser models.UserLoginRequest

		// Transforma de JSON para struct
		err := json.NewDecoder(r.Body).Decode(&reqUser)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Verifica existência do usuário
		user, err := getUserByEmail(app, reqUser.Email)
		if err != nil {
			http.Error(w, "Email or password are incorrect", http.StatusUnauthorized)
			return
		}

		// Compara a senha inserida com a senha encriptada salva no db (hash)
		validPsw := checkPasswordHash(reqUser.Password, user.Password)

		if !validPsw {
			http.Error(w, "Email or password are incorrect", http.StatusUnauthorized)
			return
		}

		key := []byte(os.Getenv("SECRET"))

		// Claim (JSON) com as infos que deseja guardar no token
		// Obs.: Informações sensíveis não devem ser armazenadas no token (Ex: psw)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   user.ID,
			"name":  user.Username,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		})

		tokenString, err := token.SignedString(key)
		if err != nil {
			http.Error(w, "Could not create JWT Token", http.StatusInternalServerError)
			return
		}

		tokenString = "Bearer " + tokenString

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", tokenString)

		userJson, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(userJson)
	}
}

func UserProfileHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID")

		fmt.Println(userID)
	}
}

// Funções privadas
func getUserByEmail(app *app.App, email string) (*models.User, error) {
	var user models.User
	result := app.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			fmt.Println("User not found")
		} else {
			fmt.Printf("Error querying user: %v\n", result.Error)
		}
		return nil, result.Error
	}
	return &user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
