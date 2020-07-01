package server

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/uswah-uswatunhahaha/building-microservices/currency/data"
	protos "github.com/uswah-uswatunhahaha/building-microservices/currency/protos/currency"
)

// Currency is a struct
type Currency struct {
	rates *data.ExchangeRates
	log   hclog.Logger
}

// NewCurrency is a constructor
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{r, l}
}

//GetRate implementation of CurrencyServer interface
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateReponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	// call our GetRate method
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	// return &protos.RateReponse{Rate: 0.5}, nil
	return &protos.RateReponse{Rate: rate}, nil
}
