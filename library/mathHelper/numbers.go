package mathHelper

import "math"

func RoundFloat64(num float64, precision int) float64 {
	coef := math.Pow(10, float64(precision))
	magnified := num * coef
	return float64(int(magnified+math.Copysign(0.5, magnified))) / coef
}
