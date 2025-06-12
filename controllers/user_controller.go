package controllers

import (
	"Forum/models"
	"Forum/services"
	"html/template"
	"net/http"
	"time" // Assurez-vous d'importer le package time

	"github.com/gorilla/mux"
)

type UserController struct {
	service *services.UserService
	tmpl    *template.Template
}

func InitUserController(service *services.UserService, tmpl *template.Template) *UserController {
	return &UserController{
		service: service,
		tmpl:    tmpl,
	}
}

// UserRouter enregistre toutes les routes liées à l'utilisateur
func (c *UserController) UserRouter(r *mux.Router) {
	r.HandleFunc("/register", c.showRegisterForm).Methods("GET")
	r.HandleFunc("/register", c.handleRegister).Methods("POST")
	r.HandleFunc("/login", c.showLoginForm).Methods("GET")
	r.HandleFunc("/login", c.handleLogin).Methods("POST")
	r.HandleFunc("/logout", c.HandleLogout).Methods("GET") // <-- Route de déconnexion
}

// Affiche le formulaire d'inscription
func (c *UserController) showRegisterForm(w http.ResponseWriter, r *http.Request) {
	c.tmpl.ExecuteTemplate(w, "register.html", nil)
}

// Traite la soumission du formulaire d'inscription
func (c *UserController) handleRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := models.User{
		Name:     r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	err := c.service.Register(&user)
	if err != nil {
		data := map[string]string{"Error": err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		c.tmpl.ExecuteTemplate(w, "register.html", data)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Affiche le formulaire de connexion
func (c *UserController) showLoginForm(w http.ResponseWriter, r *http.Request) {
	c.tmpl.ExecuteTemplate(w, "login.html", nil)
}

// Traite la soumission du formulaire de connexion
func (c *UserController) handleLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	token, err := c.service.Login(identifier, password)
	if err != nil {
		data := map[string]string{"Error": err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		c.tmpl.ExecuteTemplate(w, "login.html", data)
		return
	}

	// Stocker le token dans un cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Important pour la sécurité
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleLogout gère la déconnexion en supprimant le cookie.
func (c *UserController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Créer un cookie avec une date d'expiration passée pour que le navigateur le supprime
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
		Path:    "/",
	})

	// Rediriger vers la page d'accueil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ... dans controllers/user_controller.go

// ShowProfile affiche la page de profil de l'utilisateur connecté.
func (c *UserController) ShowProfile(w http.ResponseWriter, r *http.Request) {
	// Le middleware a déjà vérifié que l'utilisateur est connecté.
	// On récupère son ID depuis le contexte.
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		// Normalement impossible si le middleware est bien appliqué
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Récupérer les informations de l'utilisateur
	user, err := c.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}

	// Récupérer les statistiques
	topicCount, _ := c.service.GetUserTopicCount(userID)
	messageCount, _ := c.service.GetUserMessageCount(userID)

	data := map[string]interface{}{
		"User":         user,
		"TopicCount":   topicCount,
		"MessageCount": messageCount,
		"Page":         "Profile", // Pour la navigation conditionnelle
		// On passe aussi les infos d'authentification pour le _navbar.html
		"IsAuthenticated": true,
		"CurrentUser":     models.User{ID: user.ID, RoleID: user.RoleID},
	}

	c.tmpl.ExecuteTemplate(w, "profile.html", data)
}
