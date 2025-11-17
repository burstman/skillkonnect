package models

import (
	"database/sql"

	"gorm.io/gorm"
)

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
