package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model

	UserID    uint
	Token     string
	IPAddress string
	UserAgent string
	ExpiresAt time.Time
	CreatedAt time.Time
	User      User
}
