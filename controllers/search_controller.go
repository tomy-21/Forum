package controllers

import (
	"Forum/models"
	"Forum/services"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type SearchController struct {
	topicService *services.TopicService
	tmpl         *template.Template
}

func InitSearchController(ts *services.TopicService, tmpl *template.Template) *SearchController {
	return &SearchController{
		topicService: ts,
		tmpl:         tmpl,
	}
}

// HandleSearch gère l'affichage de la page de résultats de recherche.
func (c *SearchController) HandleSearch(w http.ResponseWriter, r *http.Request) {
	// Récupérer le terme de recherche depuis l'URL (ex: /search?q=mon-sujet)
	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Logique de pagination pour les résultats de recherche
	const limit = 10
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	totalTopics, err := c.topicService.GetSearchTopicCount(searchQuery)
	if err != nil {
		http.Error(w, "Impossible de traiter la recherche", http.StatusInternalServerError)
		return
	}

	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(totalTopics) / float64(limit)))

	topics, err := c.topicService.SearchTopics(searchQuery, limit, offset)
	if err != nil {
		http.Error(w, "Impossible de récupérer les résultats de la recherche", http.StatusInternalServerError)
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
		"Query":      searchQuery,
		"Topics":     topics,
		"Pagination": pagination,
	}

	c.tmpl.ExecuteTemplate(w, "search_results.html", data)
}
