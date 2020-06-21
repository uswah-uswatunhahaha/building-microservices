package server

import (
	"context"

	"github.com/hashicorp/go-hclog"

	protos "github.com/uswah-uswatunhahaha/building-microservices/currency/protos/currency"
)

// Currency is a struct
type Currency struct {
	log hclog.Logger
}

// NewCurrency is a constructor
func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{l}
}

//GetRate implementation of CurrencyServer interface
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateReponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	return &protos.RateReponse{Rate: 0.5}, nil
}
