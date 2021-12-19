package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	. "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseEngineInterface interface {
	GetDb(host, port, user, password, dbname string) *gorm.DB
	RunMigration()
}

type gormDb struct {
	db   *gorm.DB
	once sync.Once
}

func NewGormDb() DatabaseEngineInterface {
	return &gormDb{}
}

func (g *gormDb) GetDb(host, port, user, password, dbname string) *gorm.DB {
	if g.db == nil {
		g.once.Do(func() {
			InitDatabase(host, port, user, password, dbname, g)
		})
	}
	return g.db
}

func (g *gormDb) RunMigration() {
	if g.db == nil {
		panic("Initialise gorm db before running migrations")
	}
	err := g.db.AutoMigrate(&Product{}, &Category{}, &User{})
	if err != nil {
		log.Fatal(err)
	}
}

func InitDatabase(host, port, user, password, dbname string, g *gormDb) {
	dsnDefault := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable", host, user, password, port)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	// Db connection with default
	dbConn, err := gorm.Open(postgres.Open(dsnDefault), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	dbConn.Exec("CREATE DATABASE " + dbname)
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
	dbConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatal(err)
	}
	g.db = dbConn
}
