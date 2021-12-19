package main

import (
	"log"
	"os"

	"bitbucket.org/mobeen_ashraf1/go-starter-kit/db"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/routers"
	"github.com/joho/godotenv"
)

// Entrypoint
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	a := routers.App{}
	gormDb := db.NewGormDb()
	gDb :=
		gormDb.GetDb(
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	gormDb.RunMigration()
	a.DB = gDb
	app := routers.InitRoutes(&a)
	if err := app.Engine.Run(":4000"); err != nil {
		log.Fatal("Failed to start server")
	}
}
