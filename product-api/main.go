package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-openapi/runtime/middleware"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/uswah-uswatunhahaha/building-microservices/product-api/data"
	"github.com/uswah-uswatunhahaha/building-microservices/product-api/db"
	"github.com/uswah-uswatunhahaha/building-microservices/product-api/handlers"
	"google.golang.org/grpc"

	protos "github.com/uswah-uswatunhahaha/building-microservices/currency/protos/currency"
)

func main() {
	// l := log.New(os.Stdout, "Mysql - Go", log.LstdFlags)
	l := hclog.Default()
	l.Info("Running Product-Api main")
	database, err := db.CreateDatabase()

	if err != nil {
		l.Error("Database connection failed: %s", "error", err.Error())
	}

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Create client
	cc := protos.NewCurrencyClient(conn)

	// create database instance
	db := data.NewProductsDB(database, cc, l)
	// data.NewProductsDB(cc, l)

	//create the handlers
	ph := handlers.NewProducts(l, db)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/product/{id:[0-9]+}", ph.ListSingle).Queries("currency", "[A-Z]{3}")
	getR.HandleFunc("/product/{id:[0-9]+}", ph.ListSingle)

	getR.HandleFunc("/product", ph.ListAll).Queries("currency", "[A-Z]{3}")
	getR.HandleFunc("/product", ph.ListAll)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/product", ph.Create)
	postR.Use(ph.MiddlewareValidateProduct)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/product/{id:[0-9]+}", ph.Edit)
	putR.Use(ph.MiddlewareValidateProduct)

	delR := sm.Methods(http.MethodDelete).Subrouter()
	delR.HandleFunc("/product/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := &http.Server{
		Addr:         ":9090",                                          // configure the address
		Handler:      ch(sm),                                           //set the default handler
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}), // set the logger for the server
		ReadTimeout:  5 * time.Second,                                  // max time to read request from the client
		WriteTimeout: 10 * time.Second,                                 //max time ro write response to the client
		IdleTimeout:  120 * time.Second,                                //max time for connections using TCP Keep=Alive
	}

	// start the server
	go func() {
		l.Info("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Error("Error starting server: %s\n", "error ", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)

}
