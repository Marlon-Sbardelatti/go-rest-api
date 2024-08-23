package handlers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"main.go/app"
	"main.go/models"
	"net/http"
)

// func CheckPasswordHash(password string, hash string) bool  {
//     err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//     return err == nil
// }

func CreateUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Json inválido", http.StatusBadRequest)
			return
		}

		//encrypt
		hash, _ := hashPassword(user.Password)
		user.Password = hash
		// bytes, _ := bcrypt.GenerateFromPassword([]byte(secret), 10)
		// user.Password = string(bytes)

		result := app.DB.Create(&user)
		if result.Error != nil {
			http.Error(w, "Usuário já existe ou dados incorretos", http.StatusBadRequest)
			return
		}

		log.Println("User created successfully!")
		w.Write([]byte("User created!"))
	}

}

func LoginUserHandler(app *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        //email and password
		var reqUser models.UserLoginRequest

        //transforma em json
		err := json.NewDecoder(r.Body).Decode(&reqUser)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		//verificando existência do usuário
		user, err := getUserByEmail(app, reqUser.Email)
		if err != nil {
			http.Error(w, "Email or password are incorrect", http.StatusUnauthorized)
			return
		}

        //compara a senha inserida com a senha salva no db (hash)
        validUser := checkPasswordHash(reqUser.Password, user.Password)
        
        if !validUser {
			http.Error(w, "Email or password are incorrect", http.StatusUnauthorized)
			return
        }

		w.Header().Set("Content-Type", "application/json")

		userJson, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Write(userJson)
	}
}

//fns privadas
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
