package main

import (
	"os"

	controllers "bitbucket.org/mobeen_ashraf1/go-starter-kit/controllers"
	"github.com/joho/godotenv"
)

// Entrypoint
func main() {
	godotenv.Load(".env")
	a := controllers.App{}
	a.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	a.Run(":4000")
}