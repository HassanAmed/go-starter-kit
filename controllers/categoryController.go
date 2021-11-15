package controllers

import (
	"errors"
	"net/http"
	"strconv"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (a *App) GetCategory(c *gin.Context) {
	id := c.Param("id")

	type Product struct {
		ID         uint
		Name       string
		Price      float64
		CategoryId uint `json:"-"`
	}
	type Category struct {
		ID       uint
		Name     string
		Products []Product
	}

	ctg := Category{}
	err := a.DB.Table("categories").Preload("Products").Find(&ctg, id).Error
	if err != nil {
		switch err.Error() {
		case "record not found":
			c.JSON(http.StatusNotFound, errorResponse(errors.New("Category not found")))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Unexpected error while fetching data from db")))
		}
		return
	}
	if ctg.Name == "" {
		c.JSON(http.StatusNotFound, errorResponse(errors.New("Category not found")))
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": ctg})
}

// Create Product Handler
func (a *App) CreateCategory(c *gin.Context) {
	var ctg m.Category
	if err := c.ShouldBindJSON(&ctg); err != nil || ctg.Name == "" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid Payload")))
		return
	}

	if err := a.DB.Create(&ctg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error creating record")))
		return
	}
	response := map[string]interface{}{
		"ID":   ctg.ID,
		"name": ctg.Name,
	}

	c.JSON(http.StatusOK, gin.H{"result": response})
}

// Update Handler
func (a *App) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var ctg m.Category
	if err := c.ShouldBindJSON(&ctg); err != nil || ctg.Name == "" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid payload")))
		return
	}
	result := a.DB.Model(&ctg).Where("id = ?", id).Update("name", ctg.Name)
	switch {
	case result.Error != nil:
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while trying to update category")))
		return
	case result.RowsAffected < 1:
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("No rows affected from update")))
		return
	}
	const base = 10
	const bitsize = 64
	u64, _ := strconv.ParseUint(id, base, bitsize)
	ctg.ID = uint(u64)

	response := map[string]interface{}{
		"ID":   ctg.ID,
		"name": ctg.Name,
	}
	c.JSON(http.StatusOK, gin.H{"result": response})
}

func (a *App) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	ctg := m.Category{}
	products := m.Product{}

	txErr := a.DB.Transaction(func(tx *gorm.DB) error {
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

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
		//commit
		return nil
	})
	if txErr != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(txErr))
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
