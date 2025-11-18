package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
}

// CategorySwagger is used for Swagger documentation
type CategorySwagger struct {
	ID          uint   `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
