package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE request
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	var prod data.Product

	p.l.Debug("[DELETE] delete records", prod.ID)

	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		p.l.Error("Error", err)
	}

	err = p.productDB.DeleteProduct(ctx, prod, id)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		p.l.Error("Error", err)
	}

}
