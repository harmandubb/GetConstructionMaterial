package api

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx" // the _ allows the line to stay as is and doesn't disappear because we are telling the system that we want to intialize the package but not nessesarily use function
	_ "github.com/lib/pq"
)

type Product struct {
	Name      string
	Category  string
	DataSheet byte   // FUTURE: URL to the data sheet
	Pictures  []byte // FUTURE: URLs to product images
	Price     float64
}

// adds a new product to the data base
func AddProduct() {
	var db *sql.DB
	var err error
	// db_url := os.Getenv("db_URL")
	db, err = sql.Open("pgx", "")
	if err != nil {
		log.Fatal("Error:", err)
	}

	defer db.Close()

	pingErr := db.Ping() //verifies a connection to the database is still alive, establishing a connection if necessary.
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

}
