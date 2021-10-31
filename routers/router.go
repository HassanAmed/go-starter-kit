package routers

import (
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/controllers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(a *controllers.App) *controllers.App {
	// a := controllers.App{}
	a.Router = gin.Default()
	router := a.Router
	router.GET("/products", a.GetAllProducts)
	router.GET("/product/:id", a.GetProduct)
	router.POST("/product", a.CreateProduct)
	router.PUT("/product/:id", a.UpdateProduct)
	router.DELETE("/product/:id", a.DeleteProduct)
	router.GET("/category/:id", a.GetCategory)
	router.POST("/category", a.CreateCategory)
	return a
}
