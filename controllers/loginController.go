package controllers

import (
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/service"
	"github.com/gin-gonic/gin"
)

type LoginController interface {
	Login(ctx *gin.Context) string
}

type loginController struct {
	jWtService service.JWTService
}

func NewLoginController(jWtService service.JWTService) LoginController {
	return &loginController{
		jWtService: jWtService,
	}
}

func (controller *loginController) Login(ctx *gin.Context) string {
	authCreds := models.Credentials{}
	err := ctx.ShouldBind(&authCreds)
	if err != nil || authCreds.Username == "" || authCreds.Password == "" {
		return ""
	}

	return controller.jWtService.GenerateToken(authCreds.Username)

}
