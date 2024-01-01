package database

import (
	"context"
	g "docstruction/getconstructionmaterial/GCalls"

	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerInquiry struct {
	ID                       int
	Inquiry_ID               string
	Email                    string
	Time_Inquired            time.Time
	Material                 string
	Loc                      string
	Present                  bool
	Supplier_Email_Thread_ID string //Include the supplier thread ID that has the best overall offering based on price
	Price                    float64
	Currency                 string
	Data_Sheet               *[]uint8
}

// Purpose: create a customer inquiry row that has all basic info about the user and material filled
// Parameters:
// p *pgxpool.Pool -->
// matForm g.MaterialFormInfo --> struct holding the info inputted by the user
// database string --> database name that you want to store the info in
// tableNmae string --> name of the table that you want to input the data into
// Return:
// Error if present
func AddBlankCustomerInquiry(p *pgxpool.Pool, matForm g.MaterialFormInfo, tableName, currency string) (inquiryID string, err error) {

	sqlString := fmt.Sprintf("INSERT INTO %s (email, inquiry_id, time_inquired, material, loc, present, supplier_email_thread_id, price, currency, data_sheet) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", tableName)

	inquiryID = generateInquiryID()

	err = dataBaseTransmit(p, sqlString, matForm.Email, inquiryID, time.Now(), matForm.Material, matForm.Loc, false, "", 0, currency, nil)
	if err != nil {
		return "", err
	}

	return inquiryID, nil

}

// Purpose: create a customer inquiry row that has all basic info about the user and material filled
// Parameters:
// inquiryIDStream chan<- string : channel to return inquiryID which links all things realted to this inquiry
// errStream chan<- error: channel for error return
// ct context.Context: allows for canel of the function to occur due to time exceeded
// p *pgxpool.Pool --> pool of database connnections (to be made through a pool share in sync package)
// matForm g.MaterialFormInfo --> struct holding the info inputted by the used
// tableNmae string --> name of the table that you want to input the data into
func ConcurrentAddBlankCustomerInquiry(inquiryIDStream chan<- string, errStream chan<- error, ctx context.Context,
	p *pgxpool.Pool, matForm g.MaterialFormInfo, currency, tableName string) {

	select {
	case <-ctx.Done():
		//The context has been canncelled
		errStream <- ctx.Err()
		return

	default:
		sqlString := fmt.Sprintf("INSERT INTO %s (Email, Inquiry_ID, Time_Inquired, Material, Loc, Present, supplier_email_thread_id, Price, Currency, Data_Sheet) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", tableName)

		inquiryID := generateInquiryID()

		err := dataBaseTransmit(p, sqlString, matForm.Email, inquiryID, time.Now(), matForm.Material, matForm.Loc, false, "", 0, currency, nil)
		if err != nil {
			errStream <- err
			return
		}

		inquiryIDStream <- inquiryID
	}

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

func updateCustomerInquiry(p *pgxpool.Pool, database string, tableName string, inquiry_id string, col string, val any) (err error) {
	sqlString := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE inquiry_id = $2", tableName, col)

	err = dataBaseTransmit(p, sqlString, database, val, inquiry_id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCustomerInquiryMaterial(p *pgxpool.Pool, tableName string, inquiry_id string, supplier_email_thread_id string, price float64, currency string, Data_Sheet *[]byte) (err error) {
	sqlString := fmt.Sprintf("UPDATE %s SET present = $1, supplier_email_thread_id = $2, price = $3, currency = $4, data_sheet = $5 WHERE inquiry_id = '%s'", tableName, inquiry_id)

	err = dataBaseTransmit(p, sqlString, true, supplier_email_thread_id, price, currency, Data_Sheet)
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

// Purpose: Allow a way to generate an ID per inquiry that is consistent between different tables
// Returns:
// id string --> unique ID to relate the
func generateInquiryID() (id string) {
	id = uuid.New().String()

	fmt.Println(len(id))
	return id
}

// Purpose: read the entire row that is related to a customer inquiry
// Parameters:
// tableNmae string --> table name in postgres that you want to get the informaiton from
// sqlString string --> the command that will be sent to the table to read row information
// args any --> arguements needed to fill and accomplish the readRowInfo
// Return:
// custInquiry CustomerInquiry struct --> Retruns what ever infromation is present for the customer in quiry so far
// Error if present

func ReadCustomerInquiry(p *pgxpool.Pool, tableName string, inquiryID string) (custInquiry CustomerInquiry, err error) {
	sqlString := fmt.Sprintf("SELECT * FROM %s WHERE inquiry_id = '%s'", tableName, inquiryID)
	rows, err := dataBaseRead(p, sqlString)
	if err != nil {
		return CustomerInquiry{}, err
	}

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			&custInquiry.ID,
			&custInquiry.Inquiry_ID,
			&custInquiry.Email,
			&custInquiry.Time_Inquired,
			&custInquiry.Material,
			&custInquiry.Loc,
			&custInquiry.Present,
			&custInquiry.Supplier_Email_Thread_ID,
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
