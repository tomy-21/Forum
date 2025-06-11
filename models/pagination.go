package models

// PaginationData contient toutes les informations nécessaires pour afficher
// les boutons "Précédent", "Suivant" et les numéros de page dans une vue.
type PaginationData struct {
	CurrentPage int
	TotalPages  int
	HasPrevPage bool
	HasNextPage bool
	PrevPage    int
	NextPage    int
}
