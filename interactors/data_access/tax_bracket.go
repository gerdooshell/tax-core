package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type TaxBracketData interface {
	GetTaxBrackets(ctx context.Context, year int, province canada.Province) <-chan TaxBracketsDataOut
}
