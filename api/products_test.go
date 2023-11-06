package api

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/lib/pq"
)

func resetTestDataBase() (bool, error) {
	query := "DELETE FROM products"

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

	if _, err := stmt.Exec(); err != nil {
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return true, nil

}

func writeFileFromBytes(filePath string, data []byte) error {
	// Write data to filePath using os.WriteFile
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func TestAddProductBasic(t *testing.T) {
	name := "Meta Caulk Collar"
	category := "Firestopping"
	price := 10.01

	resetTestDataBase()

	AddProductBasic(name, category, price)

	_, err := AddProductDataSheet(name, "../pdf/1.pdf")
	if err != nil {
		log.Fatalln(err)
	}

	//Read the database to see if the action occured
	query := "SELECT * FROM products"
	rows, err := dataBaseRead(query)
	if err != nil {
		log.Fatalln(err)
	}

	got := Product{}
	rows.Next()

	p := reflect.ValueOf(&got).Elem()
	numCols := p.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		field := p.Field(i)
		columns[i] = field.Addr().Interface()
	}
	err = rows.Scan(columns...)
	if err != nil {
		log.Fatal(err)
	}

	if got.Name != name {
		t.Errorf("Database name = %s; want %s", got.Name, name)
	}

	if got.Category != category {
		t.Errorf("Database category = %s; want %s", got.Category, category)
	}

	if got.Price != price {
		t.Errorf("Database price = %f; want %f", got.Price, price)
	}

	datasheet := pq.Array(&got.Data_Sheet)

	fmt.Println(reflect.TypeOf(datasheet))

	input, _ := os.ReadFile("../pdf/1.pdf")

	outputLen := len(got.Data_Sheet)
	inputLen := len(input)

	if outputLen != inputLen {
		t.Errorf("The size of the input file is not the same as output")
	}

	// if got.Data_Sheet != input {
	// 	t.Errorf("Input file doesn't equal output")
	// }

	// writeFileFromBytes("../pdf/test/output.pdf", got.Data_Sheet)

}
