package models

import "gorm.io/gorm"

type Skill struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	CategoryID  uint
	Category    Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
