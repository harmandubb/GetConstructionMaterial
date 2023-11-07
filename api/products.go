package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Product struct {
	ID            int
	Name          string
	Category      string
	Picture       []byte // FUTURE: URLs to product images
	PictureWidth  *int
	PictureHeight *int
	Data_Sheet    []byte
	Price         float64
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func connectToDataBase(database string) *pgxpool.Pool {
	err := godotenv.Load() //need to load the environmental variables in to the area before they can be used.
	if err != nil {
		log.Fatalln(err)
	}

	url := os.Getenv("DB_URL")

	db_url := url + database

	dbpool, err := pgxpool.New(context.Background(), db_url)
	if err != nil {
		log.Fatal("Error:", err)
	}

	return dbpool
}

func CheckDataBase(database string) string {
	p := connectToDataBase(database)

	tx, err := p.Begin(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	defer tx.Rollback(context.Background())

	rows, err := p.Query(context.Background(), `SELECT current_database();`)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	// defer stmt.Close()

	var tableName string

	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tableName
}

func AddProductBasic(name string, category string, price float64) {
	p := connectToDataBase("mynewdatabase")

	tx, err := p.Begin(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback(context.Background())

	sqlString := "INSERT INTO products (name, category, price) VALUES($1, $2, $3)"

	cmdTag, err := tx.Exec(context.Background(), sqlString, name, category, price)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cmdTag)

}

// func dataBaseRead(query string) (*sql.Rows, error) {
// 	db := connectToDataBase("mynewdatabase")

// 	rows, err := db.Query(query) //returns a pointer to where rows are
// 	if err != nil {
// 		return nil, err
// 	}

// 	return rows, nil
// }

// func dataBaseTransmit(query string, args ...any) (bool, error) {
// 	db := connectToDataBase("mynewdatabase")

// 	tx, err := db.Begin()
// 	if err != nil {
// 		return false, err
// 	}

// 	defer tx.Rollback()

// 	stmt, err := tx.Prepare(query)
// 	if err != nil {
// 		return false, err
// 	}

// 	defer stmt.Close() // Close the statement when we're done with it

// 	_, err = stmt.Exec(args...)
// 	if err != nil {
// 		return false, err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return false, err
// 	}

// 	return true, nil
// }

// func AddProductDataSheet(name string, pdfPath string, database string) error {
// 	db := connectToDataBase(database)
// 	Tx, err := db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	lo := pgx.LargeObjects{tx: Tx}

// 	return nil
// }
