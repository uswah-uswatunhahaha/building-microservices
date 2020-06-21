package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// KeyProduct is a struct
type KeyProduct struct{}

// Products is a struct
type Products struct {
	l        *log.Logger
	database *sql.DB
}

// NewProducts is a constructor
func NewProducts(l *log.Logger, database *sql.DB) *Products {
	return &Products{l, database}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

func getProductID(r *http.Request) int {
	// parse the product id from url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}

// ErrProductNotFound is a var
var ErrProductNotFound = fmt.Errorf("Product not found")

func (p *Products) findProductID(id int) int {
	prod := &data.Product{}
	p.database.QueryRow("SELECT EXISTS(SELECT * FROM tbl_product WHERE id=?)", id).Scan(&prod.ID)

	if prod.ID == 0 {
		log.Println("ID Not Found")
		return prod.ID
	}

	log.Printf("Found the ID %d", id)

	return prod.ID
}
