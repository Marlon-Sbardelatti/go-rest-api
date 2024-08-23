package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"main.go/app"
	"main.go/db"
	"main.go/routes"
)

func main() {

    // db := db.InitDB()
    db := db.InitDB()
    app := &app.App{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
    
    routes.RegisterUserRoutes(r, app)
	http.ListenAndServe(":3000", r)
}
