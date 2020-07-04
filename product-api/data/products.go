package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/hashicorp/go-hclog"
	protos "github.com/uswah-uswatunhahaha/building-microservices/currency/protos/currency"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`
}

// Products is a collection of Product
type Products []Product

// ProductsDB is a type of Products and call protos
type ProductsDB struct {
	database *sql.DB
	currency protos.CurrencyClient
	log      hclog.Logger
}

// NewProductsDB is an idiomatic function of Go
func NewProductsDB(db *sql.DB, c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	return &ProductsDB{db, c, l}
}

// GetProducts return all products from the database (data pkg)
func (pdb *ProductsDB) GetProducts(currency string) (*Products, error) {

	var prods Products
	rows, err := pdb.database.Query("SELECT * FROM `tbl_product`")

	if err != nil {
		return nil, err
	}

	// defer rows.Close()

	for rows.Next() {
		prod := &Product{}
		err := rows.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.SKU)
		if err != nil {
			pdb.log.Error("error", err)
			return nil, err
		}
		prods = append(prods, *prod)
	}

	if currency == "" {
		return &prods, nil
	}

	rate, err := pdb.getRate(currency)
	if err != nil {
		pdb.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	for _, p := range prods {
		np := p
		np.Price = np.Price * rate
		prods = append(prods, np)

	}

	return &prods, nil
}

// GetProductByID returns product by id
func (pdb *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	isIDExist := pdb.findIndexByProductID(id)
	prod := &Product{}

	if isIDExist == 0 {
		return nil, ErrProductNotFound
	}

	err := pdb.database.
		QueryRow("SELECT id, name, description, price, sku FROM `tbl_product` WHERE id = ?", id).
		Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.SKU)
	if err != nil {
		log.Fatal("Database SELECT failed")
		return nil, err
	}

	log.Println("You fetched a thing")

	if currency == "" {
		return prod, nil
	}

	rate, err := pdb.getRate(currency)
	if err != nil {
		pdb.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	np := *prod
	np.Price = np.Price * rate

	return &np, nil
}

// UpdateProduct is a method to update record at database
func (pdb *ProductsDB) UpdateProduct(ctx context.Context, prod Product, id int) error {
	queryText := fmt.Sprintf("UPDATE tbl_product SET name = '%s', description ='%s', price = %f, sku = '%s' where id = %d",
		prod.Name,
		prod.Description,
		prod.Price,
		prod.SKU,
		id)

	fmt.Println(queryText)

	// Check ID existance before exec update
	isIDExist := pdb.findIndexByProductID(id)
	if isIDExist == 0 {
		return ErrProductNotFound
	}

	_, err := pdb.database.ExecContext(ctx, queryText)

	if err != nil {
		return err
	}
	return nil
}

// AddProduct is a method to save record to database
func (pdb *ProductsDB) AddProduct(ctx context.Context, prod Product) error {
	table := "tbl_product"
	queryText := fmt.Sprintf("INSERT INTO %v (name, description, price, sku) values('%v','%v',%v,'%v')", table,
		prod.Name,
		prod.Description,
		prod.Price,
		prod.SKU)

	_, err := pdb.database.ExecContext(ctx, queryText)

	if err != nil {
		return err
	}
	return nil
}

// DeleteProduct is method to delete record from mysql
func (pdb *ProductsDB) DeleteProduct(ctx context.Context, prod Product, id int) error {
	queryText := fmt.Sprintf("DELETE FROM tbl_product where id = %d", id)

	fmt.Println(queryText)

	// Check ID existance before exec update
	isIDExist := pdb.findIndexByProductID(id)
	if isIDExist == 0 {
		return ErrProductNotFound
	}

	_, err := pdb.database.ExecContext(ctx, queryText)

	if err != nil {
		return err
	}
	return nil
}

func (pdb *ProductsDB) findIndexByProductID(id int) int {
	prod := &Product{}
	pdb.database.QueryRow("SELECT EXISTS(SELECT * FROM tbl_product WHERE id=?)", id).Scan(&prod.ID)

	if prod.ID == 0 {
		log.Println("ID Not Found")
		return prod.ID
	}

	log.Printf("Found the ID %d", id)

	return prod.ID
}

func (pdb *ProductsDB) getRate(destination string) (float64, error) {
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destination]),
	}

	resp, err := pdb.currency.GetRate(context.Background(), rr)

	return resp.Rate, err
}
