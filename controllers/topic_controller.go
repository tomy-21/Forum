package controllers

import (
	"Forum/models"
	"Forum/services"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductControllers struct {
	service  *services.ProductService
	template *template.Template
}

func InitProductController(service *services.ProductService, template *template.Template) *ProductControllers {
	return &ProductControllers{service: service, template: template}
}

func (c *ProductControllers) ProductRouter(r *mux.Router) {
	r.HandleFunc("/", c.DisplayList).Methods("GET")
	r.HandleFunc("/product/create", c.CreateForm).Methods("GET")
	r.HandleFunc("/product/{id}", c.DisplayById).Methods("GET")
	r.HandleFunc("/product", c.Create).Methods("POST")
}

func (c *ProductControllers) CreateForm(w http.ResponseWriter, r *http.Request) {
	if err := c.template.ExecuteTemplate(w, "product.create", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *ProductControllers) Create(w http.ResponseWriter, r *http.Request) {
	convPrice, convPriceErr := strconv.ParseFloat(r.FormValue("price"), 32)
	convCategorie, convCategorieErr := strconv.Atoi(r.FormValue("categorie"))
	if convPriceErr != nil || convCategorieErr != nil {
		http.Error(w, "Erreur : des données manquantes ou invalides", http.StatusBadRequest)
		return
	}

	newProduct := models.Topic{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Price:       float32(convPrice),
		CategorieId: convCategorie,
	}

	// Validation basique
	if newProduct.Name == "" {
		http.Error(w, "Le nom du produit est obligatoire", http.StatusBadRequest)
		return
	}

	productId, productErr := c.service.Create(newProduct)
	if productErr != nil {
		http.Error(w, productErr.Error(), http.StatusInternalServerError)
		return
	}

	// Redirection après POST → GET
	http.Redirect(w, r, fmt.Sprintf("/product/%d", productId), http.StatusSeeOther)
}

func (c *ProductControllers) DisplayList(w http.ResponseWriter, r *http.Request) {
	productList, productErr := c.service.ReadAll()
	if productErr != nil {
		http.Error(w, productErr.Error(), http.StatusInternalServerError)
		return
	}
	if err := c.template.ExecuteTemplate(w, "product.list", productList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *ProductControllers) DisplayById(w http.ResponseWriter, r *http.Request) {
	idProduct, idProductErr := strconv.Atoi(mux.Vars(r)["id"])
	if idProductErr != nil {
		http.Error(w, "Identifiant produit invalide", http.StatusBadRequest)
		return
	}

	product, productErr := c.service.ReadById(idProduct)
	if productErr != nil {
		http.Error(w, productErr.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.template.ExecuteTemplate(w, "product.detail", product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
