package api

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Product struct {
	ID            int
	Name          string
	Created       time.Time
	Category      string
	Picture       *uint32 // FUTURE: URLs to product images
	PictureWidth  *uint
	PictureHeight *uint
	Data_Sheet    *uint32
	Price         float64
	Currency      string
}

func AddProductBasic(name string, category string, price float64, p *pgxpool.Pool) {
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

// func readDataBaseRow(tableName string, productName string) (Product, error) {
// 	var product Product
// 	sqlString := fmt.Sprintf("SELECT * FROM %s WHERE name = '%s'", tableName, productName)
// 	rows, err := dataBaseRead(sqlString)
// 	if err != nil {
// 		return product, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {

// 		err := rows.Scan(
// 			&product.ID,
// 			&product.Name,
// 			&product.Created,
// 			&product.Category,
// 			&product.Picture,
// 			&product.PictureWidth,
// 			&product.PictureHeight,
// 			&product.Data_Sheet,
// 			&product.Price,
// 			&product.Currency,
// 		)

// 		if err != nil {
// 			return product, err
// 		}
// 	}

// 	return product, nil

// }

// func AddProductDataSheet(name string, pdfPath string, database string, p *pgxpool.Pool) (uint32, error) {
// 	tx, err := p.Begin(context.Background())
// 	if err != nil {
// 		return 0, err
// 	}

// 	//can start to initiative the large objects process
// 	los := tx.LargeObjects()

// 	oidVal, err := los.Create(context.Background(), 0)
// 	if err != nil {
// 		return 0, err
// 	}

// 	fmt.Println(oidVal)

// 	// should I upload the oid number to the table section
// 	lo, err := los.Open(context.Background(), oidVal, pgx.LargeObjectModeWrite)
// 	if err != nil {
// 		return 0, err
// 	}

// 	defer lo.Close()

// 	// Can write the pdf to the large object since I have the  connection established.
// 	file, err := os.Open(pdfPath)
// 	if err != nil {
// 		return 0, err
// 	}

// 	fileBytes, err := io.ReadAll(file)
// 	if err != nil {
// 		return 0, err
// 	}

// 	_, err = lo.Write(fileBytes)
// 	if err != nil {
// 		return 0, err
// 	}

// 	//store the oid value in the database table
// 	sqlString := "UPDATE products SET data_sheet=$1 WHERE name=$2"

// 	_, err = tx.Exec(context.Background(), sqlString, oidVal, name)
// 	if err != nil {
// 		return 0, err
// 	}
// 	err = tx.Commit(context.Background())
// 	if err != nil {
// 		return 0, err
// 	}

// 	return oidVal, nil
// }

// func addProduct(pInfo ProductInfo) error {
// 	//Can make the function more general for choosing different databases
// 	p := connectToDataBase("mynewdatabase")

// 	tx, err := p.Begin(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	defer tx.Rollback(context.Background())

// 	//TODO: Can make the function more general for choosing differeint tables

// 	sqlString := "INSERT INTO products (created, name, category, price, currency) VALUES($1, $2, $3, $4, $5)"

// 	_, err = tx.Exec(context.Background(), sqlString, pInfo.Date, pInfo.Name, pInfo.Category, pInfo.Price, pInfo.Currency)
// 	if err != nil {
// 		return err
// 	}

// 	tx.Commit(context.Background())

// 	//Now to upload the datasheet

// 	//TODO Implement the function to upload and save multiple datasheets
// 	if len(pInfo.DataSheet) >= 1 {
// 		_, err = AddProductDataSheet(pInfo.Name, pInfo.DataSheet[0], "mynewdatabase", p)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func getProductDataSheet(oidVal uint32, database string, outputPath string) error {
// 	p := connectToDataBase(database)
// 	tx, err := p.Begin(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	los := tx.LargeObjects()

// 	lo, err := los.Open(context.Background(), oidVal, pgx.LargeObjectModeRead)
// 	if err != nil {
// 		return err
// 	}

// 	file, err := os.Create(outputPath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	//buffer
// 	buffer := make([]byte, 1048576)

// 	//Read in a loop
// 	for {
// 		n, err := lo.Read(buffer)
// 		if err != nil {
// 			if err == io.EOF {
// 				if _, err := file.Write(buffer[:n]); err != nil {
// 					return err
// 				}
// 				break
// 			}
// 			return err
// 		}

// 		if n == 0 {
// 			break
// 		}

// 		if _, err := file.Write(buffer[:n]); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func createImageFromBytes(colorBytes []byte, img_w int, img_h int, imgOutput string) error {
// 	img := image.NewRGBA(image.Rect(0, 0, img_w, img_h))

// 	idx := 0

// 	for y := 0; y < img_h; y++ {
// 		for x := 0; x < img_w; x++ {
// 			//Extract the rgba components
// 			r := colorBytes[idx]
// 			g := colorBytes[idx+1]
// 			b := colorBytes[idx+2]
// 			a := colorBytes[idx+3]

// 			//set the pixal color in the new image
// 			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: a})

// 			//Move to the next color in the array
// 			idx += 4
// 		}
// 	}

// 	//Encode the new image to a file (will go with png as a standard)
// 	outFile, err := os.Create(imgOutput + ".png")
// 	if err != nil {
// 		return err
// 	}

// 	defer outFile.Close()

// 	if err := png.Encode(outFile, img); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func AddProductPicture(name string, imgPath string, database string) (uint32, int, int, error) {
// 	p := connectToDataBase(database)
// 	tx, err := p.Begin(context.Background())
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}

