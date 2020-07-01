package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/uswah-uswatunhahaha/building-microservices/currency/data"
	"github.com/uswah-uswatunhahaha/building-microservices/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	protos "github.com/uswah-uswatunhahaha/building-microservices/currency/protos/currency"
)

const (
	port = ":9092"
)

func main() {
	log := hclog.Default()
	rates, err := data.NewRates(log)
	if err != nil {
		log.Error("Unable to generate rate", "error", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()
	cs := server.NewCurrency(rates, log)

	protos.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}

	// gs.Serve(l)
	if err := gs.Serve(l); err != nil {
		log.Error("failed to serve: %v", err)
	}
}
