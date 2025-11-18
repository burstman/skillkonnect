package models

import (
	"database/sql"

	"gorm.io/gorm"
)

// UserSwagger is used for Swagger documentation
type UserSwagger struct {
	ID              uint    `json:"id"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	DeletedAt       string  `json:"deleted_at,omitempty"`
	Email           string  `json:"email"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Role            string  `json:"role"`
	PasswordHash    string  `json:"password_hash"`
	EmailVerifiedAt string  `json:"email_verified_at"`
	Suspended       bool    `json:"suspended"`
	Approved        bool    `json:"approved"`
	Bio             string  `json:"bio"`
	Rating          float64 `json:"rating"`
}

type User struct {
	gorm.Model
	Email           string
	FirstName       string
	LastName        string
	Role            string `gorm:"default:'client'"` // "admin", "worker", "client"
	PasswordHash    string
	EmailVerifiedAt sql.NullTime
	Suspended       bool    `gorm:"default:false"`
	Approved        bool    `gorm:"default:false"`
	Bio             string  `gorm:"type:text"`
	Rating          float64 `gorm:"default:0"`
}
