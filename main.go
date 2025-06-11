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
		log.Fatalf("Erreur BDD : %v", dbErr)
	}
	defer db.Close()

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Erreur templates : %v", err)
	}

	// --- Initialisation des Services ---
	userService := services.NewUserService(db)
	categoryService := services.NewCategoryService(db)
	topicService := services.NewTopicService(db)
	messageService := services.NewMessageService(db)
	reactionService := services.NewReactionService(db)

	// --- Initialisation des Contrôleurs ---
	userController := controllers.InitUserController(userService, tmpl)
	homeController := controllers.InitHomeController(categoryService, topicService, tmpl)
	topicController := controllers.InitTopicController(topicService, messageService, reactionService, tmpl)
	reactionController := controllers.InitReactionController(reactionService)
	categoryController := controllers.InitCategoryController(categoryService, topicService, tmpl) // <-- NOUVEAU

	// --- Configuration du Routeur ---
	r := mux.NewRouter()

	// -- Routes Publiques --
	r.HandleFunc("/", homeController.DisplayHomepage).Methods("GET")
	r.HandleFunc("/topic/{id:[0-9]+}", topicController.ShowTopic).Methods("GET")
	r.HandleFunc("/category/{id:[0-9]+}", categoryController.ShowCategoryPage).Methods("GET") // <-- NOUVEAU
	userController.UserRouter(r)

	// -- Routes Protégées --
	authRouter := r.NewRoute().Subrouter()
	authRouter.Use(middleware.AuthMiddleware)
	authRouter.HandleFunc("/topics/create", topicController.ShowCreateTopicForm).Methods("GET")
	authRouter.HandleFunc("/topics/create", topicController.HandleCreateTopic).Methods("POST")
	authRouter.HandleFunc("/topic/{id:[0-9]+}/reply", topicController.PostMessage).Methods("POST")
	authRouter.HandleFunc("/topic/{id:[0-9]+}/delete", topicController.HandleDeleteTopic).Methods("POST")
	authRouter.HandleFunc("/react", reactionController.HandleReaction).Methods("POST")

	// --- Démarrage du Serveur ---
	log.Println("🚀 Le serveur écoute sur http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Erreur ListenAndServe : %v", err)
	}
}
