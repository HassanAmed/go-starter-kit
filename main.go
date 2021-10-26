package main

import (
	"log"
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
	func() {
		if _, err := a.DB.Exec(tableCreationQuery); err != nil {
			log.Fatal(err)
		}
	}()
	a.Run(":4000")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)`
