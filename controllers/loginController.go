package controllers

import (
	"errors"
	"net/http"

	"bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
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
	jWtService service.JWTService
}

func NewLoginController(jWtService service.JWTService) LoginController {
	return &loginController{
		jWtService: jWtService,
	}
}

func (controller *loginController) Login(ctx *gin.Context) {
	authCreds := models.User{}
	err := ctx.ShouldBind(&authCreds)
	if err != nil || authCreds.Username == "" || authCreds.Password == "" {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}
	db := GetDB()
	type Result struct {
		Username string
		Password string
	}
	var r Result
	if err := db.Table("users").Where("username = ?", authCreds.Username).Scan(&r).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while Sign up try again")))
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(authCreds.Password)); err != nil {
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
	db := GetDB()
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New(err.Error())))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New(err.Error())))
	}
	u.Password = string(hashedPassword)

	if err := db.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(errors.New("Error while Sign up try again")))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"result": "success"})
}
