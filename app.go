package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {

	defaultConnectionString := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable", user, password)
	newConnectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	var err error
	//Db connection with default
	a.DB, err = sql.Open("postgres", defaultConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.DB.Exec("CREATE DATABASE " + dbname)
	//Re Db connection
	a.DB, err = sql.Open("postgres", newConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	// router
	a.Router = mux.NewRouter()
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
	log.Println("App Started at", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))

}

// Get Product Handler
func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Route hit getProduct")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"]) // parseInt

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not Found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

// Create Product Handler
func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Route hit createProduct")
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		fmt.Println("error", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := p.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

// Update Handler
func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Route hit updateProduct")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	p.ID = id

	if err := p.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

// Handler for delete
func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Route hit deleteProduct")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := product{ID: id}
	if err := p.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Response Makers
func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	log.Printf("Errored :\\%d", statusCode)
	respondWithJSON(w, statusCode, map[string]string{"error": msg})
}
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
