package service

import (
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/repository"
)

type CategoryServiceInterface interface {
	GetCategoryService(id int) (interface{}, error)
	CreateCategoryService(ctg *models.Category) (interface{}, error)
	UpdateCategoryService(ctg *models.Category, id int) (interface{}, error)
	DeleteCategoryService(id int) error
}

type categoryService struct{}

var (
	categoryRepository repository.CategoryRepositoryInterface
)

func NewCategoryService(repository repository.CategoryRepositoryInterface) CategoryServiceInterface {
	categoryRepository = repository
	return &categoryService{}
}

func (cs *categoryService) GetCategoryService(id int) (interface{}, error) {
	return categoryRepository.GetCategoryById(id)
}

func (cs *categoryService) CreateCategoryService(ctg *models.Category) (interface{}, error) {
	return categoryRepository.CreateCategory(ctg)
}

func (cs *categoryService) UpdateCategoryService(ctg *models.Category, id int) (interface{}, error) {
	_, err := categoryRepository.GetCategoryById(id)
	if err != nil {
		return nil, err
	}
	return categoryRepository.UpdateCategory(id, ctg)
}

func (cs *categoryService) DeleteCategoryService(id int) error {
	_, err := categoryRepository.GetCategoryById(id)
	if err != nil {
		return err
	}
	return categoryRepository.DeleteCategory(id)
}
