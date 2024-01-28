package main

import (
	"net/http"

	"github.com/Drinnn/go-expert-api/configs"
	"github.com/Drinnn/go-expert-api/internal/entity"
	"github.com/Drinnn/go-expert-api/internal/infra/database"
	"github.com/Drinnn/go-expert-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Product{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, configs.TokenAuth, configs.JWTExpirationTime)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/", productHandler.GetProducts)
		r.Post("/", productHandler.CreateProduct)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", productHandler.GetProduct)
			r.Put("/", productHandler.UpdateProduct)
			r.Delete("/", productHandler.DeleteProduct)
		})
	})

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/login", userHandler.Login)

	http.ListenAndServe(":3000", r)
}
