package handlers

import (
	"context"
	"net/http"

	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse

// Edit handles PUT request
func (p *Products) Edit(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	p.l.Debug("[PUT] update records", prod.ID)

	err := p.productDB.UpdateProduct(ctx, *prod, id)

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
		p.l.Error("[ERROR] serializing product", err)
	}

}
