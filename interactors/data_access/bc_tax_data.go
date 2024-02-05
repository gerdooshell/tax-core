package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/bc/credits"
)

type BCTaxData interface {
	AllCanadaTaxData
	GetBCBPA(ctx context.Context, year int) (<-chan credits.BasicPersonalAmount, <-chan error)
}
