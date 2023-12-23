package database

import "time"

type CustomerInquiry struct {
	ID            int
	Email         string
	Time_Inquired time.Time
	Material      string
	Present       bool
	Price         float64
	Currency      string
	Data_Sheet    *uint32
}

func AddBlankInquiry