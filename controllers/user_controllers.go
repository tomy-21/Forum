package controllers

import (
	"html/template"
	"net/http"

	"Forum/models"   // REMPLACEZ forum_app par votre nom de module réel
	"Forum/services" // REMPLACEZ forum_app par votre nom de module réel

	"github.com/gorilla/mux"
)

type UserController struct {
	UserService *services.UserService // Intégrez votre service utilisateur
	Templates   *template.Template
}

// NewUserController crée un nouveau UserController
func NewUserController(us *services.UserService, temp *template.Template) *UserController {
	return &UserController{UserService: us, Templates: temp}
}

// RegisterUserRoutes configure les routes pour les actions liées à l'utilisateur
func (uc *UserController) RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/register", uc.ShowRegistrationForm).Methods("GET")
	router.HandleFunc("/register", uc.HandleRegistration).Methods("POST")
	// Ajoutez les routes de connexion, etc.
	// router.HandleFunc("/login", uc.ShowLoginForm).Methods("GET")
	// router.HandleFunc("/login", uc.HandleLogin).Methods("POST")
}

func (uc *UserController) ShowRegistrationForm(w http.ResponseWriter, r *http.Request) {
	// Assurez-vous que "register.html" existe dans votre dossier ./templates/
	err := uc.Templates.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		http.Error(w, "Failed to render registration page: "+err.Error(), http.StatusInternalServerError)
	}
}

func (uc *UserController) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	registrationPayload := models.RegistrationPayload{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	_, err := uc.UserService.CreateUser(&registrationPayload)
	if err != nil {
		// Vous pourriez vouloir afficher un message plus convivial ou logger l'erreur
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Redirige vers la page de connexion (assurez-vous que cette route existe)
	// Si votre login est sous /users/login, changez la redirection.
	http.Redirect(w, r, "/users/login", http.StatusSeeOther) // Ajusté pour correspondre au préfixe /users
}

// Ajoutez d'autres handlers : showLoginForm, handleLogin etc.
// func (uc *UserController) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
// 	err := uc.Templates.ExecuteTemplate(w, "login.html", nil)
// 	if err != nil {
// 		http.Error(w, "Failed to render login page: "+err.Error(), http.StatusInternalServerError)
// 	}
// }

// func (uc *UserController) HandleLogin(w http.ResponseWriter, r *http.Request) {
//	 // Logique de connexion
// 	w.Write([]byte("Traitement de la connexion"))
// }
