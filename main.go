package main

import (
	"os"

	"github.com/joho/godotenv"
)

// Entrypoint
func main() {
	godotenv.Load(".env")
	a := App{}
	a.Initialize(
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	a.Run(":4000")
}
