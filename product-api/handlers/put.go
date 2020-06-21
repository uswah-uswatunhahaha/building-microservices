package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// Edit handles PUT request
func (p *Products) Edit(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("[PUT] update records")
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err := p.Update(ctx, *prod, id)

	switch err {
	case nil:
	case ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never here but log the error just incase
		p.l.Println("[ERROR] serializing product", err)
	}

}

// Update is a method to update record at database
func (p *Products) Update(ctx context.Context, prod data.Product, id int) error {
	queryText := fmt.Sprintf("UPDATE tbl_product SET name = '%s', description ='%s', price = %d, sku = '%s' where id = %d",
		prod.Name,
		prod.Description,
		prod.Price,
		prod.SKU,
		id)

	fmt.Println(queryText)

	// Check ID existance before exec update
	isIDExist := p.findProductID(id)
	if isIDExist == 0 {
		return ErrProductNotFound
	}

	_, err := p.database.ExecContext(ctx, queryText)

	if err != nil {
		return err
	}
	return nil
}
