package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"main.go/app"
	"main.go/db"
	"main.go/routes"
)

func main() {
	// Inicializa conexão com banco e cria DAO
	db := db.InitDB()
	app := &app.App{DB: db}

	// Cria o router e registra as rotas do servidor
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	routes.RegisterRoutes(r, app)

	log.Println("Server running on Port 3000")
	http.ListenAndServe(":3000", r)
}
