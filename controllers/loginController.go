package controllers

import (
	"errors"
	"net/http"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginController interface {
	Login(ctx *gin.Context)
	SignUp(c *gin.Context)
}

type loginController struct {
	jWtService   service.JWTServiceInterface
	loginService service.LoginServiceInterface
}

func NewLoginController(jWtService service.JWTServiceInterface, loginService service.LoginServiceInterface) LoginController {
	return &loginController{
		jWtService:   jWtService,
		loginService: loginService,
	}
}

func (controller *loginController) Login(ctx *gin.Context) {
	authCreds := m.User{}
	err := ctx.ShouldBind(&authCreds)
	if err != nil || authCreds.Username == "" || authCreds.Password == "" {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}
	pHash, srvErr := controller.loginService.GetPasswordService(authCreds.Username)
	if srvErr != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while logging in")))
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(pHash), []byte(authCreds.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("Incorrect Credentials")))
		return
	}

	token := controller.jWtService.GenerateToken(authCreds.Username)
	if token != "" {
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
		return
	}
	ctx.JSON(http.StatusInternalServerError, nil)
}

func (controller *loginController) SignUp(c *gin.Context) {
	var u m.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("Invalid Payload")))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New(err.Error())))
	}
	u.Password = string(hashedPassword)
	_, srvErr := controller.loginService.CreateUserService(u)
	if srvErr != nil {
		switch {
		case IsErrorCode(err, UniqueViolationErrCode):
			c.JSON(http.StatusBadRequest, errorResponse(errors.New("Already Registered")))
			return
		default:
			c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while Signing up try again")))
			return
		}

	}
	c.JSON(http.StatusCreated, gin.H{"result": "success"})
}
