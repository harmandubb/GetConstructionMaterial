package database

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"googlemaps.github.io/maps"
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
}

type SupplierDetails struct {
	MapsID            string
	Supplier_Name     string
	Supplier_Geometry maps.AddressGeometry
	Supplier_Email    string
}

// Purpose: Create a standard bare mimium infromation row needed for emails sent
// Parameters:
// p *pgxpool.Pool -->
// inquiry_id string --> ID associated with each customer request for material
// customer_email string --> email that was handed in requesting construction material with
// material string --> material the customer has requested
// suppDetails SupplierDetails struct --> Supplier deatils struct
// Sent_Out bool --> flag confirming that the email has been sent out successfully
func AddBlackEmailInquiryEntry(p *pgxpool.Pool, inquiry_id, client_email, material string, suppDetails SupplierDetails, sent_out bool, tableName string) (err error) {
	str := "INSERT INTO %s (inquiry_id, client_email, time_sent, material, supplier_map_id, supplier_name, supplier_lat, supplier_lng, send_out, present, price, currency, data_sheet) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)"
	sqlString := fmt.Sprintf(str, tableName)

	err = dataBaseTransmit(p, sqlString,
		inquiry_id,
		client_email,
		time.Now(),
		material,
		suppDetails.MapsID,
		suppDetails.Supplier_Name,
		suppDetails.Supplier_Geometry.Location.Lat,
		suppDetails.Supplier_Geometry.Location.Lng,
		sent_out,
		false, //sent_out
		0,     //price
		"",    //currency
		nil,   //data sheet pointer
	)
	if err != nil {
		return err
	}

	return nil
}
