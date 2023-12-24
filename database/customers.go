package database

import (
	g "docstruction/getconstructionmaterial/GCalls"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerInquiry struct {
	ID            int
	Email         string
	Time_Inquired time.Time
	Material      string
	Loc           string
	Present       bool
	Price         float64
	Currency      string
	Data_Sheet    *uint32
}

// Purpose: create a customer inquiry row that has all basic info about the user and material filled
// Parameters:
// matForm g.MaterialFormInfo --> struct holding the info inputted by the user
// database string --> database name that you want to store the info in
// tableNmae string --> name of the table that you want to input the data into
// Return:
// Error if present
func AddBlankCustomerInquiry(p *pgxpool.Pool, matForm g.MaterialFormInfo, database string, tableName string) (err error) {
	sqlString := fmt.Sprintf("INSERT INTO %s (Email, Time_Inquired, Material, Loc, Present, Price, Currency, Data_Sheet) VALUES($1, $2, $3, $4, $5, $6, $7, $8)", tableName)

	err = dataBaseTransmit(p, sqlString, database, matForm.Email, time.Now(), matForm.Material, matForm.Loc, false, 0, "", nil)
	if err != nil {
		return err
	}

	// Implement the read for this test

	return nil
}

// Purpose: read the entire row that is related to a customer inquiry
// Parameters:
// tableNmae string --> table name in postgres that you want to get the informaiton from
// sqlString string --> the command that will be sent to the table to read row information
// args any --> arguements needed to fill and accomplish the readRowInfo
// Return:
// custInquiry CustomerInquiry struct --> Retruns what ever infromation is present for the customer in quiry so far
// Error if present

func readCustomerInquiry(tableName string, customerEmail string) (custInquiry CustomerInquiry, err error) {
	sqlString := fmt.Sprintf("SELECT * FROM %s WHERE email = '%s'", tableName, customerEmail)
	rows, err := dataBaseRead(sqlString)
	if err != nil {
		return CustomerInquiry{}, err
	}

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			&custInquiry.ID,
			&custInquiry.Email,
			&custInquiry.Time_Inquired,
			&custInquiry.Material,
			&custInquiry.Loc,
			&custInquiry.Present,
			&custInquiry.Price,
			&custInquiry.Currency,
			&custInquiry.Data_Sheet,
		)

		if err != nil {
			return CustomerInquiry{}, err
		}
	}

	return custInquiry, nil

}

// Purpose: update the customer inquiry row given a customer email
// Parameters:
// p *pgxpool.Pool --> pointer to the databse connection
// databse string --> database that you want to connect to
// tableName string --> table that hods the data that you would like to change
// customerEmailer string --> email will be used to pull out a row to update
// col string --> column you want to update
// val any --> variable to update the col with
// Return:
// Error if any present

func updateCustomerInquiry(p *pgxpool.Pool, database string, tableName string, customerEmail string, col string, val any) (err error) {
	sqlString := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE email = $2", tableName, col)

	err = dataBaseTransmit(p, sqlString, database, val, customerEmail)
	if err != nil {
		return err
	}

	return nil
}

// Purpose: Speciailized version to update Present of customer inquiry given an email
func updateCustomerInquiryPresent(p *pgxpool.Pool, database string, tableName string, customerEmail string, present bool) (err error) {
	err = updateCustomerInquiry(p, database, tableName, customerEmail, "present", present)
	if err != nil {
		return err
	}

	return nil
}

// Purpose: Speciailized version to update Price and the currency of customer inquiry given an email
func updateCustomerInquiryPrice(p *pgxpool.Pool, database string, tableName string, customerEmail string, price float64, currency string) (err error) {
	err = updateCustomerInquiry(p, database, tableName, customerEmail, "price", price)
	if err != nil {
		return err
	}

	err = updateCustomerInquiry(p, database, tableName, customerEmail, "currency", currency)
	if err != nil {
		return err
	}

	return nil
}

// Purpose: Speciailized version to update DATASHEET of customer inquiry given an emai.
// ---------TODO------

func updateCustomerInquiryDataSheet(p *pgxpool.Pool, database string, tableName string, customerEmail string, price float64, currency string) (err error) {
	return nil
}
