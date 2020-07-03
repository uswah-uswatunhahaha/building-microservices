package handlers

import (
	"net/http"

	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// ListSingle is a handler to get product by id
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)
	cur := r.URL.Query().Get("currency")

	p.l.Debug("[DEBUG] get record id", id)

	// prod, err := p.getProductByID(id)
	prod, err := p.productDB.GetProductByID(id, cur)

	switch err {
	case nil:
	case data.ErrProductNotFound:
		p.l.Error("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never here but log the error just incase
		p.l.Error("[ERROR] serializing product", "error", err)
	}
}

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	p.l.Debug("[DEBUG] get all records")

	cur := r.URL.Query().Get("currency")

	// prod, err := p.getProducts()
	prod, err := p.productDB.GetProducts(cur)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
	}

	err = data.ToJSON(prod, rw)

	if err != nil {
		// we should never here but log the error just incase
		p.l.Error("[ERROR] serializing product", "error ", err)
	}
}
