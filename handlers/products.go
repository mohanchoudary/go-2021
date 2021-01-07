package handlers

import (
	"context"
	"log"
	dbfunc "main/product-api/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	intid, _ := strconv.Atoi(id)
	p.l.Println("Update Product", id)

	prod := r.Context().Value(KeyProduct{}).(*dbfunc.Product)
	// err := prod.FromJSON(r.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(rw, err.Error(), http.StatusBadRequest)
	// }
	// p.l.Printf("%#v", prod)

	err := dbfunc.UpdateProduct(intid, prod)
	if err == dbfunc.ErrProductNotFound {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(rw, "product not found", http.StatusInternalServerError)
		return
	}

}
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp, _ := dbfunc.GetProducts()
	err := lp.ToJSON(rw)
	// d, err := json.Marshal(lp)
	if err != nil {
		http.Error(rw, "JSON issue", http.StatusInternalServerError)

	}
	// rw.Write(d)
}
func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Add Product")
	prod := r.Context().Value(KeyProduct{}).(*dbfunc.Product)
	p.l.Printf("%#v", prod)
	err := dbfunc.AddProducts(prod)
	if err != nil {
		p.l.Println("Error while adding product")
	}

}

type KeyProduct struct{}

func (p *Products) MidllewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &dbfunc.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "BadJSON input", http.StatusBadRequest)
		}
		if err = prod.Validate(); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return

		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		//context.WithValue(context.Background(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
