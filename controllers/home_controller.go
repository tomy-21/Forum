package controllers

import (
	"Forum/services"
	"html/template"
	"log"
	"net/http"
)

type HomeController struct {
	categoryService *services.CategoryService
	topicService    *services.TopicService // <-- AJOUTER LE SERVICE DES SUJETS
	tmpl            *template.Template
}

// MODIFIER LE CONSTRUCTEUR
func InitHomeController(cs *services.CategoryService, ts *services.TopicService, tmpl *template.Template) *HomeController {
	return &HomeController{
		categoryService: cs,
		topicService:    ts, // <-- AJOUTER LE SERVICE DES SUJETS
		tmpl:            tmpl,
	}
}

// DisplayHomepage gère l'affichage de la page d'accueil.
func (c *HomeController) DisplayHomepage(w http.ResponseWriter, r *http.Request) {
	// Récupérer les catégories
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		log.Printf("Erreur lors de la récupération des catégories: %v", err)
		http.Error(w, "Impossible de charger les données du forum", http.StatusInternalServerError)
		return
	}

	// Récupérer les sujets
	topics, err := c.topicService.GetAllTopics() // <-- AJOUTER CET APPEL
	if err != nil {
		log.Printf("Erreur lors de la récupération des sujets: %v", err)
		http.Error(w, "Impossible de charger les données du forum", http.StatusInternalServerError)
		return
	}

	// Combiner toutes les données pour le template
	data := map[string]interface{}{
		"Categories": categories,
		"Topics":     topics, // <-- PASSER LES SUJETS AU TEMPLATE
	}

	c.tmpl.ExecuteTemplate(w, "index.html", data)
}
