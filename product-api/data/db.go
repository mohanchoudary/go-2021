package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "helloworld"
	hostname = "localhost:3308"
	dbname   = "ecommerce"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func dbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return nil, err
	}
	//defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return nil, err
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return nil, err
	}
	log.Printf("rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return nil, err
	}
	//defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return nil, err
	}
	log.Printf("Connected to DB %s successfully\n", dbname)
	return db, nil
}

var (
	ID          int
	Name        string
	Description string
	Price       float32
	SKU         string
	CreatedOn   string
	UpdatedOn   string
	DeletedOn   string
)

// Following function gets the products from Product_api table and adds to productList slice
func GetProducts() (Products, error) {
	sqlStatement := `SELECT * FROM product_api;`
	db, err := dbConnection()
	if err != nil {
		fmt.Println("error while openng db")
	}

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Printf("Error %s when retrieving row from products table", err)
		return nil, err
	}
	productList = productList[:0]
	for rows.Next() {
		rows.Scan(&ID, &Name, &Description, &Price, &SKU, &CreatedOn, &UpdatedOn)

		p := Product{
			ID:          ID,
			Name:        Name,
			Description: Description,
			Price:       Price,
			SKU:         SKU,
			CreatedOn:   CreatedOn,
			UpdatedOn:   UpdatedOn,
		}
		fmt.Println(ID, Name, Description, Price, CreatedOn, UpdatedOn)
		productList = append(productList, &p)
	}
	db.Close()

	return productList, nil
}

// Following function gets the product details from request body and add it to
//Product_api table

func AddProducts(p *Product) error {
	sqlStatement := ` insert into product_api(product_id,name, Description, Price, SKU) 
	values (?,?,?,?,?);`
	db, err := dbConnection()
	if err != nil {
		fmt.Println("error while openng db")
	}
	rows, err := db.Query("SELECT COUNT(*) as count FROM  product_api")
	if err != nil {
		return fmt.Errorf("Add product failed")
	}
	p.ID = checkCount(rows) + 1

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, sqlStatement)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, p.ID, p.Name, p.Description, p.Price, p.SKU)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}
	a, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d products created ", a)

	productList = append(productList, p)
	return nil

}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Println("Error while finding the number of rows.")
		}
	}
	return count
}
