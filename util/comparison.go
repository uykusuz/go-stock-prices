package util

import (
	"math"
)

const (
	Epsilon = float64(0.00001)
)

func EqualToEp(lhs float64, rhs float64) bool {
	return math.Abs(lhs-rhs) < Epsilon
}

func GreaterThanEp(lhs float64, rhs float64) bool {
	return (lhs - rhs) > Epsilon
}

func LessThanEp(lhs float64, rhs float64) bool {
	return (lhs - rhs) < -Epsilon
}
