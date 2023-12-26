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
	Data_Sheet      *uint32
} //15 items

// Purpose: Create a standard bare mimium infromation row needed for emails sent
// Parameters:
// p *pgxpool.Pool -->
// inquiry_id string --> ID associated with each customer request for material
// customer_email string --> email that was handed in requesting construction material with
// material string --> material the customer has requested
// suppDetails SupplierDetails struct --> Supplier deatils struct
// Sent_Out bool --> flag confirming that the email has been sent out successfully
func AddBlankEmailInquiryEntry(p *pgxpool.Pool, inquiry_id, client_email, material string, suppDetails g.SupplierInfo, sent_out bool, tableName string) (err error) {
	str := "INSERT INTO %s (inquiry_id, client_email, time_sent, material, supplier_map_id, supplier_name, supplier_lat, supplier_lng, supplier_email, sent_out, present, price, currency, data_sheet) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"
	sqlString := fmt.Sprintf(str, tableName)

	err = dataBaseTransmit(p, sqlString,
		inquiry_id,
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
		0,     //price
		"",    //currency
		nil,   //data sheet pointer
	)
	if err != nil {
		return err
	}

	return nil
}

func ReadEmailInquiryEntry(p *pgxpool.Pool, tableName string, inquiryID string) (emailInquiry EmailInquiries, err error) {
	sqlString := fmt.Sprintf("SELECT * FROM %s WHERE inquiry_id = '%s'", tableName, inquiryID)
	rows, err := dataBaseRead(p, sqlString)
	if err != nil {
		return EmailInquiries{}, err
	}

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			&emailInquiry.ID,
			&emailInquiry.Inquiry_ID,
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
