package service

import (
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/repository"
	"gorm.io/gorm"
)

type ProductServiceInterface interface {
	GetProductService(id int) (interface{}, error)
	GetProductsService(f repository.Filter) (int, int, interface{}, error)
	CreateProductService(p *models.Product) (interface{}, error)
	UpdateProductService(p *models.Product, id int) (interface{}, error)
	DeleteProductService(id int) error
}

type productService struct{}

var (
	productRepository repository.ProductRepositoryInterface
)

func NewProductService(repository repository.ProductRepositoryInterface) ProductServiceInterface {
	productRepository = repository
	return &productService{}
}

func (ps *productService) GetProductService(f int) (interface{}, error) {
	return productRepository.GetProductById(f)
}

func (ps *productService) GetProductsService(f repository.Filter) (int, int, interface{}, error) {
	count, err := productRepository.CountFilteredProducts(f.FilterQuery)
	if err != nil || count == 0 {
		switch {
		case count == 0:
			return int(count), 0, nil, gorm.ErrRecordNotFound
		default:
			return int(count), 0, nil, err
		}
	}
	productsCount, products, err1 := productRepository.GetFilteredProducts(f)
	return int(count), productsCount, products, err1
}

func (ps *productService) CreateProductService(p *models.Product) (interface{}, error) {
	return productRepository.CreateProduct(p)
}

func (ps *productService) UpdateProductService(p *models.Product, id int) (interface{}, error) {
	_, err := productRepository.GetProductById(id)
	if err != nil {
		return nil, err
	}
	return productRepository.UpdateProduct(id, p)
}

func (ps *productService) DeleteProductService(id int) error {
	_, err := productRepository.GetProductById(id)
	if err != nil {
		return err
	}
	return productRepository.DeleteProduct(id)
}
