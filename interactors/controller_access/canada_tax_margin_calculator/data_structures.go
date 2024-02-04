package canadaTaxMarginCalculator

import (
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type Input struct {
	Year     int
	Province canada.Province
}

type Output struct {
	MarginalBrackets []shared.TaxBracket
}
