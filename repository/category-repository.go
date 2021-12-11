package repository

import (
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	GetCategoryById(id int) (interface{}, error)
	CreateCategory(ctg *models.Category) (interface{}, error)
	UpdateCategory(id int, ctg *models.Category) (interface{}, error)
	DeleteCategory(id int) error
}

type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepositoryInterface {
	return &categoryRepository{db}
}

func (ctgRepository *categoryRepository) GetCategoryById(id int) (interface{}, error) {

	type Product struct {
		ID         uint
		Name       string
		Price      float64
		CategoryID uint `json:"-"`
	}
	type Category struct {
		ID       uint
		Name     string
		Products []Product
	}

	ctg := Category{}
	result := ctgRepository.DB.Table("categories").Preload("Products").Find(&ctg, id)
	if ctg.ID == 0 && result.Error == nil {
		result.Error = gorm.ErrRecordNotFound
	}
	return ctg, result.Error

}

func (ctgRepository *categoryRepository) CreateCategory(ctg *models.Category) (interface{}, error) {

	result := ctgRepository.DB.Create(&ctg)
	if result.Error != nil {
		return nil, result.Error
	}

	response := map[string]interface{}{
		"ID":   ctg.ID,
		"name": ctg.Name,
	}

	return response, result.Error
}

func (ctgRepository *categoryRepository) UpdateCategory(id int, ctg *models.Category) (interface{}, error) {

	result := ctgRepository.DB.Model(&ctg).Where("id = ?", id).Update("name", ctg.Name)
	ctg.ID = uint(id)
	response := map[string]interface{}{
		"ID":   ctg.ID,
		"name": ctg.Name,
	}
	return response, result.Error
}

func (ctgRepository *categoryRepository) DeleteCategory(id int) error {

	ctg := models.Category{}
	products := models.Product{}

	txErr := ctgRepository.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.First(&ctg, id).Error; err != nil {
			return err
		}

		if err := tx.Delete(&ctg, id).Error; err != nil {
			return err
		}
		err := tx.Where("category_id = ?", id).Delete(&products).Error
		if err != nil {
			return err
		}
		// commit
		return nil
	})
	return txErr
}
