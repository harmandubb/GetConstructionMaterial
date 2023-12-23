package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Purpose: provide a way to connect to a database that is running locally
// Parameters:
// database string --> name of the database you are trying to connect to
// Return:
// *pgxpool.Pool --> pointer that allows you to interface with the database
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

// Purpose: Verfies that you have connected to the databse
// Parameter:
// Database string --> name of the database
// Return:
// tableName string --> name of the table name in the database?? (TODO: what if there are multiple tables in the database?)
func CheckDataBase(database string) (tableName string) {
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

// Purpose: read infromation that is in a database for a SQL prompt
// Parameters:
// sqlString string --> Prompt that is used to get information from a table
// Return:
// pgx.Rows --> infromation in a row style for the result of the sqlString prompt
// Errors if present
func dataBaseRead(sqlString string) (pgx.Rows, error) {
	p := connectToDataBase("mynewdatabase")

	rows, err := p.Query(context.Background(), sqlString) //returns a pointer to where rows are
	if err != nil {
		return rows, err
	}

	return rows, nil
}

// Purpose: Transmit informatin to the database
// Parameters:
// sqlString string --> prompt in the style of a sqlString to trasmit the data to the database
// database string --> name of the database that you are tyring to trasmit info to
// args ..any --> the arguements to the sqlString that are to fill the infromation needed in the table that you are targetting
// Return:
// Error if any present

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
