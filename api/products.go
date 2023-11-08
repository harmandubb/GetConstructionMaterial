package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
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

	tx.Commit(context.Background())

}

func dataBaseRead(sqlString string) (pgx.Rows, error) {
	p := connectToDataBase("mynewdatabase")

	rows, err := p.Query(context.Background(), sqlString) //returns a pointer to where rows are
	if err != nil {
		return rows, err
	}

	return rows, nil
}

func dataBaseTransmit(sqlString string, database string, args ...any) error {
	p := connectToDataBase(database)

	tx, err := p.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), sqlString, args...)
	if err != nil {
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}

func AddProductDataSheet(name string, pdfPath string, database string) error {
	p := connectToDataBase(database)
	tx, err := p.Begin(context.Background())
	if err != nil {
		return err
	}

	//can start to initiative the large objects process
	los := tx.LargeObjects()

	oidVal, err := los.Create(context.Background(), 0)
	if err != nil {
		return err
	}

	fmt.Println(oidVal)

	// should I upload the oid number to the table section
	lo, err := los.Open(context.Background(), oidVal, pgx.LargeObjectModeWrite)
	if err != nil {
		return err
	}

	defer lo.Close()

	// Can write the pdf to the large object since I have the  connection established.
	file, err := os.Open(pdfPath)
	if err != nil {
		return err
	}

	var fileBytes []byte

	_, err = file.Read(fileBytes)
	if err != nil {
		return err
	}

	_, err = lo.Write(fileBytes)
	if err != nil {
		return err
	}

	//store the oid value in the database table
	sqlString := "UPDATE products SET data_sheet=$1 WHERE name='$2'"

	_, err = tx.Exec(context.Background(), sqlString, oidVal, name)
	if err != nil {
		return err
	}

	return nil
}
