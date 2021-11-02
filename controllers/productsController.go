package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Get Product Handler
func (a *App) GetProduct(c *gin.Context) {
	id := c.Param("id")
	p := m.Product{}
	if err := a.DB.First(&p, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": p})
}

// Get All Products
func (a *App) GetAllProducts(c *gin.Context) {
	var sql string
	// Filter & Search
	op := c.Query("op")
	price := c.Query("price")
	if price != "" {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Limit query parameter is invalid no"})
			return
		}
	}
	strOffset := c.Query("offset")
	offset := -1
	if strOffset != "" {
		var err error
		offset, err = strconv.Atoi(strOffset)
		if err != nil || offset < -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Offset query parameter is invalid no"})
			return
		}
	}

	products := []m.Product{}
	if err := a.DB.Where(sql).Limit(limit).Offset(offset).Order(sortQuery).Find(&products).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": products, "count": len(products)})
}

// Create Product Handler
func (a *App) CreateProduct(c *gin.Context) {
	var p m.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := a.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": p})
}

// Update Handler
func (a *App) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var p m.Product
	var ctg m.Category
	if err := c.ShouldBindJSON(&p); err != nil || p.Name == "" || p.Price == 0 || p.CategoryID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}
	// verify category exists
	if err := a.DB.First(&ctg, p.CategoryID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No category exists for such categoryId"})
		return
	}
	p.Category = ctg
	if err := a.DB.Model(&p).Where("id = ?", id).Updates(
		m.Product{
			Name:       p.Name,
			Price:      p.Price,
			CategoryID: p.CategoryID,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	const base = 10
	const bitsize = 64
	u64, _ := strconv.ParseUint(id, base, bitsize)
	p.ID = uint(u64)
	c.JSON(http.StatusOK, gin.H{"result": p})
}

// Handler for delete
func (a *App) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	p := m.Product{}
	if err := a.DB.Table("products").Where("id = ?", id).Delete(&p, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
