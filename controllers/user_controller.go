package controllers

import (
	"Forum/models"
	"Forum/services"
	"html/template"
	"net/http"
	"time"

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
	r.HandleFunc("/logout", c.HandleLogout).Methods("GET")
}

// showRegisterForm affiche le formulaire d'inscription
func (c *UserController) showRegisterForm(w http.ResponseWriter, r *http.Request) {
	// On indique au template qu'on est sur la page d'inscription
	data := map[string]interface{}{
		"Page": "Register",
	}
	c.tmpl.ExecuteTemplate(w, "register.html", data)
}

// handleRegister traite la soumission du formulaire d'inscription
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

// showLoginForm affiche le formulaire de connexion
func (c *UserController) showLoginForm(w http.ResponseWriter, r *http.Request) {
	// On indique au template qu'on est sur la page de connexion
	data := map[string]interface{}{
		"Page": "Login",
	}
	c.tmpl.ExecuteTemplate(w, "login.html", data)
}

// handleLogin traite la soumission du formulaire de connexion
func (c *UserController) handleLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	token, err := c.service.Login(identifier, password)
	if err != nil {
		data := map[string]interface{}{
			"Error": err.Error(),
			"Page":  "Login",
		}
		w.WriteHeader(http.StatusUnauthorized)
		c.tmpl.ExecuteTemplate(w, "login.html", data)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// HandleLogout gère la déconnexion de l'utilisateur.
func (c *UserController) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ShowProfile affiche la page de profil de l'utilisateur connecté.
func (c *UserController) ShowProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := c.service.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Utilisateur non trouvé", http.StatusNotFound)
		return
	}

	topicCount, _ := c.service.GetUserTopicCount(userID)
	messageCount, _ := c.service.GetUserMessageCount(userID)

	data := map[string]interface{}{
		"User":            user,
		"TopicCount":      topicCount,
		"MessageCount":    messageCount,
		"Page":            "Profile",
		"IsAuthenticated": true,
		"CurrentUser":     models.User{ID: user.ID, RoleID: user.RoleID},
	}

	c.tmpl.ExecuteTemplate(w, "profile.html", data)
}
