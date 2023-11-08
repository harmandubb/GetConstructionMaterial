package api

import (
	"context"
	"log"
	"reflect"
	"testing"
)

func resetTestDataBase() error {
	sqlString := "DELETE FROM products"

	p := connectToDataBase("mynewdatabase")

	tx, err := p.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), sqlString)
	if err != nil {
		return err
	}

	tx.Commit(context.Background())

	return nil

}

// func writeFileFromBytes(filePath string, data []byte) error {
// 	// Write data to filePath using os.WriteFile
// 	err := os.WriteFile(filePath, data, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func TestCheckDatabase(t *testing.T) {
	database := "mynewdatabase"
	result := CheckDataBase(database)

	if result != database {
		t.Errorf("Database = %s, but wanted %s", result, database)
	}
}

func TestAddProductBasic(t *testing.T) {
	name := "Meta Caulk Collar"
	category := "Firestopping"
	price := 10.01

	resetTestDataBase()

	AddProductBasic(name, category, price)

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
}
