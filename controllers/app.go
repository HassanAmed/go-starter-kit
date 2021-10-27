package controllers

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"

	m "bitbucket.org/mobeen_ashraf1/go-starter-kit/models"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
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
	a.DB.Debug().AutoMigrate(&m.Product{})
	// router
	a.Router = mux.NewRouter().StrictSlash(true)
	a.initRoutes()
}

func (a *App) initRoutes() {
	ar := a.Router
	ar.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	ar.HandleFunc("/product", a.createProduct).Methods("POST")
	ar.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	ar.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))

}
