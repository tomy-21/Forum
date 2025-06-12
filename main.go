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
	categoryController := controllers.InitCategoryController(categoryService, topicService, tmpl)
	searchController := controllers.InitSearchController(topicService, tmpl)
	adminController := controllers.InitAdminController(userService, topicService, messageService, tmpl)
	messageController := controllers.InitMessageController(messageService)

	// --- Configuration du Routeur ---
	r := mux.NewRouter()
	r.StrictSlash(true) // Gère les slashs en fin d'URL

	// -- Routes Publiques --
	r.HandleFunc("/", homeController.DisplayHomepage).Methods("GET")
	r.HandleFunc("/topic/{id:[0-9]+}", topicController.ShowTopic).Methods("GET")
	r.HandleFunc("/category/{id:[0-9]+}", categoryController.ShowCategoryPage).Methods("GET")
	r.HandleFunc("/search", searchController.HandleSearch).Methods("GET")
	userController.UserRouter(r) // Gère /login, /register, /logout

	// -- Routes Protégées (pour tous les utilisateurs connectés) --
	authRouter := r.NewRoute().Subrouter()
	authRouter.Use(middleware.AuthMiddleware)
	authRouter.HandleFunc("/profil", userController.ShowProfile).Methods("GET") // <-- AJOUT DE LA ROUTE PROFIL
	authRouter.HandleFunc("/category/{id:[0-9]+}/topics/create", topicController.ShowCreateTopicForm).Methods("GET")
	authRouter.HandleFunc("/category/{id:[0-9]+}/topics/create", topicController.HandleCreateTopic).Methods("POST")
	authRouter.HandleFunc("/topic/{id:[0-9]+}/reply", topicController.PostMessage).Methods("POST")
	authRouter.HandleFunc("/topic/{id:[0-9]+}/delete", topicController.HandleDeleteTopic).Methods("POST")
	authRouter.HandleFunc("/react", reactionController.HandleReaction).Methods("POST")
	authRouter.HandleFunc("/message/{id:[0-9]+}/delete", messageController.HandleDeleteMessage).Methods("POST")

	// -- Routes Administrateur --
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AdminMiddleware)
	adminRouter.HandleFunc("/", adminController.ShowDashboard).Methods("GET")
	adminRouter.HandleFunc("/users/ban/{id:[0-9]+}", adminController.BanUser).Methods("POST")
	adminRouter.HandleFunc("/topics/status/{id:[0-9]+}", adminController.HandleUpdateTopicStatus).Methods("POST")
	adminRouter.HandleFunc("/messages/delete/{id:[0-9]+}", adminController.HandleDeleteMessage).Methods("POST")

	// --- Démarrage du Serveur ---
	log.Println("🚀 Le serveur écoute sur http://localhost:8080")
	// On enveloppe le routeur avec le middleware qui peuple le contexte
	if err := http.ListenAndServe(":8080", middleware.PopulateContextMiddleware(r)); err != nil {
		log.Fatalf("Erreur ListenAndServe : %v", err)
	}
}
