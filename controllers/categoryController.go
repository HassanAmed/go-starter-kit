package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryControllerInterface interface {
	GetCategory(c *gin.Context)
	CreateCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

type categoryController struct{}

var (
	categoryService service.CategoryServiceInterface
)

func NewCategoryController(service service.CategoryServiceInterface) CategoryControllerInterface {
	categoryService = service
	return &categoryController{}
}

func (ctgController *categoryController) GetCategory(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("id parameter is not valid")))
		return
	}
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
	ctg, err := categoryService.GetCategoryService(id)
	if err != nil {
		fmt.Println("tyerr", err)
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, errorResponse(errors.New("category not found")))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("unknown error while fetching data from db")))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": ctg})
}

// Create Product Handler
func (ctgController *categoryController) CreateCategory(c *gin.Context) {
	var ctg m.Category
	if err := c.ShouldBindJSON(&ctg); err != nil || ctg.Name == "" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid Payload")))
		return
	}

	result, err := categoryService.CreateCategoryService(&ctg)
	if err != nil {
		switch {
		case IsErrorCode(err, UniqueViolationErrCode):
			c.JSON(http.StatusBadRequest,
				errorResponse(errors.New("category already exists. please use a unique category name")))
		default:
			c.JSON(http.StatusInternalServerError,
				errorResponse(errors.New("unknown error while trying to create category")))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// Update Handler
func (ctgController *categoryController) UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("id parameter is not valid")))
		return
	}
	var ctg m.Category
	if err := c.ShouldBindJSON(&ctg); err != nil || ctg.Name == "" {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid payload")))
		return
	}
	res, err := categoryService.UpdateCategoryService(&ctg, id)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New(fmt.Sprintf("category not found for given id: %v", id))))
		default:
			c.JSON(http.StatusInternalServerError,
				errorResponse(errors.New("unknown error while trying to update category")))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": res})
}

func (ctgController *categoryController) DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("id parameter is not valid")))
		return
	}
	txErr := categoryService.DeleteCategoryService(id)
	if txErr != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(txErr))
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
