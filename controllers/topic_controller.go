package controllers

import (
	"Forum/models"
	"Forum/services"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TopicController struct {
	topicService    *services.TopicService
	messageService  *services.MessageService
	reactionService *services.ReactionService
	tmpl            *template.Template
}

func InitTopicController(ts *services.TopicService, ms *services.MessageService, rs *services.ReactionService, t *template.Template) *TopicController {
	return &TopicController{
		topicService:    ts,
		messageService:  ms,
		reactionService: rs,
		tmpl:            t,
	}
}

func (c *TopicController) ShowCreateTopicForm(w http.ResponseWriter, r *http.Request) {
	c.tmpl.ExecuteTemplate(w, "create_topic.html", nil)
}

func (c *TopicController) HandleCreateTopic(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	r.ParseForm()
	topic := models.Topic{
		Title:  r.FormValue("title"),
		UserID: userID,
	}
	// Créer le sujet
	topicID, err := c.topicService.Create(&topic)
	if err != nil {
		http.Error(w, "Erreur lors de la création du sujet", http.StatusInternalServerError)
		return
	}
	// Créer le premier message (qui est la description)
	firstMessage := models.Message{
		TopicID: topicID,
		UserID:  userID,
		Content: r.FormValue("description"),
	}
	err = c.messageService.CreateMessage(&firstMessage)
	if err != nil {
		// Optionnel : supprimer le sujet si le premier message ne peut être créé
		c.topicService.DeleteTopic(topicID)
		http.Error(w, "Erreur lors de la création du premier message", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *TopicController) ShowTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de sujet invalide", http.StatusBadRequest)
		return
	}
	sortBy := r.URL.Query().Get("sort")
	topic, err := c.topicService.GetTopicByID(id)
	if err != nil {
		log.Printf("Erreur GetTopicByID: %v", err)
		http.Error(w, "Sujet non trouvé", http.StatusNotFound)
		return
	}
	messages, err := c.messageService.GetMessagesByTopicID(id, sortBy)
	if err != nil {
		log.Printf("Erreur GetMessagesByTopicID: %v", err)
		http.Error(w, "Impossible de charger les messages", http.StatusInternalServerError)
		return
	}
	var currentUser models.User
	userID, isAuthenticated := r.Context().Value("userID").(int)
	if isAuthenticated {
		currentUser.ID = userID
		currentUser.RoleID, _ = r.Context().Value("roleID").(int)
	}
	data := map[string]interface{}{
		"Topic":           topic,
		"Messages":        messages,
		"IsAuthenticated": isAuthenticated,
		"CurrentUser":     currentUser,
	}
	c.tmpl.ExecuteTemplate(w, "topic.html", data)
}

func (c *TopicController) PostMessage(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["id"])
	r.ParseForm()
	content := r.FormValue("content")
	if content == "" {
		http.Redirect(w, r, "/topic/"+strconv.Itoa(topicID), http.StatusSeeOther)
		return
	}
	msg := &models.Message{
		TopicID: topicID,
		UserID:  userID,
		Content: content,
	}
	err := c.messageService.CreateMessage(msg)
	if err != nil {
		http.Error(w, "Erreur lors de l'envoi du message", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/topic/"+strconv.Itoa(topicID), http.StatusSeeOther)
}

func (c *TopicController) HandleDeleteTopic(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	roleID := r.Context().Value("roleID").(int)
	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["id"])

	topic, err := c.topicService.GetTopicByID(topicID)
	if err != nil {
		http.Error(w, "Sujet non trouvé", http.StatusNotFound)
		return
	}

	if userID != topic.UserID && roleID != 1 {
		http.Error(w, "Action non autorisée", http.StatusForbidden)
		return
	}

	err = c.topicService.DeleteTopic(topicID)
	if err != nil {
		log.Printf("Erreur lors de la suppression du sujet: %v", err)
		http.Error(w, "Erreur lors de la suppression", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
