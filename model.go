package main

import (
	"database/sql"
	"log"
)

type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *product) getProduct(db *sql.DB) error {
	log.Println("retrieving product from DB")
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1", p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) createProduct(db *sql.DB) error {
	log.Println("creating product in DB")
	err := db.QueryRow("INSERT INTO products(name, price) VALUES($1, $2) RETURNING id", p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}
	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	log.Println("updating product in DB")
	_, err := db.Exec("Update products SET name=$1, price=$2 WHERE id=$3", p.Name, p.Price, p.ID)

	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products where id=$1", p.ID)
	log.Println("product deleted")
	return err
}
