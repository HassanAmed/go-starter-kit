package routers

import (
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/controllers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(a *controllers.App) *controllers.App {
	// a := controllers.App{}
	a.Router = gin.Default()
	router := a.Router
	router.GET("/product/:id", a.GetProduct)
	router.POST("/product", a.CreateProduct)
	router.PUT("/product/:id", a.UpdateProduct)
	router.DELETE("/product/:id", a.DeleteProduct)

	return a
}
