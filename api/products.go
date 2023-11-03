package api

import (
	"database/sql"
	"fmt"
	"image"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx" // the _ allows the line to stay as is and doesn't disappear because we are telling the system that we want to intialize the package but not nessesarily use function
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Product struct {
	Name          string
	Category      string
	DataSheet     []byte // FUTURE: URL to the data sheet
	Picture       []byte // FUTURE: URLs to product images
	PictureHeight int
	PictureWidth  int
	Data_Sheet 	   []byte
	Price         float64

}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func connectToDataBase() *sql.DB {

	err := godotenv.Load() //need to load the environmental variables in to the area before they can be used.

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

func AddProductBasic(name string, category string, price float64) {
	db := connectToDataBase()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO products (name, category, price) VALUES($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close() // Close the statement when we're done with it

	if _, err := stmt.Exec(name, category, price); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Product Added Sucessfully")

}

func readImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func imageEncode(img image.Image) (int, int, []byte) {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Create a buffer to store the bytes
	var result []byte

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the color of the current pixel
			color := img.At(x, y)

			// Convert the color to RGBA
			r, g, b, a := color.RGBA()

			// Convert from 16-bit color to 8-bit color
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// Append the RGB(A) values to the result slice
			result = append(result, r8, g8, b8, a8)
		}
	}
	return width, height, result
}

func dataBaseTransmit(query string, args ...any) (bool, error){
	db := connectToDataBase()

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()	

	stmt, err := tx.Prepare(query)
	if err != nil {
		return false, err
	}

	defer stmt.Close() // Close the statement when we're done with it

	if _, err := stmt.Exec(args); err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func AddProductDataSheet(name string, pdfPath string){
	file, err := os.ReadFile(pathName)

	query := "UPDATE products
				SET data_sheet = $1
				WHERE name = $2"

	_, err := dataBaseTransmit(query, file, name)
	if err != nil {
		log.Fatal(err)
	}

}

func AddProductPicture(name string, imgPath string) {
	img, err := readImage(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	w, h, imgBytes := imageEncode(img)

	query := "INSERT INTO products (name, picture, picture_w, picture_h) 
				VALUES($1, $2, $3, $4)
				ON CONFLICT (name) 
				DO UPDATE SET 
					picture = excluded.picture,
					picture_w = excluded.picture_w,
					picture_h = excluded.picture_h"

	_, err := dataBaseTransmit(query,name,imgBytes,w,h)

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Product Image Added Successfully")

}
