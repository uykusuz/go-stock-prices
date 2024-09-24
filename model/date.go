package model

import (
	"fmt"
	"strconv"
	"strings"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

func ParseDateFromString(str string) (*Date, error) {
	parts := strings.Split(str, "-")
	if len(parts) != 3 {
		return nil, fmt.Errorf("date string doesn't contain year month and value")
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse year")
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse month")
	}

	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse day")
	}

	return &Date{
		Year:  year,
		Month: month,
		Day:   day,
	}, nil
}

func (d *Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}
