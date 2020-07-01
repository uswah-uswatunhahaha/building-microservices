package handlers

import (
	"net/http"

	// protos "github.com/uswah-uswatunhahaha/building-microservices/currency/protos/currency"
	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
)

// ListSingle is a handler to get product by id
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Debug("[DEBUG] get record id", id)

	// prod, err := p.getProductByID(id)
	prod, err := p.productDB.GetProductByID(id)

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

	// rr := &protos.RateRequest{
	// 	Base:        protos.Currencies(protos.Currencies_value["EUR"]),
	// 	Destination: protos.Currencies(protos.Currencies_value["GBP"]),
	// }

	// // get exchange rate
	// resp, err := p.cc.GetRate(context.Background(), rr)

	// if err != nil {
	// 	p.l.Println("[ERROR] error getting new rate", err)
	// 	data.ToJSON(&GenericError{Message: err.Error()}, rw)
	// 	return
	// }

	// p.l.Printf("Resp %#v", resp)

	// prod.Price = prod.Price * resp.Rate

	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never here but log the error just incase
		p.l.Error("[ERROR] serializing product", err)
	}
}

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	p.l.Debug("[DEBUG] get all records")

	// prod, err := p.getProducts()
	prod, err := p.productDB.GetProducts()

	err = data.ToJSON(prod, rw)

	if err != nil {
		// we should never here but log the error just incase
		p.l.Error("[ERROR] serializing product", err)
	}
}

// // GetProductByID is an exported product
// func (p *Products) getProductByID(id int) (*data.Product, error) {
// 	isIDExist := p.findProductID(id)
// 	prod := &data.Product{}

// 	if isIDExist == 0 {
// 		return nil, ErrProductNotFound
// 	}

// 	err := p.database.
// 		QueryRow("SELECT id, name, description, price, sku FROM `tbl_product` WHERE id = ?", id).
// 		Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.SKU)
// 	if err != nil {
// 		log.Fatal("Database SELECT failed")
// 		return nil, err
// 	}

// 	log.Println("You fetched a thing")

// 	return prod, nil
// }

// // GetProducts return all products from the database
// func (p *Products) getProducts() (*data.Products, error) {
// 	p.l.Println("get products data")
// 	var prods data.Products
// 	rows, err := p.database.Query("SELECT * FROM `tbl_product`")

// 	if err != nil {
// 		return nil, err
// 	}

// 	// defer rows.Close()

// 	for rows.Next() {
// 		prod := &data.Product{}
// 		err := rows.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.SKU)
// 		if err != nil {
// 			// log.Fatal("Database SELECT ALL failed")
// 			p.l.Println(err)
// 			return nil, err
// 		}
// 		prods = append(prods, *prod)
// 	}

// 	return &prods, nil
// }
