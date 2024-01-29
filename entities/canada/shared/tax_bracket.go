package shared

import (
	"fmt"
	"math"
)

type TaxBracket struct {
	High float64
	Low  float64
	Rate float64
}

func FromArray(brackets, rates []float64) ([]TaxBracket, error) {
	result := make([]TaxBracket, 0, len(brackets))
	if len(brackets) != len(rates) {
		return result, fmt.Errorf("tax bracket error: there must be exactly one rate for each bracket")
	}
	var high float64
	for i := 0; i < len(brackets); i++ {
		if i+1 >= len(brackets) {
			high = math.MaxFloat64
		} else {
			high = brackets[i+1]
		}
		result = append(result, TaxBracket{
			High: high,
			Low:  brackets[i],
			Rate: rates[i],
		})
	}
	return result, nil
}
