package routers

import (
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
	gin.SetMode(gin.ReleaseMode)
	a.Engine = gin.New()
	a.Engine.Use(gin.Logger(), gin.CustomRecovery(middlewares.PanicHandler))
	router := a.Engine

	// Login Route to get jwt
	router.POST("/signup", loginController.SignUp)
	router.POST("/login", loginController.Login)

	rg := router.Group("/api", middlewares.AuthorizeJWT())
	{
		rg.GET("/products", a.GetAllProducts)
		rg.GET("/product/:id", a.GetProduct)
		rg.POST("/product", a.CreateProduct)
		rg.PUT("/product/:id", a.UpdateProduct)
		rg.DELETE("/product/:id", a.DeleteProduct)
		rg.GET("/category/:id", a.GetCategory)
		rg.POST("/category", a.CreateCategory)
		rg.PUT("/category/:id", a.UpdateCategory)
		rg.DELETE("/category/:id", a.DeleteCategory)

	}
	return a
}
