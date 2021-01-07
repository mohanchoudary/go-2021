package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int `json:"id" validate:"required"`
	Name        string
	Description string
	Price       float32 `validate:"gt=0"`
	SKU         string  `validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

var productList = []*Product{}

func (p *Product) Validate() error {

	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)

}

func validateSKU(fl validator.FieldLevel) bool {

	//re := regexp.MustCompile(`[a-zA-Z]+-[a-zA-Z]+[a-zA-Z]+`)
	re := regexp.MustCompile(`[a-zA-Z]+-[a-zA-Z]+[a-zA-Z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)
	if len(matches) != 1 {
		return false
	}

	return true

}

func UpdateProduct(id int, p *Product) error {
	//p.ID = id
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	db, err := dbConnection()
	query := "UPDATE product_api SET name= ?, Price=?, Description=?, SKU=? WHERE product_id =?;"
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, p.Name, p.Price, p.Description, p.SKU, p.ID)

	if err != nil {
		fmt.Println("execution")
		fmt.Println(err)
	}
	no, _ := res.RowsAffected()
	fmt.Printf("updated number of products: %d", no)
	productList[pos] = p
	fmt.Println("product updated")
	return nil

}

var ErrProductNotFound = fmt.Errorf("Product not found...")

func findProduct(id int) (*Product, int, error) {
	productList, _ := GetProducts()
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, 0, ErrProductNotFound
}

type Products []*Product

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}
