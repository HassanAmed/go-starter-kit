package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name       string   `gorm:"size:100;not null;unique"  json:"name" binding:"min=2,max=10"`
	Price      float64  `gorm:"type:decimal(12,2)"       json:"price" binding:"gt=0"`
	CategoryId uint     `json:"categoryId" binding:"required,gt=0"`
	Category   Category `gorm:"foreignkey:CategoryId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"    json:"-"`
}
