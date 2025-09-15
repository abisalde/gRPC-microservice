package model

import "time"

type Catalog struct {
	ID          string     `json:"id"`
	Name        string     `json:"name" validate:"required,name"`
	Description string     `json:"description" validate:"required,description"`
	Price       float64    `json:"price" validate:"required,price"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type CreateCatalog struct {
	ID          string  `json:"id"`
	Name        string  `json:"name" validate:"required,name"`
	Description string  `json:"description" validate:"required,description"`
	Price       float64 `json:"price" validate:"required,price"`
}
