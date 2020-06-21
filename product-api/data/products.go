package data

// Product is entities
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
}

// Products is a collection of Product
type Products []Product

// NewProduct is a constructor
func NewProduct() *Product {
	return &Product{}
}
