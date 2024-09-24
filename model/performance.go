package model

type Performance struct {
	Symbol       string
	SimpleReturn float64
	MaxDrawdown  float64
	From         *Date
	To           *Date
}
