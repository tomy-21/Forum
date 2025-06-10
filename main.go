package main

import (
	"html/template"
	"log"
	"net/http"

	"Forum/config"
	"Forum/controllers"
	"Forum/services"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv()

	db, dbErr := config.InitDB()
	if dbErr != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données : %v", dbErr)
	}
	defer db.Close()

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Erreur lors du parsing des templates : %v", err)
	}

	// --- NOUVEAU CODE ---
	// Initialisation des services
	userService := services.NewUserService(db)
	categoryService := services.NewCategoryService(db) // On crée le service des catégories

	// Initialisation des contrôleurs
	userController := controllers.InitUserController(userService, tmpl)
	homeController := controllers.InitHomeController(categoryService, tmpl) // On crée le contrôleur pour l'accueil

	r := mux.NewRouter()

	// Enregistrement de la route pour la page d'accueil
	r.HandleFunc("/", homeController.DisplayHomepage).Methods("GET")

	// Enregistrement des routes pour les utilisateurs
	userController.UserRouter(r)

	// --- ANCIEN CODE SUPPRIMÉ ---
	// Les lignes concernant ProductService et ProductController ont été retirées.

	log.Println("🚀 Le serveur écoute sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Erreur ListenAndServe : %v", err)
	}
}
