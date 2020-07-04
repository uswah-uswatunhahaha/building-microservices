package main

import (
	"fmt"
	"testing"

	"github.com/uswah-uswatunhahaha/building-microservices/product-api/sdk/client"
	"github.com/uswah-uswatunhahaha/building-microservices/product-api/sdk/client/product"
)

func TestOurClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := product.NewListProductsParams()
	prod, err := c.Product.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v", prod.GetPayload()[0])
	t.Fail()

}
