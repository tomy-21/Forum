package controllers

import (
	"Forum/services"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AdminController struct {
	userService    *services.UserService
	topicService   *services.TopicService
	messageService *services.MessageService
	tmpl           *template.Template
}

func InitAdminController(us *services.UserService, ts *services.TopicService, ms *services.MessageService, tmpl *template.Template) *AdminController {
	return &AdminController{
		userService:    us,
		topicService:   ts,
		messageService: ms,
		tmpl:           tmpl,
	}
}

// ShowDashboard affiche la page principale du panneau d'administration.
func (c *AdminController) ShowDashboard(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.GetAllUsers()
	if err != nil {
		http.Error(w, "Impossible de charger les utilisateurs", http.StatusInternalServerError)
		return
	}
	topics, err := c.topicService.GetTopics(1000, 0, 0)
	if err != nil {
		http.Error(w, "Impossible de charger les sujets", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Users":  users,
		"Topics": topics,
	}
	c.tmpl.ExecuteTemplate(w, "admin_dashboard.html", data)
}

// BanUser gère l'action de bannir un utilisateur.
func (c *AdminController) BanUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID utilisateur invalide", http.StatusBadRequest)
		return
	}
	err = c.userService.BanUser(userID)
	if err != nil {
		log.Printf("Erreur lors du bannissement de l'utilisateur %d: %v", userID, err)
		http.Error(w, "Erreur lors du bannissement", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// HandleUpdateTopicStatus gère la mise à jour du statut d'un sujet.
func (c *AdminController) HandleUpdateTopicStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topicID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de sujet invalide", http.StatusBadRequest)
		return
	}
	r.ParseForm()
	newStatus := r.FormValue("status")
	if newStatus != "ouvert" && newStatus != "ferme" && newStatus != "archive" {
		http.Error(w, "Statut invalide", http.StatusBadRequest)
		return
	}
	err = c.topicService.UpdateTopicStatus(topicID, newStatus)
	if err != nil {
		log.Printf("Erreur lors de la mise à jour du statut pour le sujet %d: %v", topicID, err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// HandleDeleteMessage gère la suppression d'un message par un admin.
func (c *AdminController) HandleDeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de message invalide", http.StatusBadRequest)
		return
	}
	err = c.messageService.DeleteMessage(messageID)
	if err != nil {
		log.Printf("Erreur lors de la suppression du message %d: %v", messageID, err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	referer := r.Header.Get("Referer")
	if referer != "" {
		http.Redirect(w, r, referer, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
