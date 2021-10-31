package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name       string   `gorm:"size:100;not null;unique"  json:"name"`
	Price      float64  `gorm:"type:decimal(12,2)"        json:"price"`
	CategoryID uint     `json:"categoryId"`
	Category   Category `gorm:"foreignkey:CategoryID" binding:"required"     json:"-"`
}