// 	//can start to initiative the large objects process
// 	los := tx.LargeObjects()

// 	oidVal, err := los.Create(context.Background(), 0)
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}

// 	// should I upload the oid number to the table section
// 	lo, err := los.Open(context.Background(), oidVal, pgx.LargeObjectModeWrite)
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}

// 	defer lo.Close()

// 	// Can write the pdf to the large object since I have the  connection established.
// 	reader, err := os.Open(imgPath)
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}

// 	defer reader.Close()

// 	img, format, err := image.Decode(reader)
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}

// 	fmt.Println(format) //not a good method becuase I need to know the format it was saved in.

// 	bounds := img.Bounds()

// 	pic_w := bounds.Max.X - bounds.Min.X
// 	pic_h := bounds.Max.Y - bounds.Min.Y

// 	var colorBytes []byte

// 	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
// 		for x := bounds.Min.X; x < bounds.Max.X; x++ {
// 			// Get the color of the pixel
// 			r, g, b, a := img.At(x, y).RGBA()

// 			// Convert color values to bytes and append to the slice
// 			// Note: RGBA returns color values in the range [0, 65535].
// 			// Convert them to [0, 255] if necessary.
// 			// lossing infromation when converting and when you get the values back I would need to decode back by making the values a significant byte
// 			colorBytes = append(colorBytes, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
// 		}
// 	}

// 	_, err = lo.Write(colorBytes)
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}

// 	//store the oid value in the database table
// 	sqlString := "UPDATE products SET picture=$1, picture_w=$2, picture_h=$3 WHERE name=$4"

// 	_, err = tx.Exec(context.Background(), sqlString, oidVal, pic_w, pic_h, name)
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}
// 	err = tx.Commit(context.Background())
// 	if err != nil {
// 		return 0, 0, 0, err
// 	}

// 	return oidVal, pic_w, pic_h, nil
// }

// func getProductPicture(oidVal uint32, img_w int, img_h int, database string, outputPath string) error {
// 	p := connectToDataBase(database)
// 	tx, err := p.Begin(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	los := tx.LargeObjects()

// 	lo, err := los.Open(context.Background(), oidVal, pgx.LargeObjectModeRead)
// 	if err != nil {
// 		return err
// 	}

// 	//buffer
// 	buffer := make([]byte, 1024)
// 	var imgBuffer []byte

// 	//Read in a loop
// 	for {
// 		n, err := lo.Read(buffer)
// 		if err != nil {
// 			if err == io.EOF {
// 				imgBuffer = append(imgBuffer, buffer[:n]...)
// 				break
// 			}
// 			return err
// 		}

// 		if n == 0 {
// 			break
// 		}

// 		imgBuffer = append(imgBuffer, buffer[:n]...)
// 	}

// 	createImageFromBytes(imgBuffer, img_w, img_h, "./output")

// 	return nil
// }
