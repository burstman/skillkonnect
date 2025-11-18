package models

import "gorm.io/gorm"

type Skill struct {
	gorm.Model
	Name        string
	Description string
	CategoryID  uint
	Category    Category
}

// SkillSwagger is used for Swagger documentation
type SkillSwagger struct {
	ID          uint   `json:"id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  uint   `json:"category_id"`
}
