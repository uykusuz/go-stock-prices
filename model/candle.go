package model

type Candle struct {
	Date             *Date
	Open             float64
	High             float64
	Low              float64
	Close            float64
	Volume           int64
	UnadjustedVolume int64
}
