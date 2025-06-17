package controllers

import (
	"Forum/models"
	"Forum/services"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// TopicController gère toutes les actions liées aux sujets.
type TopicController struct {
	topicService    *services.TopicService
	messageService  *services.MessageService
	reactionService *services.ReactionService
	tmpl            *template.Template
}

// InitTopicController est le constructeur pour le TopicController.
func InitTopicController(ts *services.TopicService, ms *services.MessageService, rs *services.ReactionService, t *template.Template) *TopicController {
	return &TopicController{
		topicService:    ts,
		messageService:  ms,
		reactionService: rs,
		tmpl:            t,
	}
}

// ShowCreateTopicForm affiche le formulaire de création de sujet pour une catégorie donnée.
func (c *TopicController) ShowCreateTopicForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, _ := strconv.Atoi(vars["id"])

	data := map[string]interface{}{
		"CategoryID": categoryID,
	}
	c.tmpl.ExecuteTemplate(w, "create_topic.html", data)
}

// HandleCreateTopic traite la soumission du formulaire de création de sujet.
func (c *TopicController) HandleCreateTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, _ := strconv.Atoi(vars["id"])

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	topic := models.Topic{
		Title:   r.FormValue("title"),
		UserID:  userID,
		ForumID: categoryID,
	}

	topicID, err := c.topicService.Create(&topic)
	if err != nil {
		http.Error(w, "Erreur lors de la création du sujet", http.StatusInternalServerError)
		return
	}

	firstMessage := models.Message{
		TopicID: topicID,
		UserID:  userID,
		Content: r.FormValue("description"),
	}
	err = c.messageService.CreateMessage(&firstMessage)
	if err != nil {
		c.topicService.DeleteTopic(topicID)
		http.Error(w, "Erreur lors de la création du premier message", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/category/"+strconv.Itoa(categoryID), http.StatusSeeOther)
}

// ShowTopic affiche un sujet, ses messages et leurs réactions.
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

	// --- LOGIQUE D'AUTHENTIFICATION MISE À JOUR ---
	var currentUser models.User
	var isAuthenticated bool

	// On vérifie si les valeurs existent bien dans le contexte.
	// C'est plus sûr que de juste vérifier si la conversion fonctionne.
	userID, okUserID := r.Context().Value("userID").(int)
	roleID, okRoleID := r.Context().Value("roleID").(int)

	if okUserID && okRoleID {
		isAuthenticated = true
		currentUser.ID = userID
		currentUser.RoleID = roleID
	}

	data := map[string]interface{}{
		"Topic":           topic,
		"Messages":        messages,
		"IsAuthenticated": isAuthenticated,
		"CurrentUser":     currentUser,
		"Page":            "Topic",
	}

	c.tmpl.ExecuteTemplate(w, "topic.html", data)
}

// PostMessage gère l'envoi d'une réponse dans un sujet.
func (c *TopicController) PostMessage(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value("userID").(int)
	vars := mux.Vars(r)
	topicID, _ := strconv.Atoi(vars["id"])

	// On parse un formulaire multipart (qui peut contenir des fichiers)
	// 32 << 20 signifie une taille max de 32MB pour la requête
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusBadRequest)
		return
	}

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

	// Traitement du fichier uploadé
	file, handler, err := r.FormFile("image")
	if err == nil { // Un fichier a été fourni
		defer file.Close()

		// Générer un nom de fichier unique pour éviter les conflits
		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(handler.Filename))
		filePath := filepath.Join("static", "uploads", fileName)

		// Créer le fichier sur le serveur
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Impossible de créer le fichier sur le serveur", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copier le contenu du fichier uploadé vers le nouveau fichier
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Impossible de sauvegarder le fichier", http.StatusInternalServerError)
			return
		}

		// On enregistre le chemin public de l'image
		msg.ImageURL = sql.NullString{String: "/" + filePath, Valid: true}
	}

	// On crée le message en base de données (avec ou sans image)
	err = c.messageService.CreateMessage(msg)
	if err != nil {
		http.Error(w, "Erreur lors de l'envoi du message", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/topic/"+strconv.Itoa(topicID), http.StatusSeeOther)
}

// HandleDeleteTopic gère la suppression d'un sujet.
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
