package repository

import (
	"errors"
	"fmt"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ProductRepositoryInterface interface {
	GetProductById(id int) (interface{}, error)
	GetFilteredProducts(f Filter) (int, interface{}, error)
	CountFilteredProducts(w string) (int64, error)
	CreateProduct(prod *m.Product) (interface{}, error)
	UpdateProduct(id int, p *m.Product) (interface{}, error)
	DeleteProduct(id int) error
}

type Filter struct {
	FilterQuery string
	Limit       int
	Offset      int
	Order       string
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepositoryInterface {
	return &productRepository{db}
}

func (prodRepo *productRepository) GetProductById(id int) (interface{}, error) {

	p := m.Product{}
	result := prodRepo.DB.First(&p, id)
	response := map[string]interface{}{
		"ID":         p.ID,
		"name":       p.Name,
		"price":      p.Price,
		"categoryID": p.CategoryID,
	}
	return response, result.Error

}

func (prodRepo *productRepository) CountFilteredProducts(w string) (int64, error) {
	var count int64
	result := prodRepo.DB.Model(&m.Product{}).Where(w).Count(&count)
	return count, result.Error
}

func (prodRepo *productRepository) GetFilteredProducts(f Filter) (int, interface{}, error) {
	products := []m.Product{}
	result := prodRepo.DB.Table("products").Preload("Category").Where(f.FilterQuery).Limit(f.Limit).Offset(f.Offset).Order(f.Order).Find(&products)
	if result.Error != nil {
		return 0, nil, result.Error
	}

	type Product struct {
		ID         uint
		Name       string
		Price      float64
		categoryID uint
		Category   string
	}

	response := make([]Product, len(products))

	for i := 0; i < len(products); i++ {
		response[i] = Product{
			products[i].ID,
			products[i].Name,
			products[i].Price,
			products[i].CategoryID,
			products[i].Category.Name,
		}
	}
	return len(response), response, result.Error
}

func (prodRepo *productRepository) CreateProduct(p *m.Product) (interface{}, error) {

	result := prodRepo.DB.Create(&p)
	fmt.Println("err2", errors.As(result.Error, &pq.Error{Message: "fill in the details.."}))
	if result.Error != nil {
		return nil, result.Error
	}
	response := map[string]interface{}{
		"ID":         p.ID,
		"name":       p.Name,
		"price":      p.Price,
		"categoryID": p.CategoryID,
	}
	return response, result.Error
}

func (prodRepo *productRepository) UpdateProduct(id int, p *m.Product) (interface{}, error) {

	result := prodRepo.DB.Model(&p).Where("id = ?", id).Updates(
		m.Product{
			Name:       p.Name,
			Price:      p.Price,
			CategoryID: p.CategoryID,
		})
	if result.Error != nil {
		return nil, result.Error
	}
	response := map[string]interface{}{
		"ID":         id,
		"name":       p.Name,
		"price":      p.Price,
		"categoryID": p.CategoryID,
	}
	return response, result.Error
}

func (prodRepo *productRepository) DeleteProduct(id int) error {
	p := m.Product{}
	result := prodRepo.DB.Delete(&p, id)
	return result.Error
}
