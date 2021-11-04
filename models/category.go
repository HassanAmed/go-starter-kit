package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string    `gorm:"size:100;not null;unique" json:"name"`
	Products []Product `json:"products,omitempty"`
}
