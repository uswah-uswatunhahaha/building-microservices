package handlers

import (
	"context"
	"net/http"

	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  501: errorResponse

// Create handles POST Request
func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	p.l.Debug("[POST] save records")

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	// err = p.insert(ctx, prod)
	err := p.productDB.AddProduct(ctx, *prod)
	if err != nil {
		p.l.Error("Error", err)
	}

}
