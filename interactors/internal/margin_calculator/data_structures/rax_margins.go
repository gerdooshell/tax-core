package marginDS

import (
	"github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type Input struct {
	Province canada.Province
	Year     int
}

type Output struct {
	Brackets []shared.TaxBracket
	Err      error
}
