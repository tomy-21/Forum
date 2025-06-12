package controllers

import (
	"Forum/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MessageController struct {
	messageService *services.MessageService
}

func InitMessageController(ms *services.MessageService) *MessageController {
	return &MessageController{messageService: ms}
}

// HandleDeleteMessage gère la suppression d'un message.
func (c *MessageController) HandleDeleteMessage(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur connecté depuis le contexte
	loggedInUserID, okUserID := r.Context().Value("userID").(int)
	loggedInRoleID, okRoleID := r.Context().Value("roleID").(int)

	if !okUserID || !okRoleID {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}

	// Récupérer l'ID du message depuis l'URL
	vars := mux.Vars(r)
	messageID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de message invalide", http.StatusBadRequest)
		return
	}

	// Récupérer l'ID du propriétaire du message
	ownerID, err := c.messageService.GetMessageOwnerID(messageID)
	if err != nil {
		http.Error(w, "Message non trouvé", http.StatusNotFound)
		return
	}

	// --- LOGIQUE D'AUTORISATION ---
	// L'utilisateur peut supprimer si c'est son message OU s'il est admin
	if loggedInUserID != ownerID && loggedInRoleID != 1 {
		http.Error(w, "Action non autorisée", http.StatusForbidden)
		return
	}

	// Autorisation accordée, on supprime le message
	err = c.messageService.DeleteMessage(messageID)
	if err != nil {
		log.Printf("Erreur lors de la suppression du message %d: %v", messageID, err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}

	// Rediriger vers la page précédente
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
