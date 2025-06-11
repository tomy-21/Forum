package controllers

import (
	"Forum/models"
	"Forum/services"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CategoryController struct {
	categoryService *services.CategoryService
	topicService    *services.TopicService
	tmpl            *template.Template
}

func InitCategoryController(cs *services.CategoryService, ts *services.TopicService, tmpl *template.Template) *CategoryController {
	return &CategoryController{
		categoryService: cs,
		topicService:    ts,
		tmpl:            tmpl,
	}
}

// ShowCategoryPage affiche la page d'une catégorie avec ses sujets paginés.
func (c *CategoryController) ShowCategoryPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de catégorie invalide", http.StatusBadRequest)
		return
	}

	// Récupérer les infos de la catégorie pour l'en-tête de la page
	category, err := c.categoryService.GetCategoryByID(categoryID)
	if err != nil {
		log.Printf("Erreur GetCategoryByID: %v", err)
		http.Error(w, "Catégorie non trouvée", http.StatusNotFound)
		return
	}

	// Logique de pagination (similaire à celle du HomeController)
	const limit = 10
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	totalTopics, err := c.topicService.GetTotalTopicCount(categoryID)
	if err != nil {
		http.Error(w, "Impossible de charger les données", http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(totalTopics) / float64(limit)))

	topics, err := c.topicService.GetTopics(limit, offset, categoryID)
	if err != nil {
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

	data := map[string]interface{}{
		"Category":   category,
		"Topics":     topics,
		"Pagination": pagination,
	}

	c.tmpl.ExecuteTemplate(w, "category.html", data)
}
