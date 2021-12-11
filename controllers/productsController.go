package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	r "bitbucket.org/mobeen_ashraf1/go-starter-kit/repository"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	UniqueViolationErrCode = "23505"
)
const (
	FkeyViolationErrCode = "23503"
)

type ProductControllerInterface interface {
	GetProduct(c *gin.Context)
	GetAllProducts(c *gin.Context)
	CreateProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	DeleteProduct(c *gin.Context)
}

type productController struct{}

var (
	pService service.ProductServiceInterface
)

type Filter struct {
	filterQuery string
	limit       int
	offset      int
	order       string
}

func NewProductController(service service.ProductServiceInterface) ProductControllerInterface {
	pService = service
	return &productController{}
}

// Get Product Handler
func (pController *productController) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("id parameter is not valid")))
		return
	}
	p, err := pService.GetProductService(id)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, errorResponse(errors.New("product not found")))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("unknown error while trying to fetch data")))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": p})
}

// Get All Products
func (pController *productController) GetAllProducts(c *gin.Context) {
	var sql string
	// Filter & Search
	op := c.Query("op")
	price := c.Query("price")
	if price != "" {
		if result := paramIsFloat(price); !result {
			c.JSON(http.StatusBadRequest, errorResponse(errors.New("price parameter is not a valid float")))
			return
		}
		switch op {
		case "gt":
			sql = fmt.Sprintf(`price > '%s'`, price)
		case "gte":
			sql = fmt.Sprintf(`price >= '%s'`, price)
		case "lt":
			sql = fmt.Sprintf(`price < '%s'`, price)
		case "lte":
			sql = fmt.Sprintf(`price <= '%s'`, price)
		default:
			sql = fmt.Sprintf(`price = '%s'`, price)
		}
	}
	name := c.Query("name")
	if name != "" && sql != "" {
		name = fmt.Sprintf("%%%s%%", name) // %value%
		sql = fmt.Sprintf(`%s AND name LIKE '%s'`, sql, name)
	}
	// SORTING
	sortby := c.Query("sortby")
	if sortby == "" {
		sortby = "id"
	}
	order := c.Query("order")
	if order == "" {
		order = "asc"
	}
	orderQuery := fmt.Sprintf("%s %s", sortby, order)
	// Limit Offset
	strLimit := c.Query("limit")
	// with a value as -1 for gorms Limit method, we'll get a request without limit as default
	limit := -1
	if strLimit != "" {
		var err error
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < -1 {
			c.JSON(http.StatusBadRequest, errorResponse(errors.New("limit query parameter is invalid num")))
			return
		}
	}
	strOffset := c.Query("offset")
	offset := -1
	if strOffset != "" {
		var err error
		offset, err = strconv.Atoi(strOffset)
		if err != nil || offset < -1 {
			c.JSON(http.StatusBadRequest, errorResponse(errors.New("offset query parameter is invalid num")))
			return
		}
	}

	type Product struct {
		ID         uint
		Name       string
		Price      float64
		categoryID uint
		Category   string
	}
	f := r.Filter{
		FilterQuery: sql,
		Limit:       limit,
		Offset:      offset,
		Order:       orderQuery,
	}
	totalCount, pCount, result, err := pService.GetProductsService(f)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("no matching record found")))
			return
		default:
			c.JSON(http.StatusInternalServerError,
				errorResponse(errors.New("unknown error while trying to fetch category")))
		}
	}
	c.JSON(http.StatusOK, gin.H{"totalCount": totalCount, "result": result, "productsCount": pCount})
}

// Create Product Handler
func (pController *productController) CreateProduct(c *gin.Context) {
	var p m.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid payload")))
		return
	}

	result, err := pService.CreateProductService(&p)
	if err != nil {
		fmt.Println("Here we are with ", err)
		switch {
		case IsErrorCode(err, UniqueViolationErrCode):
			c.JSON(http.StatusBadRequest, errorResponse(errors.New("product already exists. please use a unique product name")))
		case IsErrorCode(err, FkeyViolationErrCode):
			c.JSON(http.StatusBadRequest, errorResponse(errors.New("category with provided id does not exists")))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("unknown error while trying to create product")))
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": result})
}

// Update Handler
func (pController *productController) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("id parameter is not valid")))
		return
	}
	var p m.Product
	// var ctg m.Category
	if err := c.ShouldBindJSON(&p); err != nil || p.Name == "" || p.Price == 0 || p.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid payload")))
		return
	}
	// verify category exists
	if _, err := categoryService.GetCategoryService(int(p.CategoryID)); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New(fmt.Sprintf("no category exist for given categoryID: %v", p.CategoryID))))
		return
	}

	res, err := pService.UpdateProductService(&p, id)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New(fmt.Sprintf("product doesn't exist for given id: %v", id))))
			return
		default:
			c.JSON(http.StatusInternalServerError,
				errorResponse(errors.New("unknown error while trying to update category")))
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": res})
}

// Handler for delete
func (pController *productController) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("id parameter is not valid")))
		return
	}
	// verify product exists
	txErr := pService.DeleteProductService(id)
	if txErr != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(txErr))
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
