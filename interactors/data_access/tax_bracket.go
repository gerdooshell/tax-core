package dataAccess

import (
	"context"
	sharedEntities "github.com/gerdooshell/tax-core/entities/canada/shared"
	"github.com/gerdooshell/tax-core/library/region/canada"
)

type TaxBracketsDataOut struct {
	TaxBrackets []sharedEntities.TaxBracket
	Err         error
}

type TaxBracketData interface {
	GetTaxBrackets(ctx context.Context, year int, province canada.Province) <-chan TaxBracketsDataOut
}
