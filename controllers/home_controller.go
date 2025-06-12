package controllers

import (
	"Forum/models"
	"Forum/services"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type HomeController struct {
	categoryService *services.CategoryService
	topicService    *services.TopicService
	tmpl            *template.Template
}

func InitHomeController(cs *services.CategoryService, ts *services.TopicService, tmpl *template.Template) *HomeController {
	return &HomeController{
		categoryService: cs,
		topicService:    ts,
		tmpl:            tmpl,
	}
}

// DisplayHomepage gère l'affichage de la page d'accueil.
func (c *HomeController) DisplayHomepage(w http.ResponseWriter, r *http.Request) {
	// Récupérer les catégories
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		log.Printf("ERREUR: Impossible de récupérer les catégories: %v", err)
		http.Error(w, "Impossible de charger les données du forum", http.StatusInternalServerError)
		return
	}

	// Logique de pagination pour les sujets
	const limit = 4 // Vous pouvez ajuster cette valeur
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	totalTopics, err := c.topicService.GetTotalTopicCount(0)
	if err != nil {
		http.Error(w, "Impossible de charger les données", http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(totalTopics) / float64(limit)))

	topics, err := c.topicService.GetTopics(limit, offset, 0)
	if err != nil {
		log.Printf("ERREUR: Impossible de récupérer les sujets pour la page %d: %v", page, err)
		http.Error(w, "Impossible de charger les sujets", http.StatusInternalServerError)
		return
	}

	pagination := models.PaginationData{
		CurrentPage: page,
		TotalPages:  totalPages,
		HasPrevPage: page > 1,
		HasNextPage: page < totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
	}

	// --- LOGIQUE D'AUTHENTIFICATION SIMPLIFIÉE ---
	var isAuthenticated bool
	var currentUser models.User

	// On lit simplement le contexte, qui est mis à jour par le middleware PopulateContextMiddleware.
	userID, okUserID := r.Context().Value("userID").(int)
	roleID, okRoleID := r.Context().Value("roleID").(int)

	if okUserID && okRoleID {
		isAuthenticated = true
		currentUser.ID = userID
		currentUser.RoleID = roleID
		// Le nom d'utilisateur pourrait aussi être passé via le contexte si nécessaire
		currentUser.Name, _ = r.Context().Value("username").(string)
	}

	// Préparation de toutes les données pour la vue
	data := map[string]interface{}{
		"Categories":      categories,
		"Topics":          topics,
		"IsAuthenticated": isAuthenticated,
		"CurrentUser":     currentUser,
		"Pagination":      pagination,
		"Page":            "Home",
	}

	c.tmpl.ExecuteTemplate(w, "index.html", data)
}
