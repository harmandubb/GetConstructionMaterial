package api

import (
	"bufio"
	"database/sql"
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
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

func connectToDataBase(database string) *sql.DB {
	err := godotenv.Load() //need to load the environmental variables in to the area before they can be used.

	url := os.Getenv("DB_URL")

	db_url := url + database

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Error:", err)
	}

	pingErr := db.Ping() //verifies a connection to the database is still alive, establishing a connection if necessary.
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return db
}

func CheckDataBase() {
	db := connectToDataBase("mynewdatabase")

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
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func AddProductBasic(name string, category string, price float64) {
	db := connectToDataBase("mynewdatabase")

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

}

func readImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	//Decode image
	config, fileType, _ := image.DecodeConfig(reader)

	fmt.Println(config)

	fmt.Println(fileType)

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

func dataBaseRead(query string) (*sql.Rows, error) {
	db := connectToDataBase("mynewdatabase")

	rows, err := db.Query(query) //returns a pointer to where rows are
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func dataBaseTransmit(query string, args ...any) (bool, error) {
	db := connectToDataBase("mynewdatabase")

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

	_, err = stmt.Exec(args...)
	if err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
}

func AddProductDataSheet(name string, pdfPath string) (bool, error) {
	file, err := os.ReadFile(pdfPath)

	byteaFile := pq.Array(file)

	os.WriteFile("../pdf/test/test.pdf", file, 0644)

	if err != nil {
		return false, err
	}

	query := "UPDATE products SET data_sheet = $2 WHERE name = $1"

	_, err = dataBaseTransmit(query, name, byteaFile)
	if err != nil {
		return false, err
	}

	return true, nil

}

func AddProductPicture(name string, imgPath string) {
	img, err := readImage(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	w, h, imgBytes := imageEncode(img)

	query := "INSERT INTO products (name, picture, picture_w, picture_h) VALUES($1, $2, $3, $4) ON CONFLICT (name) DO UPDATE SET picture = excluded.picture, picture_w = excluded.picture_w, picture_h = excluded.picture_h"

	_, err = dataBaseTransmit(query, name, imgBytes, w, h)

	if err != nil {
		log.Fatal(err)
	}
}
