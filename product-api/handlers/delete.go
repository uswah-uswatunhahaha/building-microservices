package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// Delete handles DELETE request
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DELETE] delete records")
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	var prod data.Product

	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		p.l.Println(err)
	}

	err = p.hapus(ctx, prod, id)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		p.l.Println(err)
	}

}

// Hapus is method to delete record from mysql
func (p *Products) hapus(ctx context.Context, prod data.Product, id int) error {
	queryText := fmt.Sprintf("DELETE FROM tbl_product where id = %d", id)

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
