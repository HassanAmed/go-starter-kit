package main

import (
	"log"
	"os"

	"bitbucket.org/mobeen_ashraf1/go-starter-kit/controllers"
	"bitbucket.org/mobeen_ashraf1/go-starter-kit/routers"
	"github.com/joho/godotenv"
)

// Entrypoint
func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	a := controllers.App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	app := routers.InitRoutes(&a)
	app.Router.Run(":4000")
}
