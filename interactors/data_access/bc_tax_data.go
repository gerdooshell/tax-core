package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/bc/credits"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
)

type BCTaxData interface {
	AllCanadaTaxData
	GetBCBPA(ctx context.Context, year int) (<-chan credits.BasicPersonalAmount, <-chan error)
	GetBCTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error)
}
