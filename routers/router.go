package routers

import (
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/controllers"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/middlewares"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/repository"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

var (
	jwtService      service.JWTServiceInterface
	loginRepository repository.LoginRepositoryInterface
	loginService    service.LoginServiceInterface
	loginController controllers.LoginController
)
var (
	ctgRepository repository.CategoryRepositoryInterface
	ctgService    service.CategoryServiceInterface
	ctgController controllers.CategoryControllerInterface
)

var (
	pRepository repository.ProductRepositoryInterface
	pService    service.ProductServiceInterface
	pController controllers.ProductControllerInterface
)

func InitRoutes(a *App) *App {
	gin.SetMode(gin.ReleaseMode)
	a.Engine = gin.New()
	a.Engine.Use(gin.Logger(), gin.CustomRecovery(middlewares.PanicHandler))
	router := a.Engine
	initLoginService(router, a.DB)
	initCategoryService(router, a.DB)
	initProductService(router, a.DB)
	return a
}

func initLoginService(router *gin.Engine, db *gorm.DB) {
	jwtService = service.NewJWTService()
	loginRepository = repository.NewLoginRepository(db)
	loginService = service.NewLoginService(loginRepository)
	loginController = controllers.NewLoginController(jwtService, loginService)

	// Login Route to get jwt
	router.POST("/signup", loginController.SignUp)
	router.POST("/login", loginController.Login)
}

func initCategoryService(router *gin.Engine, db *gorm.DB) {
	ctgRepository = repository.NewCategoryRepository(db)
	ctgService = service.NewCategoryService(ctgRepository)
	ctgController = controllers.NewCategoryController(ctgService)

	router.GET("/api/category/:id", middlewares.AuthorizeJWT(), ctgController.GetCategory)
	router.POST("/api/category", middlewares.AuthorizeJWT(), ctgController.CreateCategory)
	router.PUT("/api/category/:id", middlewares.AuthorizeJWT(), ctgController.UpdateCategory)
	router.DELETE("/api/category/:id", middlewares.AuthorizeJWT(), ctgController.DeleteCategory)
}

func initProductService(router *gin.Engine, db *gorm.DB) {
	pRepository = repository.NewProductRepository(db)
	pService = service.NewProductService(pRepository)
	pController = controllers.NewProductController(pService)

	router.GET("/api/products", middlewares.AuthorizeJWT(), pController.GetAllProducts)
	router.GET("/api/product/:id", middlewares.AuthorizeJWT(), pController.GetProduct)
	router.POST("/api/product", middlewares.AuthorizeJWT(), pController.CreateProduct)
	router.PUT("/api/product/:id", middlewares.AuthorizeJWT(), pController.UpdateProduct)
	router.DELETE("/api/product/:id", middlewares.AuthorizeJWT(), pController.DeleteProduct)

}
