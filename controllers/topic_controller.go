package controllers

import (
	"Forum/models"
	"Forum/services"
	"html/template"
	"log"
	"net/http"
)

type TopicController struct {
	service *services.TopicService
	tmpl    *template.Template
}

func InitTopicController(s *services.TopicService, t *template.Template) *TopicController {
	return &TopicController{service: s, tmpl: t}
}

// ShowCreateTopicForm affiche le formulaire de création de sujet.
func (c *TopicController) ShowCreateTopicForm(w http.ResponseWriter, r *http.Request) {
	c.tmpl.ExecuteTemplate(w, "create_topic.html", nil)
}

// HandleCreateTopic traite la soumission du formulaire.
func (c *TopicController) HandleCreateTopic(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de l'utilisateur depuis le contexte (placé par le middleware)
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		log.Println("Erreur critique: UserID non trouvé dans le contexte")
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	r.ParseForm()

	topic := models.Topic{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		UserID:      userID, // L'auteur est l'utilisateur connecté !
	}

	_, err := c.service.Create(&topic)
	if err != nil {
		log.Printf("Erreur lors de la création du sujet: %v", err)
		http.Error(w, "Erreur lors de la création du sujet", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
