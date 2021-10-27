package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name     string  `gorm:"size:100;not null;unique" json:"name"`
	Price    float64 `gorm:"type:decimal(12,2)"        json:"price"`
	Category string  `gorm:"size:100;not null"        json:"category"`
}
