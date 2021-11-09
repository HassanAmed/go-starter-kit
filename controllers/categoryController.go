package controllers

import (
	"errors"
	"net/http"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"github.com/gin-gonic/gin"
)

func (a *App) GetCategory(c *gin.Context) {
	id := c.Param("id")

	ctg := m.Category{}
	err := a.DB.Preload("Products").Find(&ctg, id).Error
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

	c.JSON(http.StatusOK, gin.H{"result": ctg})
}

func (a *App) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	ctg := m.Category{}
	products := m.Product{}

	if err := a.DB.First(&ctg, id).Error; err != nil {
		switch err.Error() {
		case "record not found":
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Category does not exist.")))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while fetching category to delete.")))
			return
		}
	}
	if err := a.DB.Delete(&ctg, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while deleting category.")))
		return
	}
	if err := a.DB.Where("category_id = ?", id).Delete(&products).Error; err != nil {
		switch err.Error() {
		case "record not found":
			c.JSON(http.StatusOK, gin.H{"result": "success"})
			return
		default:
			if err := a.DB.Unscoped().Model(&ctg).Where("id = ?", id).Updates(map[string]interface{}{"deleted_at": "NULL"}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while deleting associated products of category")))
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while deleting category Tx Reverted")))
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
