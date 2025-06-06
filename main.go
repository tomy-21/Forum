package main

import (
	"html/template"
	"log"
	"net/http"

	"Forum/config"
	"Forum/controllers"
	"Forum/database"
	"Forum/services"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv() // Charge .env

	db, dbErr := config.InitDB() // db est maintenant correctement récupéré
	if dbErr != nil {
		log.Fatal(dbErr.Error())
	}
	defer db.Close()

	// Création des tables après l'initialisation de la DB
	if err := database.CreateTables(db); err != nil { // Assurez-vous que cette fonction existe et est correcte dans votre package database
		log.Fatalf("Error creating tables: %s", err.Error())
	}

	temp, tempErr := template.ParseGlob("./templates/*.html")
	if tempErr != nil {
		log.Fatalf("Erreur chargement templates - %s", tempErr.Error())
	}

	// Initialisez vos services et contrôleurs réels
	userService := services.NewUserService(db)                         // Passez db
	userController := controllers.NewUserController(userService, temp) // Passez service et templates

	router := mux.NewRouter()

	// Configurez vos routes réelles
	// Exemple pour les routes User
	userRouter := router.PathPrefix("/users").Subrouter() // Crée un sous-routeur pour les URL /users/*
	userController.RegisterUserRoutes(userRouter)         // Une méthode dans votre UserController

	log.Println("Serveur en écoute sur http://localhost:8080") // Log avant ListenAndServe
	serveErr := http.ListenAndServe(":8080", router)
	if serveErr != nil {
		log.Fatalf("Erreur lancement serveur - %s", serveErr.Error())
	}
}
