package database

import (
	"fmt"
	"time"

	g "docstruction/getconstructionmaterial/GCalls"

	"github.com/jackc/pgx/v5/pgxpool"
)

// This Fie contains functions and functionatiliy related to updating a database tables that holds information around the email communication sent to suppliers

type EmailInquiries struct {
	ID              int
	Inquiry_ID      string
	Thread_ID       string
	Client_Email    string
	Time_Sent       time.Time
	Material        string
	Supplier_Map_ID string
	Supplier_Name   string
	Supplier_Lat    float64
	Supplier_Lng    float64
	Supplier_Email  string
	Sent_Out        bool
	Present         bool
	Price           float64
	Currency        string
	Data_Sheet      *[]byte
} //15 items

type IDOption struct {
	Inquiry_ID string
	Thread_ID  string
}

// Purpose: Create a standard bare mimium infromation row needed for emails sent
// Parameters:
// p *pgxpool.Pool -->
// inquiry_id string --> ID associated with each customer request for material
// customer_email string --> email that was handed in requesting construction material with
// material string --> material the customer has requested
// suppDetails SupplierDetails struct --> Supplier deatils struct
// Sent_Out bool --> flag confirming that the email has been sent out successfully
func AddBlankEmailInquiryEntry(p *pgxpool.Pool, inquiry_id, client_email, material string, suppDetails g.SupplierEmailInfo, sent_out bool, tableName string) (err error) {
	str := "INSERT INTO %s (inquiry_id, thread_id, client_email, time_sent, material, supplier_map_id, supplier_name, supplier_lat, supplier_lng, supplier_email, sent_out, present, price, currency, data_sheet) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)"
	sqlString := fmt.Sprintf(str, tableName)

	err = dataBaseTransmit(
		p,
		sqlString,
		inquiry_id,
		suppDetails.Email_ThreadID,
		client_email,
		time.Now(),
		material,
		suppDetails.MapsID,
		suppDetails.Name,
		suppDetails.Geometry.Location.Lat,
		suppDetails.Geometry.Location.Lng,
		suppDetails.Email[0],
		sent_out,
		false, //Present
		0,     //Price
		"",    //Currency
		nil,   //Data Sheet Pointer
	)
	if err != nil {
		return err
	}

	return nil
}

func ReadEmailInquiryEntry(p *pgxpool.Pool, tableName string, id_opt IDOption) (emailInquiry EmailInquiries, err error) {
	var sqlString string
	if id_opt.Inquiry_ID != "" {
		sqlString = fmt.Sprintf("SELECT * FROM %s WHERE inquiry_id = '%s'", tableName, id_opt.Inquiry_ID)
	} else {
		sqlString = fmt.Sprintf("SELECT * FROM %s WHERE thread_id = '%s'", tableName, id_opt.Thread_ID)
	}
	rows, err := dataBaseRead(p, sqlString)
	if err != nil {
		return EmailInquiries{}, err
	}

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			&emailInquiry.ID,
			&emailInquiry.Inquiry_ID,
			&emailInquiry.Thread_ID,
			&emailInquiry.Client_Email,
			&emailInquiry.Time_Sent,
			&emailInquiry.Material,
			&emailInquiry.Supplier_Map_ID,
			&emailInquiry.Supplier_Name,
			&emailInquiry.Supplier_Lat,
			&emailInquiry.Supplier_Lng,
			&emailInquiry.Supplier_Email,
			&emailInquiry.Sent_Out,
			&emailInquiry.Present,
			&emailInquiry.Price,
			&emailInquiry.Currency,
			&emailInquiry.Data_Sheet,
		)

		if err != nil {
			return EmailInquiries{}, err
		}
	}

	return emailInquiry, nil

}

func UpdateEmailInquiryEntryMaterialPresent(p *pgxpool.Pool, inquiry_id, tableName string, price float64, currency string, datasheet *[]byte) (err error) {
	sqlString := fmt.Sprintf("UPDATE %s SET present=$1, price=$2, currency=$3, data_sheet=$4 WHERE inquiry_id=$5", tableName)

	err = dataBaseTransmit(p, sqlString, true, price, currency, datasheet, inquiry_id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEmailInquiryEntryPresent(p *pgxpool.Pool, inquiry_id, tableName string, present bool) (err error) {
	sqlString := fmt.Sprintf("UPDATE %s SET present=$1 WHERE inquiry_id=$2", tableName)

	err = dataBaseTransmit(p, sqlString, true, inquiry_id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEmailInquiryEntryPrice(p *pgxpool.Pool, inquiry_id, tableName string, price float64, currency string) (err error) {
	sqlString := fmt.Sprintf("UPDATE %s SET price=$1, currency=$2 WHERE inquiry_id=$3", tableName)

	err = dataBaseTransmit(p, sqlString, price, currency, inquiry_id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEmailInquiryEntryDataSheet(p *pgxpool.Pool, inquiry_id, tableName string, data_sheet *[]byte) (err error) {
	sqlString := fmt.Sprintf("UPDATE %s SET data_sheet=$1 WHERE inquiry_id=$2", tableName)

	err = dataBaseTransmit(p, sqlString, *data_sheet, inquiry_id)
	if err != nil {
		return err
	}

	return nil
}

// func UpdateEmailInquiryEntryDataSheet(name string, pdfPath string, database string, p *pgxpool.Pool) (uint32, error) {
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
