package routers

import (
	"net/http"

	"bitbucket.org/mobeen_ashraf1/go-starter-kit/controllers"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/middlewares"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/service"
	"github.com/gin-gonic/gin"
)

var (
	jwtService      service.JWTService          = service.NewJWTService()
	loginController controllers.LoginController = controllers.NewLoginController(jwtService)
)

func InitRoutes(a *controllers.App) *controllers.App {
	// a := controllers.App{}
	a.Router = gin.Default()
	router := a.Router

	// Login Route to get jwt
	router.POST("/login", func(ctx *gin.Context) {
		token := loginController.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	})

	rg := router.Group("/api", middlewares.AuthorizeJWT())
	{
		rg.GET("/products", a.GetAllProducts)
		rg.GET("/product/:id", a.GetProduct)
		rg.POST("/product", a.CreateProduct)
		rg.PUT("/product/:id", a.UpdateProduct)
		rg.DELETE("/product/:id", a.DeleteProduct)
		rg.GET("/category/:id", a.GetCategory)
		rg.POST("/category", a.CreateCategory)
	}
	return a
}
