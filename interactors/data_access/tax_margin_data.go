package dataAccess

import (
	"context"
	"github.com/gerdooshell/tax-core/entities/canada/shared"
)

type TaxMargin interface {
	GetAlbertaTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error)
	GetFederalTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error)
	GetBCTaxBrackets(ctx context.Context, year int) (<-chan []shared.TaxBracket, <-chan error)
}
