package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// Create handles POST Request
func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[POST] save records")
	rw.Header().Add("Content-Type", "application/json")

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	// err = p.insert(ctx, prod)
	err := p.insert(ctx, *prod)
	if err != nil {
		p.l.Println(err)
	}

}

// Insert is a method to save record to database
func (p *Products) insert(ctx context.Context, prod data.Product) error {
	table := "tbl_product"
	queryText := fmt.Sprintf("INSERT INTO %v (name, description, price, sku) values('%v','%v',%v,'%v')", table,
		prod.Name,
		prod.Description,
		prod.Price,
		prod.SKU)

	_, err := p.database.ExecContext(ctx, queryText)

	if err != nil {
		return err
	}
	return nil
}
