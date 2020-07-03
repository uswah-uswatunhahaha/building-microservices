package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// KeyProduct is a struct
type KeyProduct struct{}

// Products is a struct
type Products struct {
	l         hclog.Logger
	productDB *data.ProductsDB
}

// NewProducts is a constructor
func NewProducts(l hclog.Logger, pdb *data.ProductsDB) *Products {
	return &Products{l, pdb}
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
