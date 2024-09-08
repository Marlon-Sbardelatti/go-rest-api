package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main.go/models"
)

func InitDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Pega informações para acesso ao banco, guardadas em arquivo .env
	host := os.Getenv("HOST")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DATABASE")
	port := os.Getenv("PORT")
	sslmode := os.Getenv("SSLMODE")
	timezone := os.Getenv("TIMEZONE")

	// Gera String no formato aceito para acesso ao Postgres
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timezone)

	// Conecta com o banco de dados
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Ingredient{}, &models.Recipe{}, &models.IngredientsRecipes{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}
