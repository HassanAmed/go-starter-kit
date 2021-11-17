package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	UniqueViolationErrCode = pq.ErrorCode("23505")
)
const (
	FkeyViolationErrCode = pq.ErrorCode("23503")
)

// Get Product Handler
func (a *App) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if result := paramIsInt(id); !result {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("id parameter is not a valid integer")))
		return
	}
	p := m.Product{}
	err := a.DB.First(&p, id).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, errorResponse(errors.New("product not found")))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("unknown error while trying to fetch data")))
		}
		return
	}
	response := map[string]interface{}{
		"ID":         p.ID,
		"name":       p.Name,
		"price":      p.Price,
		"categoryID": p.CategoryID,
	}
	c.JSON(http.StatusOK, gin.H{"result": response})
}

// Get All Products
func (a *App) GetAllProducts(c *gin.Context) {
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
	sortQuery := fmt.Sprintf("%s %s", sortby, order)
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
	var count int64

	type Product struct {
		ID         uint
		Name       string
		Price      float64
		categoryID uint
		Category   string
	}

	if err := a.DB.Model(&m.Product{}).Where(sql).Count(&count).Error; err != nil || count == 0 {
		switch {
		case count == 0:
			c.JSON(http.StatusNotFound, errorResponse(errors.New("no record matched filter")))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while trying to fetch data")))
		}
		return
	}

	products := []m.Product{}
	if err := a.DB.Table("products").Preload("Category").Where(sql).
		Limit(limit).Offset(offset).Order(sortQuery).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			errorResponse(errors.New("unknown error while trying to fetch data")))

		return
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

	c.JSON(http.StatusOK, gin.H{"totalCount": count, "result": response, "productsCount": len(products)})
}

// Create Product Handler
func (a *App) CreateProduct(c *gin.Context) {
	var p m.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid payload")))
		return
	}

	if err := a.DB.Create(&p).Error; err != nil {
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

	response := map[string]interface{}{
		"ID":         p.ID,
		"name":       p.Name,
		"price":      p.Price,
		"categoryID": p.CategoryID,
	}
	c.JSON(http.StatusCreated, gin.H{"result": response})
}

// Update Handler
func (a *App) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if result := paramIsInt(id); !result {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("id parameter is not a valid integer")))
		return
	}
	var p m.Product
	var ctg m.Category
	if err := c.ShouldBindJSON(&p); err != nil || p.Name == "" || p.Price == 0 || p.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid payload")))
		return
	}
	// verify category exists
	if err := a.DB.First(&ctg, p.CategoryID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("no category exist for given categoryID")))
		return
	}
	// verify product exists
	if err := a.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("no product exist for given id to update")))
		return
	}
	p.Category = ctg
	if err := a.DB.Model(&p).Where("id = ?", id).Updates(
		m.Product{
			Name:       p.Name,
			Price:      p.Price,
			CategoryID: p.CategoryID,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("unknown error while trying to fetch data")))
		return
	}
	const base = 10
	const bitsize = 64
	u64, _ := strconv.ParseUint(id, base, bitsize)
	p.ID = uint(u64)
	response := map[string]interface{}{
		"ID":         p.ID,
		"name":       p.Name,
		"price":      p.Price,
		"categoryID": p.CategoryID,
	}
	c.JSON(http.StatusOK, gin.H{"result": response})
}

// Handler for delete
func (a *App) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if result := paramIsInt(id); !result {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("id parameter is not a valid integer")))
		return
	}
	p := m.Product{}
	// verify product exists
	if err := a.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("no product exist for given id to update")))
		return
	}
	if err := a.DB.Delete(&p, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("unknown error while trying to delete")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
