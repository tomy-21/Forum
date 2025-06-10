package controllers

import (
	"Forum/services"
	"html/template"
	"log"
	"net/http"
)

type HomeController struct {
	categoryService *services.CategoryService
	tmpl            *template.Template
}

func InitHomeController(cs *services.CategoryService, tmpl *template.Template) *HomeController {
	return &HomeController{
		categoryService: cs,
		tmpl:            tmpl,
	}
}

// DisplayHomepage gère l'affichage de la page d'accueil avec la liste des catégories.
func (c *HomeController) DisplayHomepage(w http.ResponseWriter, r *http.Request) {
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		log.Printf("Erreur lors de la récupération des catégories: %v", err)
		http.Error(w, "Impossible de charger les données du forum", http.StatusInternalServerError)
		return
	}

	// Nous passons les catégories au template
	data := map[string]interface{}{
		"Categories": categories,
	}

	c.tmpl.ExecuteTemplate(w, "index.html", data)
}
