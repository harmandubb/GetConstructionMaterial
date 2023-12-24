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
	sqlString := fmt.Sprintf("INSERT INTO %s (Email, Time_Inquired, Material, Loc) VALUES($1, $2, $3, $4)", tableName)

	err = dataBaseTransmit(p, sqlString, database, matForm.Email, time.Now(), matForm.Material, matForm.Loc)
	if err != nil {
		return err
	}

	return nil
}
