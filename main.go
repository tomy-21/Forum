package main

import (
	"Forum/config"
	"Forum/controllers"
	"Forum/middleware"
	"Forum/services"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv()

	db, dbErr := config.InitDB()
	if dbErr != nil {
		log.Fatalf("Erreur lors de l'initialisation de la BDD : %v", dbErr)
	}
	defer db.Close()

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Erreur lors du parsing des templates : %v", err)
	}

	// --- Initialisation des Services ---
	userService := services.NewUserService(db)
	categoryService := services.NewCategoryService(db)
	topicService := services.NewTopicService(db) // <-- Initialiser le TopicService

	// --- Initialisation des Contrôleurs ---
	userController := controllers.InitUserController(userService, tmpl)
	homeController := controllers.InitHomeController(categoryService, topicService, tmpl) // <-- Passer le TopicService
	topicController := controllers.InitTopicController(topicService, tmpl)                // <-- Initialiser le TopicController

	// --- Configuration du Routeur ---
	r := mux.NewRouter()

	// Routes publiques (accessibles par tous)
	r.HandleFunc("/", homeController.DisplayHomepage).Methods("GET")
	userController.UserRouter(r) // Gère /login et /register

	// Routes protégées (nécessitent une connexion)
	// On crée un "sous-routeur" qui utilisera notre middleware
	authRoutes := r.PathPrefix("/topics").Subrouter()
	authRoutes.Use(middleware.AuthMiddleware)
	authRoutes.HandleFunc("/create", topicController.ShowCreateTopicForm).Methods("GET")
	authRoutes.HandleFunc("/create", topicController.HandleCreateTopic).Methods("POST")

	// --- Démarrage du Serveur ---
	log.Println("🚀 Le serveur écoute sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Erreur ListenAndServe : %v", err)
	}
}
