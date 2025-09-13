package model

type Catalog struct {
	Name        string `json:"name" validate:"required,name"`
	Description string `json:"description" validate:"required,description"`
	Price       string `json:"price" validate:"required,price"`
}
