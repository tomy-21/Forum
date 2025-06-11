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

// DisplayHomepage gère l'affichage de la page d'accueil avec la pagination.
func (c *HomeController) DisplayHomepage(w http.ResponseWriter, r *http.Request) {
	// --- LOGS DE DÉBOGAGE ---
	log.Println("--- Début du traitement de la page d'accueil ---")

	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		log.Printf("ERREUR: Impossible de récupérer les catégories: %v", err)
		http.Error(w, "Impossible de charger les données du forum", http.StatusInternalServerError)
		return
	}

	const limit = 6

	pageStr := r.URL.Query().Get("page")
	log.Printf("DEBUG: Paramètre 'page' reçu de l'URL : '%s'", pageStr)

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	log.Printf("DEBUG: Page actuelle calculée : %d", page)

	totalTopics, err := c.topicService.GetTotalTopicCount()
	if err != nil {
		log.Printf("ERREUR: Impossible de récupérer le nombre total de sujets: %v", err)
		http.Error(w, "Impossible de charger les données", http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Nombre total de sujets : %d", totalTopics)

	offset := (page - 1) * limit
	log.Printf("DEBUG: Offset calculé pour la requête SQL : %d", offset)

	totalPages := int(math.Ceil(float64(totalTopics) / float64(limit)))

	topics, err := c.topicService.GetTopics(limit, offset)
	if err != nil {
		log.Printf("ERREUR: Impossible de récupérer les sujets pour la page %d: %v", page, err)
		http.Error(w, "Impossible de charger les sujets", http.StatusInternalServerError)
		return
	}
	log.Printf("DEBUG: Nombre de sujets récupérés pour cette page : %d", len(topics))

	pagination := models.PaginationData{
		CurrentPage: page,
		TotalPages:  totalPages,
		HasPrevPage: page > 1,
		HasNextPage: page < totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
	}

	var isAuthenticated bool
	if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
		isAuthenticated = true
	}

	data := map[string]interface{}{
		"Categories":      categories,
		"Topics":          topics,
		"IsAuthenticated": isAuthenticated,
		"Pagination":      pagination,
	}

	log.Println("--- Fin du traitement, envoi du template ---")
	c.tmpl.ExecuteTemplate(w, "index.html", data)
}
