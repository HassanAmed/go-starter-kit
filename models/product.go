package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID       int     `gorm:"primaryKey;unique_index"  json:"id"`
	Name     string  `gorm:"size:100;not null;unique" json:"name"`
	Price    float64 `gorm:"type:decimal(7,6)"        json:"price"`
	Category string  `gorm:"size:100;not null"        json:"category"`
}

func (p *Product) GetProduct(db *gorm.DB) (*Product, error) {
	prod := &Product{}

	if err := db.Debug().Table("products").Where("id = ?", p.ID).First(prod).Error; err != nil {
		return nil, err
	}
	return prod, nil
}

func (p *Product) CreateProduct(db *gorm.DB) (*Product, error) {
	err := db.Debug().Table("products").Create(&p).Error

	if err != nil {
		return &Product{}, err
	}
	return p, nil
}

func (p *Product) UpdateProduct(db *gorm.DB) (*Product, error) {

	if err := db.Debug().Table("products").Where("id = ?", p.ID).Updates(
		Product{
			Name:     p.Name,
			Price:    p.Price,
			Category: p.Category,
		}).Error; err != nil {
		return &Product{}, err
	}
	return p, nil
}

func (p *Product) DeleteProduct(db *gorm.DB) error {
	if err := db.Debug().Table("products").Where("id = ?", p.ID).Delete(&Product{}).Error; err != nil {
		return err
	}
	return nil
}
