package models

type Topic struct {
	ID          int
	Name        string
	Description string
	Price       float32
	CategorieId int
}
