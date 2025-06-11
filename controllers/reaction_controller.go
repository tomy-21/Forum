package controllers

import (
	"Forum/services"
	"net/http"
	"strconv"
)

type ReactionController struct {
	service *services.ReactionService
}

// InitReactionController est le constructeur pour le ReactionController.
func InitReactionController(s *services.ReactionService) *ReactionController {
	return &ReactionController{service: s}
}

// HandleReaction gère la requête HTTP pour liker/disliker un message.
func (c *ReactionController) HandleReaction(w http.ResponseWriter, r *http.Request) {
	// L'ID de l'utilisateur est récupéré depuis le contexte (placé par le middleware)
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	messageID, _ := strconv.Atoi(r.FormValue("message_id"))
	reactionType := r.FormValue("reaction_type")
	topicID := r.FormValue("topic_id") // On récupère l'ID du sujet pour la redirection

	if messageID == 0 || (reactionType != "like" && reactionType != "dislike") {
		http.Error(w, "Requête invalide", http.StatusBadRequest)
		return
	}

	err := c.service.HandleReaction(userID, messageID, reactionType)
	if err != nil {
		http.Error(w, "Erreur lors du traitement de la réaction", http.StatusInternalServerError)
		return
	}

	// Rediriger l'utilisateur vers la page du sujet d'où il vient
	http.Redirect(w, r, "/topic/"+topicID, http.StatusSeeOther)
}
