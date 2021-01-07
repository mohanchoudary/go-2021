package main

import (
	"context"
	"log"
	"main/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := mux.NewRouter()
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.Use(ph.MidllewareProductValidation)
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.Use(ph.MidllewareProductValidation)
	postRouter.HandleFunc("/", ph.AddProducts)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)

		}
	}()
	//s.ListenAndServe()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	l.Println("recived termiation,graceful shutdown", <-sigChan)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
