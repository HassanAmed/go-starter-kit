package controllers

import (
	"fmt"
	"log"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

func (a *App) Initialize(host, port, user, password, dbname string) {

	dsnDefault := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable", host, user, password, port)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	var err error
	//Db connection with default
	a.DB, err = gorm.Open(postgres.Open(dsnDefault), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	a.DB.Exec("CREATE DATABASE " + dbname)

	//Re Db connection
	a.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	a.DB.AutoMigrate(&m.Product{}, &m.Category{})
}
