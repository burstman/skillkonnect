package models

import "gorm.io/gorm"

// WorkerProfileSwagger is used for Swagger documentation
type WorkerProfileSwagger struct {
	ID         uint    `json:"id"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	DeletedAt  string  `json:"deleted_at,omitempty"`
	UserID     uint    `json:"user_id"`
	Name       string  `json:"name"`
	Profession string  `json:"profession"`
	Rating     float64 `json:"rating"`
	Distance   float64 `json:"distance"` // in km
	Reviews    int     `json:"reviews"`
	Price      float64 `json:"price"` // price per hour
	Available  bool    `json:"available"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type WorkerProfile struct {
	gorm.Model
	UserID     uint    `gorm:"uniqueIndex;not null"`
	User       User    `gorm:"foreignKey:UserID"`
	Name       string  `gorm:"not null"`
	Profession string  `gorm:"not null"`
	Rating     float64 `gorm:"default:0"`
	Distance   float64 `gorm:"default:0"` // in km (calculated field, can be dynamically computed)
	Reviews    int     `gorm:"default:0"`
	Price      float64 `gorm:"default:0"` // price per hour
	Available  bool    `gorm:"default:true"`
	Latitude   float64 `gorm:"default:0"`
	Longitude  float64 `gorm:"default:0"`
}
