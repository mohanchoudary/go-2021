// package main

// import (
// 	"context"
// 	"log"
// 	"main/handlers"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"time"
// )

// func main() {
// 	l := log.New(os.Stdout, "product-api", log.LstdFlags)

// 	hh := handlers.NewHello(l)
// 	gh := handlers.NewGoodbye(l)
// 	sm := http.NewServeMux()

// 	ph := handlers.NewProducts(l)
// 	sm.Handle("/products/", ph)

// 	sm.Handle("/", hh)
// 	sm.Handle("/bye", gh)
// 	//http.Handle("/", hh)
// 	//http.HandleFunc("/", hh)

// 	s := &http.Server{
// 		Addr:         ":9090",
// 		Handler:      sm,
// 		IdleTimeout:  120 * time.Second,
// 		ReadTimeout:  120 * time.Second,
// 		WriteTimeout: 120 * time.Second,
// 	}

// 	go func() {
// 		err := s.ListenAndServe()
// 		if err != nil {
// 			l.Fatal(err)

// 		}
// 	}()
// 	//s.ListenAndServe()
// 	sigChan := make(chan os.Signal)
// 	signal.Notify(sigChan, os.Interrupt)
// 	signal.Notify(sigChan, os.Kill)

// 	l.Println("recived termiation,graceful shutdown", <-sigChan)

// 	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
// 	s.Shutdown(tc)
// }
