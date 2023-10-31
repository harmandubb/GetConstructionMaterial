package api

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx" // the _ allows the line to stay as is and doesn't disappear because we are telling the system that we want to intialize the package but not nessesarily use function
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Product struct {
	Name      string
	Category  string
	DataSheet byte   // FUTURE: URL to the data sheet
	Pictures  []byte // FUTURE: URLs to product images
	Price     float64
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func connectToDataBase() *sql.DB {

	err := godotenv.Load() //need to load the environmental variables in to the area before they can be used.

	// host := os.Getenv("HOST")
	// port := os.Getenv("PORT")
	// user := os.Getenv("USER")
	// password := os.Getenv("PASSWORD")
	// dbname := os.Getenv("DB_NAME")

	db_url := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Error:", err)
	}

	pingErr := db.Ping() //verifies a connection to the database is still alive, establishing a connection if necessary.
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	return db
}

func CheckDataBase() {
	db := connectToDataBase()

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	rows, err := db.Query(`SELECT current_database();`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rows)

	defer rows.Close()
	// defer stmt.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tableName)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Test Transmission is sucessful")

}

func AddProduct() {
	db := connectToDataBase()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	// stmt, err := tx.Prepare("INSERT INTO products (name) VALUES($1)")
	rows, err := db.Query(`SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rows)

	defer rows.Close()
	// defer stmt.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tableName)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	// if _, err := stmt.Exec("Fire Stop Collar", "Fire Stopping", 10.96); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := tx.Commit(); err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Test Transmission is sucessful")

}
