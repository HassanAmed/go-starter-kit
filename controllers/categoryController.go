package controllers

import (
	"errors"
	"net/http"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func (a *App) GetCategory(c *gin.Context) {
	id := c.Param("id")

	ctg := m.Category{}
	err := a.DB.Preload("Products").Find(&ctg, id).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
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
