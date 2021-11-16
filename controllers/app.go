package controllers

import (
	"fmt"
	"log"
	"os"
	"time"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	Engine *gin.Engine
	DB     *gorm.DB
}

func (a *App) Initialize(host, port, user, password, dbname string) {
	dsnDefault := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable", host, user, password, port)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	var err error
	// Db connection with default
	a.DB, err = gorm.Open(postgres.Open(dsnDefault), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	a.DB.Exec("CREATE DATABASE " + dbname)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color
		},
	)
	// Re Db connection
	a.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal(err)
	}
	if err = a.DB.AutoMigrate(&m.Product{}, &m.Category{}); err != nil {
		log.Fatal(err)
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
