package handlers

import (
	"encoding/json"
	"fmt"
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

		// Transforma body da request para uma struct, sem o ID
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&user)
		if err != nil || user.Username == "" || user.Email == "" || user.Password == "" {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Encrypt da senha, convertida para hash
		hash, _ := hashPassword(user.Password)
		user.Password = hash

		result := app.DB.Create(&user)
		if result.Error != nil {
			http.Error(w, "Usuário já existe ou dados incorretos", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created!"))
	}

}

func GetAllUsersHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User

		result := app.DB.Find(&users)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				fmt.Println("Users not found")
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				fmt.Printf("Error querying user: %v\n", result.Error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		userJson, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "Error encoding users to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userJson)
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

		userJson, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(userJson)
	}
}

func GetUserRecipesHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")
		var recipes []models.Recipe

		// Query que seleciona as receitas através do id de usuário associado
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

		// Verifica se foram encontradas linhas com aquele UserID
		if len(recipes) == 0 {
			http.Error(w, "No recipes found for this user", http.StatusNotFound)
			return
		}

		// Transforma o array de structs do tipo Recipe em JSON
		recipesJson, err := json.Marshal(recipes)
		if err != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
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

		w.Write([]byte("User deleted!"))
	}
}

func UpdateUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var reqUser models.User

		// Transforma body da request para uma struct, sem o ID
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&reqUser)

		// Verifica a validez do JSON e se todos os campos necessários foram preenchidos
		if err != nil || reqUser.Username == "" || reqUser.Email == "" || reqUser.Password == "" {
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

		// Atribui à struct do usuário resgatado as novas informações passadas no JSON da request
		user.Username = reqUser.Username
		user.Email = reqUser.Email
		hash, _ := hashPassword(reqUser.Password)
		user.Password = hash
		app.DB.Save(&user)

		w.Write([]byte("User updated!"))
	}
}

func LoginUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Criação de uma struct de usuário com apenas e-mail e password
		var reqUser models.UserLoginRequest

		// Transforma o JSON da requisição para struct do tipo UserLoginRequest
		err := json.NewDecoder(r.Body).Decode(&reqUser)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Verifica existência do usuário no banco, guardando-o na struct, se existir
		user, err := getUserByEmail(app, reqUser.Email)
		if err != nil {
			http.Error(w, "Email or password are incorrect", http.StatusUnauthorized)
			return
		}

		// Compara a senha inserida com a senha encriptada salva no banco (em hash)
		validPsw := checkPasswordHash(reqUser.Password, user.Password)
		if !validPsw {
			http.Error(w, "Email or password are incorrect", http.StatusUnauthorized)
			return
		}

		// Nova chave token é gerada, utilizando informações do usuário como claims
		// Obs.: Informações sensíveis, como a senha, não devem ser armazenadas no token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   user.ID,
			"name":  user.Username,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 72).Unix(), // Tempo de expiração do token
		})

		// Token é assinado utilizando o SECRET
		key := []byte(os.Getenv("SECRET"))
		tokenString, err := token.SignedString(key)
		if err != nil {
			http.Error(w, "Could not create JWT Token", http.StatusInternalServerError)
			return
		}

		// Adiciona o prefixo à String de retorno do token
		tokenString = "Bearer " + tokenString

		// Converte o usuário logado (resgatado do banco e convertido em struct) para JSON
		userJson, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Retorna JSON do usuário e o token de autenticação no Header
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", tokenString)
		w.Write(userJson)
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
